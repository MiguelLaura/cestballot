package restserver_test

/*
func TestVotesWorking(t *testing.T) {

	// Create a ballot
	rep, err := voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1, 2, 4, 3},
	)

	if err != nil {
		t.Fatal("doNewBallot should work; maybe the server is offline")
	}

	agt1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
		servAddr,
	)

	agt2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		"ag_id2",
		[]comsoc.Alternative{4, 1, 3, 2},
		nil,
		servAddr,
	)

	agt3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		"ag_id3",
		[]comsoc.Alternative{1, 2, 4, 3},
		nil,
		servAddr,
	)

	agts := []*voteragent.RestVoterAgent{agt1, agt2, agt3}

	for _, agt := range agts {
		go func(agt voteragent.RestVoterAgent) {
			agt.Start(rep.Id)
		}(*agt)
	}

	time.Sleep(5 * time.Second)

	resRep, err := voteragent.DoResult(servAddr, rep.Id)
	if err != nil {
		t.Fatal("doResult should work")
	}

	if resRep.Winner != 1 {
		t.Fatalf("The voting method should return 1; but returns %d", resRep.Winner)
	}

}

*/
