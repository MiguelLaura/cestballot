package restserver_test

import (
	"fmt"
	"testing"
	"time"

	"ia04-tp/agt/voteragent"
	"ia04-tp/comsoc"
)

func TestMajority(t *testing.T) {

	// Create a ballot
	rep, err := voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{4, 2, 3, 1},
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

	if resRep.Winner != 4 {
		t.Fatalf("The voting method should return 4; but returns %d", resRep.Winner)
	}
}

func TestBorda(t *testing.T) {

	// Create a ballot
	rep, err := voteragent.DoNewBallot(
		servAddr,
		"borda",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{2, 3, 1},
	)

	if err != nil {
		t.Fatal("doNewBallot should work; maybe the server is offline")
	}

	agt1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{1, 2, 3},
		nil,
		servAddr,
	)

	agt2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		"ag_id2",
		[]comsoc.Alternative{1, 2, 3},
		nil,
		servAddr,
	)

	agt3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		"ag_id3",
		[]comsoc.Alternative{3, 2, 1},
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

func TestApproval(t *testing.T) {

	// Create a ballot
	rep, err := voteragent.DoNewBallot(
		servAddr,
		"approval",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1, 2, 3},
	)

	if err != nil {
		t.Fatal("doNewBallot should work; maybe the server is offline")
	}

	agt1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{1, 2, 3},
		[]int{2},
		servAddr,
	)

	agt2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		"ag_id2",
		[]comsoc.Alternative{1, 3, 2},
		[]int{1},
		servAddr,
	)

	agt3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		"ag_id3",
		[]comsoc.Alternative{2, 3, 1},
		[]int{2},
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

func TestSTV(t *testing.T) {

	prefs := comsoc.Profile{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},

		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},

		{4, 3, 1, 2},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
	}

	agtsNames := make([]string, len(prefs))
	for idx := range agtsNames {
		agtsNames[idx] = fmt.Sprintf("ag_id%d", idx+1)
	}

	// Create a ballot
	rep, err := voteragent.DoNewBallot(
		servAddr,
		"stv",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		agtsNames,
		[]comsoc.Alternative{4, 2, 3, 1},
	)

	if err != nil {
		t.Fatal("doNewBallot should work; maybe the server is offline")
	}

	for idx := range agtsNames {

		agt := voteragent.NewRestVoterAgent(
			agtsNames[idx],
			agtsNames[idx],
			prefs[idx],
			nil,
			servAddr,
		)

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
