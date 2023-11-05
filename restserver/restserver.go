// Package restserver handling the REST server allowing to vote
package restserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	ba "ia04-tp/agt/ballotagent"
)

// RestServerAgent the agent representing the REST server
type RestServerAgent struct {
	sync.Mutex
	id      string                         // The ID of the REST server
	addr    string                         // The url of the server
	ballots map[string]*ba.RestBallotAgent // The ballots
	verbose bool                           // <tt>true</tt> if the server is verbose, <tt>false</tt> else
}

// NewRestServerAgent creates a new REST Server at a given address.
func NewRestServerAgent(id string, addr string) *RestServerAgent {
	rst := RestServerAgent{id: id, addr: addr}
	rst.ballots = make(map[string]*ba.RestBallotAgent)
	return &rst
}

// checkMethod tests the http verb used in the request.
func (rst *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(HTTP_VERB_INCORRECT)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (rst *RestServerAgent) SetVerbose(verbose bool) {
	rst.verbose = verbose
}

func (rst *RestServerAgent) log(messages ...any) {
	if rst.verbose {
		log.Println(messages...)
	}
}

/*
	----------------------------------------
					Decoders
	----------------------------------------
*/

// decodeRequest decodes the request to a specific structure.
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

// doNewBallot handles the new-ballot request by creating a new ballot.
//
// The server returns the following status code :
// - VOTE_CREATED (201)        : OK
// - BAD_REQUEST (400)         : incorrect Request
// - HTTP_VERB_INCORRECT (405) : the request is not a POST request
// - NOT_IMPL (501)            : if the voting method is not implemented
func (rst *RestServerAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	if !rst.checkMethod("POST", w, r) {
		rst.log("doNewBallot : request is not of type POST")
		return
	}

	// ---
	// Handle request

	var req NewBallotRequest

	err := decodeRequest(r, &req)
	if err != nil {
		rst.log("doNewBallot : ", err.Error())
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprint(w, "JSON request string incorrect format")
		return
	}

	rst.Lock()
	defer rst.Unlock()

	newBallotId := fmt.Sprintf("scrutin%d", len(rst.ballots)+1)

	theNewBallot, err := ba.NewRestBallotAgent(
		newBallotId,
		req.Rule,
		req.Deadline,
		req.Voters,
		req.Alts,
		req.TieBreak,
	)

	if err != nil {
		rst.log("doNewBallot : ", strings.Split(err.Error(), "::")[1])
		switch strings.Split(err.Error(), "::")[0] {
		case "1", "2", "3", "4":
			w.WriteHeader(BAD_REQUEST)
			fmt.Fprint(w, "JSON incorrect content")
		case "5":
			w.WriteHeader(NOT_IMPL)
			fmt.Fprintf(w, "vote method %q not supported", req.Rule)
		}
		return
	}

	// ---
	// Handle response

	err = theNewBallot.Start()
	if err != nil {
		rst.log("doNewBallot : ", err.Error())
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprint(w, "The deadline is in the past")
		return
	}

	rst.ballots[theNewBallot.ID] = theNewBallot

	var resp NewBallotResponse
	resp.Id = theNewBallot.ID

	w.WriteHeader(VOTE_CREATED)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
	rst.log("doNewBallot : ballot '" + newBallotId + "' created")
}

// doVote handles the vote request from a voting agent.
//
// The server returns the following status code :
// - VOTE_TAKEN (200)          : OK
// - BAD_REQUEST (400)         : incorrect Request
// - VOTE_ALREADY_DONE (403)   : vote already done by the agent
// - HTTP_VERB_INCORRECT (405) : the request is not a POST request
// - NOT_IMPL (501)            : if the voting method is not implemented
// - DEADLINE_OVER (503)       : when the deadline is over
func (rst *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	if !rst.checkMethod("POST", w, r) {
		rst.log("doVote : request is not of type POST")
		return
	}

	// ---
	// Handle request

	var req VoteRequest

	err := decodeRequest(r, &req)
	if err != nil {
		rst.log("doVote : ", err.Error())
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprint(w, "JSON request string incorrect format")
		return
	}

	rst.Lock()
	defer rst.Unlock()

	ballotAgent, exists := rst.ballots[req.Ballot]

	if !exists {
		rst.log("doVote : ballot '" + req.Ballot + "' does not exist")
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprintf(w, "JSON ballot %q does not exist", req.Ballot)
		return
	}

	// ---
	// Handle response

	_, err = ballotAgent.Vote(req.Agent, req.Prefs, req.Options)
	if err != nil {
		rst.log("doVote : ", strings.Split(err.Error(), "::")[1])
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

	w.WriteHeader(VOTE_TAKEN)
	fmt.Fprint(w, "Vote accepted !")
	rst.log(fmt.Sprintf("doVote : vote accepted for agent %q to ballot %q", req.Agent, req.Ballot))
}

// doResult handles the result request to a specific ballot.
//
// The server returns the following status code :
// - RESULT_OBTAINED (200)     : OK
// - NOT_FOUND (404)           : the ballot or the result is not found
// - HTTP_VERB_INCORRECT (405) : the request is not a POST request
// - TOO_EARLY (425)           : the deadline of the ballot is not passed
func (rst *RestServerAgent) doResult(w http.ResponseWriter, r *http.Request) {
	if !rst.checkMethod("POST", w, r) {
		rst.log("doResult : request is not of type POST")
		return
	}

	// ---
	// Handle request

	var req ResultRequest

	err := decodeRequest(r, &req)
	if err != nil {
		rst.log("doResult : ", err.Error())
		w.WriteHeader(BAD_REQUEST)
		fmt.Fprint(w, "JSON request string incorrect format")
		return
	}

	rst.Lock()
	defer rst.Unlock()

	ballotAgent, exists := rst.ballots[req.Ballot]

	if !exists {
		rst.log("doVote : ballot '" + req.Ballot + "' does not exist")
		w.WriteHeader(NOT_FOUND)
		fmt.Fprintf(w, "JSON ballot %q does not exist", req.Ballot)
		return
	}

	// ---
	// Handle response

	winner, ranking, err := ballotAgent.GetVoteResult()
	if err != nil {
		switch strings.Split(err.Error(), "::")[0] {
		case "1":
			w.WriteHeader(TOO_EARLY)
			fmt.Fprint(w, "It's too early for the result, wait a bit")
		case "2":
			w.WriteHeader(NOT_FOUND)
			fmt.Fprint(w, "Vote not found")
			fallthrough
		default:
			rst.log("doResult : ", strings.Split(err.Error(), "::")[1])
		}
		return
	}

	var resp ResultResponse
	resp.Winner = winner
	resp.Ranking = ranking

	w.WriteHeader(RESULT_OBTAINED)
	serial, _ := json.Marshal(resp)
	w.Write(serial)

	rst.log(fmt.Sprintf("doResult : for ballot %q, the result is %d with the following ranking : %v", req.Ballot, winner, ranking))
}

/*
----------------------------------------

	Start

----------------------------------------
*/

// Start starts the REST server.
func (rst *RestServerAgent) Start() {
	// Create the multiplexer to redirect the requests
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rst.doNewBallot)
	mux.HandleFunc("/vote", rst.doVote)
	mux.HandleFunc("/result", rst.doResult)

	// Creates the HTTP server
	s := &http.Server{
		Addr:           rst.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// Starts the server
	log.Println("Listening on", rst.addr)
	go log.Fatal(s.ListenAndServe())
}
