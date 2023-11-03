// Package ballotagent contains an agent representing a ballot
package ballotagent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
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
	ctx        context.Context                 // The context in which the ballot runs
	tiebreak   []comsoc.Alternative            // The tiebreak used to determine which vote to choose when SCF cannot decide
	res        comsoc.Alternative              // The result of the vote; available once the deadline is passed
	ranking    []comsoc.Alternative            // The ranking (aka SWF) associated with the result)
}

// NewRestBallotAgent creates a new BallotAgent.
// The deadline is given as a string in the RFC3339 format and a list of alternative is given as a tiebreak
// to choose an alternative when the ballot cannot decide which one is better
//
// If the Agent was not created, an error is returned witch a specific format : errType::errMessage.
// The errType can be one of the follow :
// - 1 : incorrect date format
// - 2 : deadline in the past
// - 3 : not the correct number of voter
// - 4 : error in the alternatives and tiebreak
// - 5 : vote type not supported
func NewRestBallotAgent(
	id string,
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

	// Checks if there's at least one voters
	if len(voters) == 0 {
		return nil, errors.New("3::error no voter given")
	}

	// Checks the number of alternatives
	if nbAlts < 2 {
		return nil, errors.New("4::error less than 2 alternatives")
	}

	if len(tieBreaks) != nbAlts {
		return nil, errors.New("4::error not the same number of alternatives in nbAlts and tieBreaks")
	}

	// Check if tieBreaks has the right values (no duplicate and values between 1 and nbAlts)
	valueCheck := make([]comsoc.Alternative, nbAlts)
	for _, tbv := range tieBreaks {
		if tbv < 1 || tbv > comsoc.Alternative(nbAlts) {
			return nil, errors.New("4::error one value in the tieBreaks is not between 1 and nbAlts")
		}
		if valueCheck[tbv-1] > 0 {
			return nil, errors.New("4::error there's two times the same value in tieBreaks")
		}

		valueCheck[tbv-1]++
	}

	// Checks if the voting method is supported
	if !isVoteSupported(voteType) {
		return nil, errors.New("5::error " + voteType + " not supported")
	}

	theVoters := make(map[string][]comsoc.Alternative)
	theVotersOpts := make(map[string][]int)
	for _, voter := range voters {
		theVoters[voter] = nil
		theVotersOpts[voter] = nil
	}

	ctx := context.Background()

	return &RestBallotAgent{
		sync.Mutex{},
		id,
		voteType,
		deadlineTime,
		theVoters,
		theVotersOpts,
		nbAlts,
		ctx,
		tieBreaks,
		0,
		nil,
	}, nil
}

// isVoteSupported tell if the given voting method is supported.
func isVoteSupported(voteType string) bool {
	return slices.Contains(supportedVotingMethods[:], voteType)
}

// Strings gives a string representing the Ballot.
func (agent *RestBallotAgent) String() string {
	return fmt.Sprintf("Bureau %s (%s)", agent.ID, agent.VoteType)
}

// IsClosed indicates if the ballot is closed or not.
func (agent *RestBallotAgent) IsClosed() bool {
	select {
	case <-agent.ctx.Done():
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
		return false, fmt.Errorf("1::%q is closed", agent.String())
	}

	if _, exists := agent.Voters[voterId]; !exists {
		return false, fmt.Errorf("2::Agent %q cannot vote here", voterId)
	}

	if agent.Voters[voterId] != nil {
		return false, fmt.Errorf("3::Agent %q already voted", voterId)
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

// GetVoteResult returns the result of the vote and the associated ranking.
// If the vote is not made, the method returns 0.
//
// If the result cannot be returned, an error is given in a specific format : errType::errMessage.
// The errType can be one of the follow :
// - 1 : If the ballot is not closed
// - 2 : If the vote is not processed
func (agent *RestBallotAgent) GetVoteResult() (comsoc.Alternative, []comsoc.Alternative, error) {
	agent.Lock()
	defer agent.Unlock()

	if !agent.IsClosed() {
		return 0, nil, fmt.Errorf("1::%q is not closed", agent.String())
	}

	if !agent.isVoteDone() {
		return 0, nil, fmt.Errorf("2::vote not processed yet for %q", agent.String())
	}

	return agent.res, agent.ranking, nil
}

// isVoteDone indicates if the vote has been processed or not
func (agent *RestBallotAgent) isVoteDone() bool {
	return agent.res != 0
}

// processVote computes the vote.
func (agent *RestBallotAgent) processVote() (err error) {
	agent.Lock()
	defer agent.Unlock()

	voteProfile := make(comsoc.Profile, 0)
	voteOpts := make([][]int, 0)
	for voterId, prefs := range agent.Voters {
		if len(prefs) != 0 {
			voteProfile = append(voteProfile, prefs)
			voteOpts = append(voteOpts, agent.VotersOpts[voterId])
		}
	}

	// Gets the SCF corresponding to the ballot voting method
	var chosenSCF func(comsoc.Profile) ([]comsoc.Alternative, error) = nil
	var chosenSWF func(comsoc.Profile) (comsoc.Count, error) = nil
	switch agent.VoteType {
	case "majority":
		chosenSCF = comsoc.MajoritySCF
		chosenSWF = comsoc.MajoritySWF
	case "borda":
		chosenSCF = comsoc.BordaSCF
		chosenSWF = comsoc.BordaSWF
	case "stv":
		chosenSCF = comsoc.STV_SCF
		chosenSWF = comsoc.STV_SWF
	case "copeland":
		chosenSCF = comsoc.CopelandSCF
		chosenSWF = comsoc.CopelandSWF
	case "approval":

		thresholds := make([]int, len(voteOpts))
		for idx := range voteOpts {
			if len(voteOpts[idx]) == 0 {
				thresholds[idx] = 1
			} else {
				thresholds[idx] = voteOpts[idx][0]
			}
		}

		chosenSCF = func(p comsoc.Profile) ([]comsoc.Alternative, error) {
			return comsoc.ApprovalSCF(p, thresholds)
		}

		chosenSWF = func(p comsoc.Profile) (comsoc.Count, error) {
			return comsoc.ApprovalSWF(p, thresholds)
		}
	}

	if chosenSCF == nil || chosenSWF == nil {
		return fmt.Errorf("vote method %q not supported", agent.VoteType)
	}

	tiebreak := comsoc.TieBreakFactory(agent.tiebreak)
	scf := comsoc.SCFFactory(chosenSCF, tiebreak)
	swf := comsoc.SWFFactory(chosenSWF, tiebreak)

	resSCF, err := scf(voteProfile)

	if err != nil {
		log.Printf("Error while voting...")
		return
	}
	agent.res = resSCF

	resSWF, err := swf(voteProfile)

	if err == nil {
		agent.ranking = resSWF
	}

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
		ctx, cancelCtx := context.WithDeadline(agent.ctx, agent.Deadline)
		agent.ctx = ctx
		defer cancelCtx()

		// Wait until the Deadline is reached
		<-ctx.Done()
		log.Println(agent.ID, "Vote closed")

		err := agent.processVote()
		if err != nil {
			log.Println(agent.ID, err.Error())
		}
	}()

	return nil
}
