package ballotagent

import (
	"errors"
	"fmt"
	"ia04/comsoc"
	"ia04/utils/concurrent"
	"strings"
	"time"
)

type RestBallotAgent struct {
	ID       string
	Type     string
	Deadline time.Time
	Voters   []string
	NbAlts   int
	opened   bool
	scf      func(comsoc.Profile) (comsoc.Alternative, error)
}

func NewRestBallotAgent(id string,
	voteType string,
	deadline string,
	voters []string,
	nbAlts int,
) (*RestBallotAgent, error) {
	deadlineTime, err := time.Parse(time.UnixDate, deadline)
	if err != nil {
		return nil, errors.New("400::error while parsing date : " + deadline)
	}

	allAlts := make([]comsoc.Alternative, nbAlts)
	for idx := range allAlts {
		allAlts[idx] = comsoc.Alternative(idx + 1)
	}

	var chosenSCF func(comsoc.Profile) ([]comsoc.Alternative, error) = nil

	switch strings.ToLower(voteType) {
	case "majority":
		chosenSCF = comsoc.MajoritySCF
	case "borda":
		chosenSCF = comsoc.BordaSCF
	}

	if chosenSCF == nil {
		return nil, errors.New("501::" + voteType + " not supported")
	}

	return &RestBallotAgent{id, voteType, deadlineTime, voters, nbAlts, false, comsoc.SCFFactory(chosenSCF, comsoc.TieBreakFactory(allAlts))}, nil
}

func (agent *RestBallotAgent) Equal(agent2 RestBallotAgent) bool {
	return agent.ID == agent2.ID && agent.Type == agent2.Type && agent.Deadline == agent2.Deadline && agent.NbAlts == agent2.NbAlts
}

func (agent *RestBallotAgent) DeepEqual(agent2 RestBallotAgent) bool {
	return agent.Equal(agent2) && concurrent.Equal(agent.Voters, agent2.Voters)
}

func (agent *RestBallotAgent) String() string {
	return fmt.Sprintf("Bureau %s (%s)", agent.ID, agent.Type)
}

func (agent *RestBallotAgent) Start() error {
	timeout := agent.Deadline.Sub(time.Now())
	if timeout < 0 {
		return errors.New("the ballot deadline is set before now")
	}

	go func() {
		agent.opened = true

		<-time.After(timeout)

		agent.opened = false
	}()

	return nil
}
