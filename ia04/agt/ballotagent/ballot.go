package ballotagent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ia04/agt"
	"ia04/comsoc"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"
)

type BallotAgent struct {
	ballotId    string
	rule        string
	deadline    time.Time
	voterIds    []string
	alts        int
	tieBreak    []comsoc.Alternative
	options     []int
	profile     comsoc.Profile
	voterIdDone []string
}

type RestServerAgent struct {
	sync.Mutex
	addr    string
	ballots map[string]*BallotAgent
}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{addr: addr, ballots: make(map[string]*BallotAgent)}
}

func NewBallotAgent(rsa *RestServerAgent) *BallotAgent {
	ballotId := fmt.Sprintf("scrutin%02d", len(rsa.ballots)+1)
	return &BallotAgent{ballotId: ballotId}
}

// Test de la méthode
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestServerAgent) decodeRequestBallot(r *http.Request) (req agt.RequestBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) decodeRequestVoter(r *http.Request) (req agt.RequestVoter, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) decodeRequestResult(r *http.Request) (req agt.RequestResult, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServerAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestBallot(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// traitement de la requête

	// Méthode de vote
	if !slices.Contains([]string{"majority", "borda", "approval", "condorcet", "copeland", "STV"}, req.Rule) {
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("erreur : règle de vote inconnue '%s'", req.Rule)
		w.Write([]byte(msg))
		return
	}

	// Deadline
	timeParsed, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		err = errors.New("erreur : mauvais format de deadline")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if time.Now().After(timeParsed) {
		err = errors.New("erreur : la deadline est déjà passée")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// VoterIds
	if len(req.VoterIds) == 0 {
		err = errors.New("erreur : pas assez de voters")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Alts
	if req.Alts <= 1 {
		err = errors.New("erreur : pas suffisamment d'alternatives : " + fmt.Sprint(req.Alts))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// TieBreak
	err = checkTieBreak(req.Alts, req.TieBreak)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Créer ballot
	ballot := NewBallotAgent(rsa)
	rsa.ballots[ballot.ballotId] = ballot
	ballot.rule = req.Rule
	ballot.deadline = timeParsed
	ballot.voterIds = req.VoterIds
	ballot.alts = req.Alts
	ballot.tieBreak = req.TieBreak

	// Gérer la réponse sans erreur
	var resp agt.ResponseBallot
	resp.BallotId = ballot.ballotId
	w.WriteHeader(http.StatusCreated)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestVoter(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// BallotId
	if req.BallotId == "" {
		err = errors.New("erreur : il manque l'id du ballot")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	ballot := rsa.ballots[req.BallotId]
	if ballot == nil {
		err = errors.New("erreur : le ballot n'existe pas encore")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// AgentId
	if req.AgentId == "" {
		err = errors.New("erreur : il manque l'id de l'agent")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if slices.Contains(ballot.voterIdDone, req.AgentId) {
		err = errors.New("erreur : l'agent " + fmt.Sprint(req.AgentId) + " a déjà voté.")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, err.Error())
		return
	}

	// Deadline
	if time.Now().After(ballot.deadline) {
		err = errors.New("erreur : la deadline est déjà passée")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, err.Error())
		return
	}

	// Prefs
	err = checkPrefs(ballot.alts, req.Prefs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Options
	if ballot.rule == "approval" && len(req.Options) != 1 {
		err = errors.New("erreur : Options est forcément de la forme [1]int")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	ballot.profile = append(ballot.profile, req.Prefs)
	ballot.options = append(ballot.options, req.Options...)
	ballot.voterIdDone = append(ballot.voterIdDone, req.AgentId)

	w.WriteHeader(http.StatusOK)
}

func (rsa *RestServerAgent) doResult(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestResult(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// BallotId
	if req.BallotId == "" {
		err = errors.New("erreur : il manque l'id du ballot")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	ballot := rsa.ballots[req.BallotId]
	if ballot == nil {
		err = errors.New("erreur : le ballot n'existe pas")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		return
	}

	// Deadline
	if time.Now().Before(ballot.deadline) {
		err = errors.New("erreur : la deadline n'est pas encore passée")
		w.WriteHeader(http.StatusTooEarly)
		fmt.Fprint(w, err.Error())
		return
	}

	// Methode de vote
	var resp agt.ResponseResult
	tieBreak := comsoc.TieBreakFactory(ballot.tieBreak)
	switch ballot.rule {
	case "majority":
		resp.Ranking, _ = comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)(ballot.profile)
		resp.Winner, _ = comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)(ballot.profile)
	case "borda":
		resp.Ranking, _ = comsoc.SWFFactory(comsoc.BordaSWF, tieBreak)(ballot.profile)
		resp.Winner, _ = comsoc.SCFFactory(comsoc.BordaSCF, tieBreak)(ballot.profile)
	// [A VERIFIER]
	case "approval":
		resp.Ranking, _ = comsoc.SWFFactoryApproval(comsoc.ApprovalSWF, tieBreak)(ballot.profile, ballot.options)
		resp.Winner, _ = comsoc.SCFFactoryApproval(comsoc.ApprovalSCF, tieBreak)(ballot.profile, ballot.options)
	case "condorcet":
		// [A VERIFIER]
		resp.Ranking = nil
		resp.Winner, _ = comsoc.SCFFactory(comsoc.CondorcetWinner, tieBreak)(ballot.profile)
	case "copeland":
		resp.Ranking, _ = comsoc.SWFFactory(comsoc.CopelandSWF, tieBreak)(ballot.profile)
		resp.Winner, _ = comsoc.SCFFactory(comsoc.CopelandSCF, tieBreak)(ballot.profile)
	case "STV":
		resp.Ranking, _ = comsoc.SWFFactory(comsoc.STV_SWF, tieBreak)(ballot.profile)
		resp.Winner, _ = comsoc.SCFFactory(comsoc.STV_SCF, tieBreak)(ballot.profile)
	}

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doNewBallot)
	mux.HandleFunc("/vote", rsa.doVote)
	mux.HandleFunc("/result", rsa.doResult)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
