package ballotagent

import (
	"fmt"
	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
	"gitlab.utc.fr/mennynat/ia04-tp/utils/concurrent"
	"log"
	"time"
)

const NB_DEFAULT_VOTERS int = 100

type BallotAgent struct {
	ID        agt.AgentID
	Name      string
	VoteType  string
	voters    []agt.AgentID
	votes     comsoc.Profile
	responses []chan string
	scf       func(comsoc.Profile) (comsoc.Alternative, error)
	c         chan agt.Vote
}

func NewBallotAgent(id agt.AgentID, name string, allAlts []comsoc.Alternative, chnl chan agt.Vote) *BallotAgent {
	agent := BallotAgent{}
	agent.ID = id
	agent.Name = name
	agent.VoteType = "Majority"
	agent.voters = make([]agt.AgentID, 0, NB_DEFAULT_VOTERS)
	agent.votes = make(comsoc.Profile, 0, NB_DEFAULT_VOTERS)
	agent.responses = make([]chan string, 0, 100)
	agent.scf = comsoc.SCFFactory(comsoc.MajoritySCF, comsoc.TieBreakFactory(allAlts))
	agent.c = chnl

	return &agent
}

func (agent *BallotAgent) Equal(agent2 BallotAgent) bool {
	return agent.ID == agent2.ID && agent.Name == agent2.Name
}

func (agent *BallotAgent) DeepEqual(agent2 BallotAgent) bool {
	return agent.Equal(agent2) && concurrent.Equal(agent.voters, agent2.voters)
}

func (agent *BallotAgent) String() string {
	return fmt.Sprintf("Bureau %q", agent.Name)
}

func (agent *BallotAgent) Start() {
	go func() {
		timeout := time.After(8 * time.Second)
		loop := true
		for loop {
			select {
			case vote := <-agent.c:
				if idx, _ := concurrent.Find(agent.voters, func(id agt.AgentID) bool { return id == vote.Agt.ID }); idx == -1 {
					log.Printf("%s a voté !", vote.Agt.Name)
					agent.votes = append(agent.votes, vote.Agt.Prefs)
					agent.voters = append(agent.voters, vote.Agt.ID)
					agent.responses = append(agent.responses, vote.C)
				}
			case <-timeout:
				loop = false
			}
		}

		theVote, err := agent.scf(agent.votes)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s a pour resultat : %d\n", agent.String(), theVote)

		for _, repChnl := range agent.responses {
			repChnl <- fmt.Sprintf("%s %d", agent.VoteType, theVote)
		}
	}()
}
