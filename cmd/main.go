package main

import (
	"ia04/agt"
	"ia04/agt/ballotagent"
	"ia04/agt/voteragent"
	"ia04/comsoc"
	"time"
)

func main() {
	// prefs := [][]comsoc.Alternative{
	// 	{1, 2, 3, 4},
	// 	{1, 2, 3, 4},
	// 	{1, 2, 3, 4},
	// 	{1, 2, 3, 4},
	// 	{1, 2, 3, 4},

	// 	{2, 3, 4, 1},
	// 	{2, 3, 4, 1},
	// 	{2, 3, 4, 1},
	// 	{2, 3, 4, 1},

	// 	{4, 3, 1, 2},
	// 	{4, 3, 1, 2},
	// 	{4, 3, 1, 2},
	// }

	// fmt.Println(comsoc.STV_SCF(prefs))

	voteChnl := make(chan agt.Vote, ballotagent.NB_DEFAULT_VOTERS)
	defer close(voteChnl)

	bureauVote := ballotagent.NewBallotAgent(1, "Bureau Valée", []comsoc.Alternative{1, 2, 3, 4}, voteChnl)

	voter1 := voteragent.NewVoterAgent(1, "Didier", []comsoc.Alternative{2, 4, 3, 1}, voteChnl)
	voter2 := voteragent.NewVoterAgent(2, "Julie", []comsoc.Alternative{1, 4, 3, 2}, voteChnl)
	voter3 := voteragent.NewVoterAgent(3, "Clara", []comsoc.Alternative{3, 2, 1, 4}, voteChnl)
	voter4 := voteragent.NewVoterAgent(4, "Pierre", []comsoc.Alternative{1, 2, 3, 4}, voteChnl)
	voter5 := voteragent.NewVoterAgent(5, "Lisa", []comsoc.Alternative{4, 2, 3, 1}, voteChnl)

	voter1.Start()
	voter2.Start()
	voter3.Start()
	voter4.Start()
	voter5.Start()

	bureauVote.Start()

	<-time.After(10 * time.Second)

}
