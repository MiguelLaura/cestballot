// Package voteragent contains an agent that can vote.
package voteragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
	rs "gitlab.utc.fr/mennynat/ia04-tp/restserver"
	"gitlab.utc.fr/mennynat/ia04-tp/utils/concurrent"
)

// The voting agent structure.
type RestVoterAgent struct {
	agt.Agent
	url  string // The url of the server where the voter votes
	opts []int  // The optional data to send to specific voting methods
}

// NewRestVoterAgent creates a new voting agent.
func NewRestVoterAgent(id string, name string, prefs []comsoc.Alternative, opts []int, url string) *RestVoterAgent {
	return &RestVoterAgent{agt.Agent{ID: id, Name: name, Prefs: prefs}, url, opts}
}

// likeLevelOfVote gives a level of likeness of the agent and the given alternative.
//
// The likeness is an int from -1 to 3 with the following meaning :
// - -1 : unknown
// - 0 : likes it a lot
// - 1 : likes a bit
// - 2 : dislikes it a bit
// - 3 : dislikes it a lot
func (agent *RestVoterAgent) likeLevelOfAlt(alt comsoc.Alternative) int {
	altRank := 0
	for altRank < len(agent.Prefs) && agent.Prefs[altRank] != alt {
		altRank++
	}

	if altRank >= len(agent.Prefs) {
		return -1
	}

	quarter, half := int(len(agent.Prefs)/4), int(len(agent.Prefs)/2)

	if altRank < quarter {
		return 0
	} else if altRank < half {
		return 1
	} else if altRank < half+quarter {
		return 2
	} else {
		return 3
	}
}

/*
	Basic Agent methods (from AgentI interface)
*/

// Equal indicates if two RestVoterAgents have the same ID, Name and server url.
func (agent *RestVoterAgent) Equal(agent2 RestVoterAgent) bool {
	return agent.ID == agent2.ID && agent.Name == agent2.Name && agent.url == agent2.url
}

// DeepEqual indicates if two RestVoterAgents have the same ID, Name, server url and preferences.
func (agent *RestVoterAgent) DeepEqual(agent2 RestVoterAgent) bool {
	return agent.Equal(agent2) && concurrent.Equal(agent.Prefs, agent2.Prefs)
}

// Clone creates a duplicate of the current agent.
func (agent *RestVoterAgent) Clone() RestVoterAgent {
	return *NewRestVoterAgent(agent.ID, agent.Name, agent.Prefs, agent.opts, agent.url)
}

// String gives a string representation of the voting agent.
func (agent *RestVoterAgent) String() string {
	return fmt.Sprintf("%s : %s", agent.ID, agent.Name)
}

// Prefers indicates if the agent prefers the alternative a to the b.
func (agent *RestVoterAgent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	return isPref(a, b, agent.Prefs)
}

/*
	HTTP response decoder
*/

// decodeResponse decodes the received http response into a specific struct
func decodeResponse[T any](r *http.Response, req *T) (err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

/*
	HTTP requesters
*/

// DoNewBallot send a new-ballot request to the server at servUrl.
//
// If the request fails or cannot be made, an error is returned with a specific format : errType::errMessage.
// The errType can be one of the following :
// - 0   : if there's no voter
// - 1   : if there's a problem when converting request to json
// - 2   : if the request cannot be made
// - 400 : if the request is incorrect
// - 501 : if the requested voting rule is not supported
func DoNewBallot(servUrl string, rule string, deadline string, votersID []string, tiebreak []comsoc.Alternative) (res rs.NewBallotResponse, err error) {

	if len(votersID) == 0 {
		return res, errors.New("0::Cannot create new ballot without any voters")
	}

	nbAlts := len(tiebreak)

	req := rs.NewBallotRequest{
		Rule:     rule,
		Deadline: deadline,
		Voters:   votersID,
		Alts:     nbAlts,
		TieBreak: tiebreak,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("1::%s", err.Error())
	}

	url := fmt.Sprintf("%s/new-ballot", servUrl)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return res, fmt.Errorf("2::%s", err.Error())
	}

	if resp.StatusCode != 201 {
		return res, fmt.Errorf("%d", resp.StatusCode)
	}

	decodeResponse[rs.NewBallotResponse](resp, &res)

	return
}

// DoVote allows the voter to vote to a specific server at the given servUrl.
// It returns the status of the response (200 if vote made, other else)
//
// If the request fails, an error is returned with a specific format : errType::errMessage.
// The errType can be one of the following :
// - 1   : if there's a problem when converting request to json
// - 2   : if the request cannot be made
// - 400 : if the request is incorrect
// - 403 : if the voter already voted
// - 501 : if the requested voting rule is not supported
// - 503 : if the deadline is passed
func (agent *RestVoterAgent) DoVote(ballot string) (res int, err error) {
	req := rs.VoteRequest{
		Agent:   fmt.Sprint(agent.ID),
		Ballot:  ballot,
		Prefs:   agent.Prefs,
		Options: agent.opts,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("1::%s", err.Error())
	}

	url := fmt.Sprintf("%s/vote", agent.url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return res, fmt.Errorf("2::%s", err.Error())
	}

	res = resp.StatusCode

	if res != 200 {
		return res, fmt.Errorf("%d", res)
	}

	return
}

// DoResult gets the result of the vote of a specific ballot in a given server.
//
// If the request fails, an error is returned with a specific format : errType::errMessage.
// The errType can be one of the following :
// - 1   : if there's a problem when converting request to json
// - 2   : if the request cannot be made
// - 404 : if the ballot does not exist
// - 425 : if the vote is not done yet
func DoResult(servUrl string, ballot string) (res rs.ResultResponse, err error) {
	req := rs.ResultRequest{
		Ballot: ballot,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("1::%s", err.Error())
	}

	url := fmt.Sprintf("%s/result", servUrl)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return res, fmt.Errorf("2::%s", err.Error())
	}

	if resp.StatusCode != 200 {
		return res, fmt.Errorf("%d", resp.StatusCode)
	}

	decodeResponse[rs.ResultResponse](resp, &res)

	return
}

/*
	Handler
*/

// Starts the voter to vote to the specified ballot.
func (agent *RestVoterAgent) Start(ballot string) {
	voteRes, err := agent.DoVote(ballot)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[%d] : vote success for agent %q to ballot %q\n", voteRes, agent.ID, ballot)

	// Wait for the results
	res, err := DoResult(agent.url, ballot)
	for err != nil {
		time.Sleep(10 * time.Second)
		res, err = DoResult(agent.url, ballot)
	}

	log.Printf("%d has won the vote on ballot %q !\n", res.Winner, ballot)
	switch agent.likeLevelOfAlt(res.Winner) {
	case -1:
		log.Printf("agent %q does not know where this result came from 🤨", agent.ID)
	case 0:
		log.Printf("agent %q is very happy with the result 😄", agent.ID)
	case 1:
		log.Printf("agent %q is quite happy with the result 🙂", agent.ID)
	case 2:
		log.Printf("agent %q is quite sad with the result 🙁", agent.ID)
	case 3:
		log.Printf("agent %q is very sad with the result 😢", agent.ID)
	}
}
