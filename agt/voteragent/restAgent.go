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

type RestVoterAgent struct {
	agt.Agent
}

func NewRestVoterAgent(id agt.AgentID, name string, prefs []comsoc.Alternative) *RestVoterAgent {
	return &RestVoterAgent{agt.Agent{ID: id, Name: name, Prefs: prefs}}
}

/*
	Basic Agent methods (from AgentI interface)
*/

func (agent *RestVoterAgent) Equal(agent2 RestVoterAgent) bool {
	return agent.ID == agent2.ID && agent.Name == agent2.Name
}

func (agent *RestVoterAgent) DeepEqual(agent2 RestVoterAgent) bool {
	return agent.Equal(agent2) && concurrent.Equal(agent.Prefs, agent2.Prefs)
}

func (agent *RestVoterAgent) Clone() RestVoterAgent {
	return *NewRestVoterAgent(agent.ID, agent.Name, agent.Prefs)
}

func (agent *RestVoterAgent) String() string {
	return fmt.Sprintf("%d : %s", agent.ID, agent.Name)
}

func (agent *RestVoterAgent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	return isPref(a, b, agent.Prefs)
}

/*
	HTTP response decoder
*/

func decodeResponse[T any](r *http.Response, req *T) (err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

/*
	HTTP requesters
*/

func DoNewBallot(url string, rule string, deadline string, voters []RestVoterAgent, tieBreak []comsoc.Alternative) (res rs.NewBallotResponse, err error) {

	if len(voters) == 0 {
		return res, errors.New("0:Cannot create new ballot without any voters")
	}

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
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	if resp.StatusCode != 201 {
		return res, fmt.Errorf("%d:%s", resp.StatusCode, resp.Status)
	}

	decodeResponse[rs.NewBallotResponse](resp, &res)

	return
}

func (agent *RestVoterAgent) DoVote(url string, ballot string, opts []int) (res int, err error) {
	req := rs.VoteRequest{
		Agent:   fmt.Sprint(agent.ID),
		Ballot:  ballot,
		Prefs:   agent.Prefs,
		Options: opts,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	res = resp.StatusCode

	if res != 200 {
		return res, fmt.Errorf("%d:%s", res, resp.Status)
	}

	return
}

func (agent *RestVoterAgent) DoResult(url string, ballot string) (res rs.ResultResponse, err error) {
	req := rs.ResultRequest{
		Ballot: ballot,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return res, fmt.Errorf("%d:%s", resp.StatusCode, resp.Status)
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
