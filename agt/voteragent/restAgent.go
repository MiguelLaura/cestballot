package voteragent

import (
	"fmt"

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

func (agent *RestVoterAgent) Start() {
	go func() {

	}()
}
