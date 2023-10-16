package ballotagent

import (
	"context"
	"errors"
	"fmt"
	"ia04/comsoc"
	"ia04/utils/concurrent"
	"log"
	"strings"
	"time"
)

type RestBallotAgent struct {
	ID       string
	Type     string
	Deadline time.Time
	Voters   []string
	NbAlts   int
	Ctx      context.Context
	scf      func(comsoc.Profile) (comsoc.Alternative, error)
}

func NewRestBallotAgent(id string,
	voteType string,
	deadline string,
	voters []string,
	nbAlts int,
	tieBreaks []comsoc.Alternative,
) (*RestBallotAgent, error) {
	deadlineTime, err := time.Parse(time.RFC3339, deadline)
	if err != nil {
		return nil, errors.New("400::error while parsing date : " + deadline)
	}

	// Check number of alternatives
	if len(tieBreaks) != nbAlts {
		return nil, errors.New("400::error not the same number of alternatives in nbAlts and tieBreaks")
	}

	// Check if tieBreaks has the right values
	valueCheck := make([]comsoc.Alternative, nbAlts)
	for _, tbv := range tieBreaks {
		if tbv < 1 || tbv > comsoc.Alternative(nbAlts) {
			return nil, errors.New("400::error one value in the tieBreaks is not between 1 and nbAlts")
		}
		if valueCheck[tbv-1] > 0 {
			return nil, errors.New("400::error theres two times the same value in tieBreaks")
		}

		valueCheck[tbv-1]++
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

	ctx := context.Background()

	return &RestBallotAgent{id, voteType, deadlineTime, voters, nbAlts, ctx, comsoc.SCFFactory(chosenSCF, comsoc.TieBreakFactory(tieBreaks))}, nil
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
	tnow := time.Now()

	if tnow.After(agent.Deadline) {
		return errors.New("the ballot deadline is set before now")
	}

	go func() {
		ctx, cancelCtx := context.WithDeadline(agent.Ctx, agent.Deadline)
		agent.Ctx = ctx
		defer cancelCtx()

		<-ctx.Done()
		log.Println(agent.ID, "Vote closed")
	}()

	return nil
}
