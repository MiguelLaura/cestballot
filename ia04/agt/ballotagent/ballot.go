// package ballotagent

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"ia04/agt"
// 	"ia04/comsoc"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// type RestServerAgent struct {
// 	sync.Mutex
// 	id          string
// 	reqCount    int
// 	addr        string
// 	methodName  string
// 	methodSWF   func(comsoc.Profile) (comsoc.Count, error)
// 	methodSCF   func(comsoc.Profile) (comsoc.Alternative, error)
// 	orderedAlts []comsoc.Alternative
// 	profile     comsoc.Profile
// 	voterCount  int
// }

// func NewRestServerAgent(addr string) *RestServerAgent {
// 	return &RestServerAgent{id: addr, addr: addr}
// }

// // Test de la méthode
// func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
// 	if r.Method != method {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		fmt.Fprintf(w, "method %q not allowed", r.Method)
// 		return false
// 	}
// 	return true
// }

// func (*RestServerAgent) decodeRequestInit(r *http.Request) (req agt.RequestInit, err error) {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r.Body)
// 	err = json.Unmarshal(buf.Bytes(), &req)
// 	return
// }

// func (*RestServerAgent) decodeRequestVote(r *http.Request) (req agt.RequestVoter, err error) {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r.Body)
// 	err = json.Unmarshal(buf.Bytes(), &req)
// 	return
// }

// func (rsa *RestServerAgent) doInit(w http.ResponseWriter, r *http.Request) {
// 	// mise à jour du nombre de requêtes
// 	rsa.Lock()
// 	defer rsa.Unlock()
// 	rsa.reqCount++

// 	// vérification de la méthode de la requête
// 	if !rsa.checkMethod("POST", w, r) {
// 		return
// 	}

// 	// décodage de la requête
// 	req, err := rsa.decodeRequestInit(r)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	// traitement de la requête
// 	var voteMethSWF func(comsoc.Profile) (comsoc.Count, error)
// 	var voteMethSCF func(comsoc.Profile) (comsoc.Alternative, error)
// 	tieBreak := comsoc.TieBreakFactory(req.Candidates)

// 	switch req.Method {
// 	case "majority":
// 		voteMethSWF = comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)
// 		voteMethSCF = comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)
// 	case "borda":
// 		voteMethSWF = comsoc.SWFFactory(comsoc.BordaSWF, tieBreak)
// 		voteMethSCF = comsoc.SCFFactory(comsoc.BordaSCF, tieBreak)
// 	// case "approval":
// 	// 	vote = comsoc.SCFFactory(comsoc.ApprovalSCF, tieBreak)
// 	// case "condorcet":
// 	// 	voteMeth = comsoc.SCFFactory(comsoc.CondorcetWinner, tieBreak)
// 	case "copeland":
// 		voteMethSWF = comsoc.SWFFactory(comsoc.CopelandSWF, tieBreak)
// 		voteMethSCF = comsoc.SCFFactory(comsoc.CopelandSCF, tieBreak)
// 	case "STV":
// 		voteMethSWF = comsoc.SWFFactory(comsoc.STV_SWF, tieBreak)
// 		voteMethSCF = comsoc.SCFFactory(comsoc.STV_SCF, tieBreak)
// 	default:
// 		w.WriteHeader(http.StatusNotImplemented)
// 		msg := fmt.Sprintf("Unkonwn vote method '%s'", req.Method)
// 		w.Write([]byte(msg))
// 		return
// 	}

// 	rsa.methodName = req.Method
// 	rsa.methodSWF = voteMethSWF
// 	rsa.methodSCF = voteMethSCF

// 	// Check pas deux fois même candidat A FAIRE
// 	rsa.orderedAlts = req.Candidates

// 	w.WriteHeader(http.StatusOK)
// 	serial, _ := json.Marshal("Bureau initialisé")
// 	w.Write(serial)
// }

// func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
// 	// mise à jour du nombre de requêtes
// 	rsa.Lock()
// 	defer rsa.Unlock()
// 	rsa.voterCount++
// 	rsa.reqCount++

// 	// vérification de la méthode de la requête
// 	if !rsa.checkMethod("POST", w, r) {
// 		return
// 	}

// 	// décodage de la requête
// 	req, err := rsa.decodeRequestVote(r)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	// check prefs A FAIRE
// 	rsa.profile = append(rsa.profile, req.Prefs)

// 	w.WriteHeader(http.StatusOK)
// 	serial, _ := json.Marshal("Vote pris en compte")
// 	w.Write(serial)

// }

// func (rsa *RestServerAgent) doResult(w http.ResponseWriter, r *http.Request) {
// 	if !rsa.checkMethod("GET", w, r) {
// 		return
// 	}

// 	var resp agt.Response
// 	resp.Ranking, _ = rsa.methodSWF(rsa.profile)
// 	resp.Winner, _ = rsa.methodSCF(rsa.profile)

// 	w.WriteHeader(http.StatusOK)
// 	rsa.Lock()
// 	defer rsa.Unlock()
// 	serial, _ := json.Marshal(resp)
// 	w.Write(serial)
// }

// func (rsa *RestServerAgent) doReqcount(w http.ResponseWriter, r *http.Request) {
// 	if !rsa.checkMethod("GET", w, r) {
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	rsa.Lock()
// 	defer rsa.Unlock()
// 	serial, _ := json.Marshal(rsa.reqCount)
// 	w.Write(serial)
// }

// func (rsa *RestServerAgent) doVoterscount(w http.ResponseWriter, r *http.Request) {
// 	if !rsa.checkMethod("GET", w, r) {
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	rsa.Lock()
// 	defer rsa.Unlock()
// 	serial, _ := json.Marshal(rsa.voterCount)
// 	w.Write(serial)
// }

// func (rsa *RestServerAgent) Start() {
// 	// création du multiplexer
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/init", rsa.doInit)
// 	mux.HandleFunc("/vote", rsa.doVote)
// 	mux.HandleFunc("/result", rsa.doResult)
// 	mux.HandleFunc("/reqcount", rsa.doReqcount)
// 	mux.HandleFunc("/voterscount", rsa.doVoterscount)

// 	// création du serveur http
// 	s := &http.Server{
// 		Addr:           rsa.addr,
// 		Handler:        mux,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20}

// 	// lancement du serveur
// 	log.Println("Listening on", rsa.addr)
// 	go log.Fatal(s.ListenAndServe())
// }
