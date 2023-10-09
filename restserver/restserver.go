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

func NewRestServerAgent(id string, addr string) *RestServerAgent {
	return &RestServerAgent{id: id, addr: addr}
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
		return
	}

	var req NewBallotRequest

	err := decodeRequest(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "JSON request string incorrect format")
		return
	}

	theNewBallot, err := ba.NewRestBallotAgent(
		rst.genBallotAgentId(),
		req.Rule,
		req.Deadline,
		req.Voters,
		req.Alts,
	)

	if err != nil {
		switch strings.Split(err.Error(), "::")[0] {
		case "400":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "JSON date string incorrect format")
		case "501":
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprintf(w, "vote method %q not supported", req.Rule)
		}
		return
	}

	err = theNewBallot.Start()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "The deadline is in the past")
	}

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
