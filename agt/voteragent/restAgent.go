// Package voteragent contains an agent that can vote.

package voteragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
	rs "gitlab.utc.fr/mennynat/ia04-tp/restserver"
	"gitlab.utc.fr/mennynat/ia04-tp/utils/concurrent"
)

// The voting agent structure.
type RestVoterAgent struct {
	agt.Agent
	url string // The url of the server where the voter votes
}

// NewRestVoterAgent creates a new voting agent.
func NewRestVoterAgent(id agt.AgentID, name string, prefs []comsoc.Alternative, url string) *RestVoterAgent {
	return &RestVoterAgent{agt.Agent{ID: id, Name: name, Prefs: prefs}, url}
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
	return *NewRestVoterAgent(agent.ID, agent.Name, agent.Prefs, agent.url)
}

// String gives a string representation of the voting agent.
func (agent *RestVoterAgent) String() string {
	return fmt.Sprintf("%d : %s", agent.ID, agent.Name)
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
func DoNewBallot(servUrl string, rule string, deadline string, voters []RestVoterAgent, tieBreak []comsoc.Alternative) (res rs.NewBallotResponse, err error) {

	if len(voters) == 0 {
		return res, errors.New("0::Cannot create new ballot without any voters")
	}

	// Creates a slice with every voter's ID
	votersIDs := make([]string, len(voters))
	for voterIdx, currVoter := range voters {
		votersIDs[voterIdx] = fmt.Sprint(currVoter.ID)
	}

	req := rs.NewBallotRequest{
		Rule:     rule,
		Deadline: deadline,
		Voters:   votersIDs,
		Alts:     int(slices.Max(voters[0].Prefs)),
		TieBreak: tieBreak,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("1::%s", err.Error())
	}

	url := fmt.Sprintf("%s/new_ballot", servUrl)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return res, fmt.Errorf("2::%s", err.Error())
	}

	if resp.StatusCode != 201 {
		return res, fmt.Errorf("%d::%s", resp.StatusCode, resp.Status)
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
func (agent *RestVoterAgent) DoVote(ballot string, opts []int) (res int, err error) {
	req := rs.VoteRequest{
		Agent:   fmt.Sprint(agent.ID),
		Ballot:  ballot,
		Prefs:   agent.Prefs,
		Options: opts,
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
		return res, fmt.Errorf("%d::%s", res, resp.Status)
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
		return res, fmt.Errorf("%d::%s", resp.StatusCode, resp.Status)
	}

	decodeResponse[rs.ResultResponse](resp, &res)

	return
}

/*
	Handler
*/

func (agent *RestVoterAgent) Start() {
	go func() {

	}()
}
