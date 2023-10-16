package restserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	ba "ia04/agt/ballotagent"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RestServerAgent struct {
	sync.Mutex
	id       string
	nbBallot int
	addr     string
	ballots  map[string]*ba.RestBallotAgent
}

// Test du verbe HTTP utilisé
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func NewRestServerAgent(id string, addr string) *RestServerAgent {
	rst := RestServerAgent{id: id, addr: addr}
	rst.ballots = make(map[string]*ba.RestBallotAgent)
	return &rst
}

func (rst *RestServerAgent) genBallotAgentId() string {
	rst.Lock()
	defer rst.Unlock()
	defer func() { rst.nbBallot++ }()

	return fmt.Sprintf("vote%d", rst.nbBallot)
}

/*
	----------------------------------------
					Decoders
	----------------------------------------
*/

func decodeRequest[T any](r *http.Request, req *T) (err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

/*
	----------------------------------------
					Handlers
	----------------------------------------
*/

func (rst *RestServerAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	if !rst.checkMethod("POST", w, r) {
		log.Println("doNewBallot : request is not of type POST")
		return
	}

	var req NewBallotRequest

	err := decodeRequest(r, &req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "JSON request string incorrect format")
		return
	}

	newBallotId := rst.genBallotAgentId()

	theNewBallot, err := ba.NewRestBallotAgent(
		newBallotId,
		req.Rule,
		req.Deadline,
		req.Voters,
		req.Alts,
		req.TieBreak,
	)

	if err != nil {
		log.Println("Ballot newBallotId : " + err.Error())
		switch strings.Split(err.Error(), "::")[0] {
		case "400":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "JSON incorrect content")
		case "501":
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprintf(w, "vote method %q not supported", req.Rule)
		}
		return
	}

	err = theNewBallot.Start()
	if err != nil {
		log.Println("Ballot newBallotId : " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "The deadline is in the past")
		return
	}

	rst.ballots[theNewBallot.ID] = theNewBallot

	var resp NewBallotResponse
	resp.Id = theNewBallot.ID

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

/*
----------------------------------------

	Start

----------------------------------------
*/
func (rst *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new-ballot", rst.doNewBallot)

	// création du serveur http
	s := &http.Server{
		Addr:           rst.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rst.addr)
	go log.Fatal(s.ListenAndServe())
}
