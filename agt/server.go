package agt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ia04-tp/comsoc"
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
	winner      comsoc.Alternative
	ranking     []comsoc.Alternative
}

type RestServerAgent struct {
	sync.Mutex
	addr    string
	ballots map[string]*BallotAgent
}

func NewBallotAgent(rsa *RestServerAgent) *BallotAgent {
	return &BallotAgent{ballotId: fmt.Sprintf("scrutin%02d", len(rsa.ballots)+1)}
}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{addr: addr, ballots: make(map[string]*BallotAgent)}
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

func (*RestServerAgent) decodeRequestBallot(r *http.Request) (req RequestBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) decodeRequestVoter(r *http.Request) (req RequestVoter, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) decodeRequestResult(r *http.Request) (req RequestResult, err error) {
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
	if !slices.Contains([]string{"majority", "borda", "approval", "condorcet", "copeland", "stv"}, req.Rule) {
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
		err = errors.New("erreur : la deadline est déjà passée : " + fmt.Sprint(timeParsed))
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
	ballot.winner = -1

	// Gérer la réponse sans erreur
	var resp ResponseBallot
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
	if !slices.Contains(ballot.voterIds, req.AgentId) {
		err = errors.New("erreur : l'agent " + fmt.Sprint(req.AgentId) + " ne peut pas voter dans ce ballot.")
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
		err = errors.New("erreur : la deadline est déjà passée : " + fmt.Sprint(ballot.deadline))
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
		err = errors.New("erreur : la deadline n'est pas encore passée : " + fmt.Sprint(ballot.deadline))
		w.WriteHeader(http.StatusTooEarly)
		fmt.Fprint(w, err.Error())
		return
	}

	// Methode de vote
	var resp ResponseResult
	if ballot.winner != -1 {
		resp.Ranking = ballot.ranking
		resp.Winner = ballot.winner
	} else if ballot.profile == nil {
		if ballot.rule == "condorcet" {
			resp.Ranking = nil
			resp.Winner = ballot.tieBreak[0]
		} else {
			resp.Ranking = ballot.tieBreak
			resp.Winner = ballot.tieBreak[0]
		}
	} else {
		tieBreak := comsoc.TieBreakFactory(ballot.tieBreak)
		switch ballot.rule {
		case "majority":
			resp.Ranking, err = comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
			resp.Winner, err = comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
		case "borda":
			resp.Ranking, err = comsoc.SWFFactory(comsoc.BordaSWF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
			resp.Winner, err = comsoc.SCFFactory(comsoc.BordaSCF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
		case "approval":
			resp.Ranking, err = comsoc.SWFFactoryApproval(comsoc.ApprovalSWF, tieBreak)(ballot.profile, ballot.options)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
			resp.Winner, err = comsoc.SCFFactoryApproval(comsoc.ApprovalSCF, tieBreak)(ballot.profile, ballot.options)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
		case "condorcet":
			var res []comsoc.Alternative
			res, err = comsoc.CondorcetWinner(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
			if len(res) != 0 {
				resp.Winner = res[0]
			} else {
				resp.Winner = 0
			}
		case "copeland":
			resp.Ranking, err = comsoc.SWFFactory(comsoc.CopelandSWF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
			resp.Winner, err = comsoc.SCFFactory(comsoc.CopelandSCF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
		case "stv":
			resp.Ranking, err = comsoc.SWFFactory(comsoc.STV_SWF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
			resp.Winner, err = comsoc.SCFFactory(comsoc.STV_SCF, tieBreak)(ballot.profile)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}
		}
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
