package restserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	ba "gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
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
		w.WriteHeader(METH_NOT_IMPL)
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
		w.WriteHeader(BAD_REQUEST)
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
			w.WriteHeader(BAD_REQUEST)
			fmt.Fprint(w, "JSON incorrect content")
		case "501":
			w.WriteHeader(NOT_IMPL)
			fmt.Fprintf(w, "vote method %q not supported", req.Rule)
		}
		return
	}

	err = theNewBallot.Start()
	if err != nil {
		log.Println("Ballot newBallotId : " + err.Error())
		w.WriteHeader(BAD_REQUEST)
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

func (rst *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	if !rst.checkMethod("POST", w, r) {
		log.Println("doVote : request is not of type POST")
		return
	}

	var req VoteRequest

	err := decodeRequest(r, &req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprint(w, "JSON request string incorrect format")
		return
	}

	rst.Lock()
	defer rst.Unlock()

	ballotAgent, exists := rst.ballots[req.Ballot]

	if !exists {
		log.Println("Error, ballot " + req.Ballot + " does not exist")
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprintf(w, "JSON ballot %q does not exist", req.Ballot)
		return
	}

	_, err = ballotAgent.Vote(req.Agent, req.Prefs, req.Options)
	if err != nil {
		log.Println("Vote : " + err.Error())
		switch strings.Split(err.Error(), "::")[0] {
		case "1":
			w.WriteHeader(DEADLINE_OVER)
			fmt.Fprint(w, "The Deadline is now over")
		case "2":
			w.WriteHeader(BAD_REQUEST)
			fmt.Fprint(w, "The voter cannot vote here")
		case "3":
			w.WriteHeader(VOTE_ALREADY_DONE)
			fmt.Fprint(w, "The voter has already voted here")
		case "4":
			w.WriteHeader(BAD_REQUEST)
			fmt.Fprint(w, "The voter has not the right preferences")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Vote accepted !")
}

func (rst *RestServerAgent) doResult(w http.ResponseWriter, r *http.Request) {
	if !rst.checkMethod("POST", w, r) {
		log.Println("doVote : request is not of type POST")
		return
	}

	var req ResultRequest

	err := decodeRequest(r, &req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprint(w, "JSON request string incorrect format")
		return
	}

	rst.Lock()
	defer rst.Unlock()

	ballotAgent, exists := rst.ballots[req.Ballot]

	if !exists {
		log.Println("Error, ballot " + req.Ballot + " does not exist")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "JSON ballot %q does not exist", req.Ballot)
		return
	}

	res, err := ballotAgent.GetVoteResult()
	if err != nil {
		log.Println("Result : " + err.Error())
		switch strings.Split(err.Error(), "::")[0] {
		case "1":
			w.WriteHeader(TOO_EARLY)
			fmt.Fprint(w, "It's too early for the result, wait a bit")
		}
		return
	}

	if res == 0 {
		log.Println("Result not found...")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Result not found")
	}

	var resp ResultResponse
	resp.Winner = res

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
	mux.HandleFunc("/vote", rst.doVote)
	mux.HandleFunc("/result", rst.doResult)

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
