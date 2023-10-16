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

// [A FAIRE]

// 503 la deadline est dépassée
// 425 Too early
// 404 Not Found

type BallotAgent struct {
	ballotId    string
	rule        string
	deadline    string
	voterIds    []string
	alts        int
	tieBreak    []comsoc.Alternative
	methodSWF   func(comsoc.Profile) (comsoc.Count, error)
	methodSCF   func(comsoc.Profile) (comsoc.Alternative, error)
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
	ballotId := fmt.Sprintf("scrutin%02d", len(rsa.ballots))
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

	// Méthode de vite
	var voteMethSWF func(comsoc.Profile) (comsoc.Count, error)
	var voteMethSCF func(comsoc.Profile) (comsoc.Alternative, error)
	tieBreak := comsoc.TieBreakFactory(req.TieBreak)
	switch req.Rule {
	case "majority":
		voteMethSWF = comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)
	case "borda":
		voteMethSWF = comsoc.SWFFactory(comsoc.BordaSWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.BordaSCF, tieBreak)
	// [A FAIRE]
	// case "approval":
	// 	vote = comsoc.SCFFactory(comsoc.ApprovalSCF, tieBreak)
	// case "condorcet":
	// 	voteMeth = comsoc.SCFFactory(comsoc.CondorcetWinner, tieBreak)
	case "copeland":
		voteMethSWF = comsoc.SWFFactory(comsoc.CopelandSWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.CopelandSCF, tieBreak)
	case "STV":
		voteMethSWF = comsoc.SWFFactory(comsoc.STV_SWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.STV_SCF, tieBreak)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("Règle de vote inconnue '%s'", req.Rule)
		w.Write([]byte(msg))
		return
	}

	// Deadline
	_, err = time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		err = errors.New("erreur : mauvais format de deadline")
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
	ballot.methodSWF = voteMethSWF
	ballot.methodSCF = voteMethSCF
	ballot.deadline = req.Deadline
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

	// Prefs
	err = checkPrefs(ballot.alts, req.Prefs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	ballot.profile = append(ballot.profile, req.Prefs)
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
		err = errors.New("erreur : le ballot n'existe pas encore")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		return
	}

	var resp agt.ResponseResult
	resp.Ranking, _ = ballot.methodSWF(ballot.profile)
	resp.Winner, _ = ballot.methodSCF(ballot.profile)

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
