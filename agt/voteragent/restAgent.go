package voteragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
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

func decodeRequest[T any](r *http.Request, req *T) (err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

/*
	HTTP requesters
*/

func (voter *RestVoterAgent) doNewBallot(rule string, deadline string, voters []RestVoterAgent, tieBreak []comsoc.Alternative) {

}

/*
	Handler
*/

func (agent *RestVoterAgent) Start() {
	go func() {

	}()
}
