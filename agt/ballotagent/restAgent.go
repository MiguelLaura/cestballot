// Package ballotagent contains an agent representing a ballot

package ballotagent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

// The voting methods supported by the ballot.
// Must be updated when a new voting method is added
var supportedVotingMethods = [...]string{"majority", "borda", "stv", "copeland", "approval"}

// The BallotAgent handling the complete voting process
type RestBallotAgent struct {
	sync.Mutex                                 // Contains a mutex to avoid concurrency problems
	ID         string                          // The ID of the agent, must be unique
	VoteType   string                          // The type of vote done
	Deadline   time.Time                       // The deadline after which the ballot is closed
	Voters     map[string][]comsoc.Alternative // All the voters and their preferences
	VotersOpts map[string][]int                // The optional voting parameters per agent, depends on the vote type
	NbAlts     int                             // Number of alternatives available
	Ctx        context.Context                 // The context in which the ballot runs
	tiebreak   []comsoc.Alternative            // The tiebreak used to determine which vote to choose when SCF cannot decide
	res        comsoc.Alternative              // The result of the vote; available once the deadline is passed
}

// NewRestBallotAgent creates a new BallotAgent.
// The deadline is given as a string in the RFC3339 format and a list of alternative is given as a tiebreak
// to choose an alternative when the ballot cannot decide which one is better
//
// If the Agent was not created, an error is returned witch a specific format : errType::errMessage.
// The errType can be one of the follow :
// - 1 : incorrect date format
// - 2 : deadline in the past
// - 3 : error in the alternatives and tiebreak
// - 4 : vote type not supported
func NewRestBallotAgent(id string,
	voteType string,
	deadline string,
	voters []string,
	nbAlts int,
	tieBreaks []comsoc.Alternative,
) (*RestBallotAgent, error) {

	// Converts the deadline in a time
	deadlineTime, err := time.Parse(time.RFC3339, deadline)
	if err != nil {
		return nil, errors.New("1::error while parsing date : " + deadline)
	}

	// Checks if the deadline is not in the past
	if deadlineTime.Before(time.Now()) {
		return nil, errors.New("2::error " + deadline + " is in the past")
	}

	// Checks the number of alternatives
	if len(tieBreaks) != nbAlts {
		return nil, errors.New("3::error not the same number of alternatives in nbAlts and tieBreaks")
	}

	// Check if tieBreaks has the right values (no duplicate and values between 1 and nbAlts)
	valueCheck := make([]comsoc.Alternative, nbAlts)
	for _, tbv := range tieBreaks {
		if tbv < 1 || tbv > comsoc.Alternative(nbAlts) {
			return nil, errors.New("3::error one value in the tieBreaks is not between 1 and nbAlts")
		}
		if valueCheck[tbv-1] > 0 {
			return nil, errors.New("3::error theres two times the same value in tieBreaks")
		}

		valueCheck[tbv-1]++
	}

	// Checks if the voting method is supported
	if !isVoteSupported(voteType) {
		return nil, errors.New("4::" + voteType + " not supported")
	}

	theVoters := make(map[string][]comsoc.Alternative)

	ctx := context.Background()

	return &RestBallotAgent{
		sync.Mutex{},
		id,
		voteType,
		deadlineTime,
		theVoters,
		make(map[string][]int),
		nbAlts,
		ctx,
		tieBreaks,
		0,
	}, nil
}

// isVoteSupported tell if the given voting method is supported.
func isVoteSupported(voteType string) bool {
	voteTypeIdx := 0
	for voteTypeIdx < len(supportedVotingMethods) && supportedVotingMethods[voteTypeIdx] != voteType {
		voteTypeIdx++
	}
	return voteTypeIdx < len(supportedVotingMethods)
}

// Strings gives a string representing the Ballot.
func (agent *RestBallotAgent) String() string {
	return fmt.Sprintf("Bureau %s (%s)", agent.ID, agent.VoteType)
}

// IsClosed indicates if the ballot is closed or not.
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

// Vote allows a voter to vote in the ballot.
//
// If the vote has not succeed, an error is returned in a specific format : errType::errMessage.
// The errType can be one of the follow :
// - 1 : If the ballot is closed
// - 2 : If the agent cannot vote in this ballot
// - 3 : If the agent already vote
// - 4 : If the preferences are incorrect
func (agent *RestBallotAgent) Vote(voterId string, prefs []comsoc.Alternative, opts []int) (bool, error) {
	agent.Lock()
	defer agent.Unlock()

	if agent.IsClosed() {
		return false, errors.New(fmt.Sprintf("1::%q is closed", agent.String()))
	}

	if _, exists := agent.Voters[voterId]; !exists {
		return false, errors.New(fmt.Sprintf("2::Agent %q cannot vote here", voterId))
	}

	if agent.Voters[voterId] != nil {
		return false, errors.New(fmt.Sprintf("3::Agent %q already voted", voterId))
	}

	// Checks if the preferences are correct (no duplicate, all the values between 1 and agent.nbAlts)
	allAlts := make([]comsoc.Alternative, agent.NbAlts)
	for idx := range allAlts {
		allAlts[idx] = comsoc.Alternative(idx + 1)
	}

	if comsoc.CheckProfile(prefs, allAlts) != nil {
		return false, errors.New("4::Preferences not correct")
	}

	// Registers the vote
	agent.Voters[voterId] = prefs
	if opts != nil {
		agent.VotersOpts[voterId] = opts
	}

	return true, nil
}

// GetVoteResult returns the result of the vote.
//
// If the result cannot be returned, an error is given in a specific format : errType::errMessage.
// The errType can be one of the follow :
// - 1 : If the ballot is not closed
func (agent *RestBallotAgent) GetVoteResult() (comsoc.Alternative, error) {
	agent.Lock()
	defer agent.Unlock()

	if !agent.IsClosed() {
		return 0, errors.New(fmt.Sprintf("1::%q is not closed", agent.String()))
	}

	return agent.res, nil
}

// processVote computes the vote.
func (agent *RestBallotAgent) processVote() (err error) {
	agent.Lock()
	defer agent.Unlock()

	voteProfile := make(comsoc.Profile, len(agent.Voters))
	voteOpts := make([][]int, len(agent.VotersOpts))
	idx := 0
	for voterId, prefs := range agent.Voters {
		voteProfile[idx] = prefs
		voteOpts[idx] = agent.VotersOpts[voterId]
		idx++
	}

	// Gets the SCF corresponding to the ballot voting method
	var chosenSCF func(comsoc.Profile) ([]comsoc.Alternative, error) = nil
	switch agent.VoteType {
	case "majority":
		chosenSCF = comsoc.MajoritySCF
	case "borda":
		chosenSCF = comsoc.BordaSCF
	case "stv":
		chosenSCF = comsoc.STV_SCF
	case "copeland":
		chosenSCF = comsoc.CopelandSCF
	case "approval":
		chosenSCF = func(p comsoc.Profile) ([]comsoc.Alternative, error) {
			thresholds := make([]int, len(voteOpts))
			for idx := range voteOpts {
				thresholds[idx] = voteOpts[idx][0]
			}

			return comsoc.ApprovalSCF(p, thresholds)
		}
	}

	if chosenSCF == nil {
		return fmt.Errorf("vote method %q not supported", agent.VoteType)
	}

	scf := comsoc.SCFFactory(chosenSCF, comsoc.TieBreakFactory(agent.tiebreak))

	resVote, err := scf(voteProfile)

	if err != nil {
		log.Printf("Error while voting...")
		return
	}
	agent.res = resVote

	return nil
}

// Start opens the votes to the ballot
//
// If the deadline is in the past, an error is returned.
func (agent *RestBallotAgent) Start() error {
	if agent.Deadline.Before(time.Now()) {
		return errors.New("the ballot deadline is set before now")
	}

	go func() {
		ctx, cancelCtx := context.WithDeadline(agent.Ctx, agent.Deadline)
		agent.Ctx = ctx
		defer cancelCtx()

		// Wait until the Deadline is reached
		<-ctx.Done()
		log.Println(agent.ID, "Vote closed")

		err := agent.processVote()
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	return nil
}
