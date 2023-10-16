package ballotagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ia04/agt"
	"ia04/comsoc"
	"log"
	"net/http"
	"sync"
	"time"
)

type RestServerAgent struct {
	sync.Mutex
	id        string
	addr      string
	rule      string
	deadline  string
	voterIds  []string
	alts      int
	tieBreak  []comsoc.Alternative
	methodSWF func(comsoc.Profile) (comsoc.Count, error)
	methodSCF func(comsoc.Profile) (comsoc.Alternative, error)
}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{id: addr, addr: addr}
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

func (*RestServerAgent) decodeRequestNewBallot(r *http.Request) (req agt.RequestBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServerAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestNewBallot(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// traitement de la requête
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

	rsa.rule = req.Rule
	rsa.methodSWF = voteMethSWF
	rsa.methodSCF = voteMethSCF

	// [A FAIRE]
	// Check pas deux fois même candidat
	rsa.tieBreak = req.TieBreak

	w.WriteHeader(http.StatusCreated)
	serial, _ := json.Marshal("Vote créé")
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doNewBallot)

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
