package ballotagent

import (
	"context"
	"errors"
	"fmt"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
	"log"
	"strings"
	"time"
)

type RestBallotAgent struct {
	ID         string
	Type       string
	Deadline   time.Time
	Voters     map[string][]comsoc.Alternative
	VotersOpts map[string][]int
	NbAlts     int
	Ctx        context.Context
	scf        func(comsoc.Profile) (comsoc.Alternative, error)
	res        comsoc.Alternative
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

	theVoters := make(map[string][]comsoc.Alternative)
	for _, voter := range voters {
		theVoters[voter] = nil
	}

	ctx := context.Background()

	return &RestBallotAgent{
		id,
		voteType,
		deadlineTime,
		theVoters,
		make(map[string][]int),
		nbAlts,
		ctx,
		comsoc.SCFFactory(chosenSCF, comsoc.TieBreakFactory(tieBreaks)),
		0,
	}, nil
}

func (agent *RestBallotAgent) String() string {
	return fmt.Sprintf("Bureau %s (%s)", agent.ID, agent.Type)
}

func (agent *RestBallotAgent) IsClosed() bool {
	select {
	case <-agent.Ctx.Done():
		if agent.res == 0 {
			return false
		}
		return true
	default:
		return false
	}
}

func (agent *RestBallotAgent) Vote(voterId string, prefs []comsoc.Alternative, opts []int) (bool, error) {
	if agent.IsClosed() {
		return false, errors.New(fmt.Sprintf("1::%q is closed", agent.String()))
	}

	if _, exists := agent.Voters[voterId]; !exists {
		return false, errors.New(fmt.Sprintf("2::Agent %q cannot vote here", voterId))
	}

	if agent.Voters[voterId] != nil {
		return false, errors.New(fmt.Sprintf("3::Agent %q already voted", voterId))
	}

	allAlts := make([]comsoc.Alternative, agent.NbAlts)
	for idx := range allAlts {
		allAlts[idx] = comsoc.Alternative(idx + 1)
	}

	if comsoc.CheckProfile(prefs, allAlts) != nil {
		return false, errors.New("4::Preferences not correct")
	}

	agent.Voters[voterId] = prefs
	if opts != nil {
		agent.VotersOpts[voterId] = opts
	}

	return true, nil
}

func (agent *RestBallotAgent) GetVoteResult() (comsoc.Alternative, error) {
	if !agent.IsClosed() {
		return 0, errors.New(fmt.Sprintf("1::%q is not closed", agent.String()))
	}

	return agent.res, nil
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

		// Proceed to vote...
		voteProfile := make(comsoc.Profile, len(agent.Voters))
		idx := 0
		for _, prefs := range agent.Voters {
			voteProfile[idx] = prefs
			idx++
		}

		resVote, err := agent.scf(voteProfile)
		if err != nil {
			log.Printf("Error while voting...")
			return
		}
		agent.res = resVote
	}()

	return nil
}
