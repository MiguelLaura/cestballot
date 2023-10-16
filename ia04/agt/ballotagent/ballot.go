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
	"sync"
	"time"
)

type RestBallotAgent struct {
	sync.Mutex
	ballotId  string
	addr      string
	rule      string
	deadline  string
	voterIds  []string
	alts      int
	tieBreak  []comsoc.Alternative
	methodSWF func(comsoc.Profile) (comsoc.Count, error)
	methodSCF func(comsoc.Profile) (comsoc.Alternative, error)
}

// [A FAIRE] numéro de scrutin
func NewRestBallotAgent(addr string) *RestBallotAgent {
	return &RestBallotAgent{ballotId: "scrutin", addr: addr}
}

// Test de la méthode
func (rba *RestBallotAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestBallotAgent) decodeRequestNewBallot(r *http.Request) (req agt.RequestBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rba *RestBallotAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	rba.Lock()
	defer rba.Unlock()

	// vérification de la méthode de la requête
	if !rba.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rba.decodeRequestNewBallot(r)
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
		rba.rule = req.Rule
		voteMethSWF = comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)
	case "borda":
		rba.rule = req.Rule
		voteMethSWF = comsoc.SWFFactory(comsoc.BordaSWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.BordaSCF, tieBreak)
	// [A FAIRE]
	// case "approval":
	// 	vote = comsoc.SCFFactory(comsoc.ApprovalSCF, tieBreak)
	// case "condorcet":
	// 	voteMeth = comsoc.SCFFactory(comsoc.CondorcetWinner, tieBreak)
	case "copeland":
		rba.rule = req.Rule
		voteMethSWF = comsoc.SWFFactory(comsoc.CopelandSWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.CopelandSCF, tieBreak)
	case "STV":
		rba.rule = req.Rule
		voteMethSWF = comsoc.SWFFactory(comsoc.STV_SWF, tieBreak)
		voteMethSCF = comsoc.SCFFactory(comsoc.STV_SCF, tieBreak)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("Règle de vote inconnue '%s'", req.Rule)
		w.Write([]byte(msg))
		return
	}
	rba.rule = req.Rule
	rba.methodSWF = voteMethSWF
	rba.methodSCF = voteMethSCF

	// Deadline
	_, err = time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		err = errors.New("erreur : mauvais format de deadline")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	rba.deadline = req.Deadline

	// VoterIds
	if len(req.VoterIds) == 0 {
		err = errors.New("erreur : pas assez de voters")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	rba.voterIds = req.VoterIds

	// Alts
	if req.Alts <= 1 {
		err = errors.New("erreur : pas suffisamment d'alternatives : " + fmt.Sprint(req.Alts))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	rba.alts = req.Alts

	// TieBreak
	err = checkTieBreak(rba.alts, req.TieBreak)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	rba.tieBreak = req.TieBreak

	// Gérer la réponse sans erreur
	var resp agt.ResponseBallot
	resp.BallotId = rba.ballotId
	w.WriteHeader(http.StatusCreated)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rba *RestBallotAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rba.doNewBallot)

	// création du serveur http
	s := &http.Server{
		Addr:           rba.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rba.addr)
	go log.Fatal(s.ListenAndServe())
}
