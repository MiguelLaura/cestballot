package voteragent

import (
	"fmt"
	"ia04/agt"
	"ia04/comsoc"
	"ia04/utils/concurrent"
	"log"
	"strconv"
	"strings"
)

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt comsoc.Alternative, prefs []comsoc.Alternative) int {
	for index, value := range prefs {
		if value == alt {
			return index
		}
	}

	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 comsoc.Alternative, prefs []comsoc.Alternative) bool {
	rk1, rk2 := rank(alt1, prefs), rank(alt2, prefs)
	return rk1 != -1 && rk2 != -1 && rk1 < rk2
}

type VoterAgent struct {
	agt.Agent
	c chan agt.Vote
}

func NewVoterAgent(id agt.AgentID, name string, prefs []comsoc.Alternative, chnlVot chan agt.Vote) *VoterAgent {
	agent := VoterAgent{}
	agent.ID = id
	agent.Name = name
	agent.Prefs = prefs
	agent.c = chnlVot

	return &agent
}

func (agent *VoterAgent) Equal(agent2 VoterAgent) bool {
	return agent.ID == agent2.ID && agent.Name == agent2.Name
}

func (agent *VoterAgent) DeepEqual(agent2 VoterAgent) bool {
	return agent.Equal(agent2) && concurrent.Equal(agent.Prefs, agent2.Prefs)
}

func (agent *VoterAgent) Clone() VoterAgent {
	return *NewVoterAgent(agent.ID, agent.Name, agent.Prefs, agent.c)
}

func (agent *VoterAgent) String() string {
	return fmt.Sprintf("%d : %s", agent.ID, agent.Name)
}

func (agent *VoterAgent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	return isPref(a, b, agent.Prefs)
}

func (agent *VoterAgent) Start() {
	go func() {
		cRep := make(chan string, 1)
		defer close(cRep)

		agent.c <- agt.Vote{Agt: agt.Agent{ID: agent.ID, Name: agent.Name, Prefs: agent.Prefs}, C: cRep}

		resVote := <-cRep
		theFinalVote, _ := strconv.Atoi(strings.Split(resVote, " ")[1])
		vote := comsoc.Alternative(theFinalVote)

		if vote == agent.Prefs[0] || agent.Prefers(vote, agent.Prefs[0]) {
			log.Printf("%s est content.e du vote !! 😃", agent.String())
		} else {
			log.Printf("%s est déçu.e du vote !! 😢", agent.String())
		}
	}()
}
