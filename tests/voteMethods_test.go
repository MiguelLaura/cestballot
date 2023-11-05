package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func TestMajority(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{4, 2, 3, 1},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	ag_v1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	ag_v2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{4, 1, 3, 2},
		nil,
	)
	ag_v3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 4, 3},
		nil,
	)

	agts := []*voteragent.RestVoterAgent{ag_v1, ag_v2, ag_v3}

	for _, ag_v := range agts {
		wg.Add(1)
		go func(ag_v voteragent.RestVoterAgent) {
			_ = ag_v.DoVote()
			defer wg.Done()
		}(*ag_v)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if int(res.Winner) != 4 {
		t.Fatalf("Incorrect result, winner should be 4 instead of %d", int(res.Winner))
	}
	correct := []comsoc.Alternative{4, 2, 1, 3}
	for idx, val := range res.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result, ranking should be [4, 2, 1, 3] instead of %d", res.Ranking)
		}
	}
}

func TestBorda(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"borda",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		3,
		[]comsoc.Alternative{2, 3, 1},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	ag_v1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		nil,
	)
	ag_v2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		nil,
	)
	ag_v3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{3, 2, 1},
		nil,
	)

	agts := []*voteragent.RestVoterAgent{ag_v1, ag_v2, ag_v3}

	for _, ag_v := range agts {
		wg.Add(1)
		go func(ag_v voteragent.RestVoterAgent) {
			_ = ag_v.DoVote()
			defer wg.Done()
		}(*ag_v)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if int(res.Winner) != 1 {
		t.Fatalf("Incorrect result, winner should be 1 instead of %d", int(res.Winner))
	}
	correct := []comsoc.Alternative{1, 2, 3}
	for idx, val := range res.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result, ranking should be [1, 2, 3] instead of %d", res.Ranking)
		}
	}
}

func TestApproval(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"borda",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		3,
		[]comsoc.Alternative{1, 2, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	ag_v1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		[]int{2},
	)
	ag_v2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 3, 2},
		[]int{1},
	)
	ag_v3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 1},
		[]int{2},
	)

	agts := []*voteragent.RestVoterAgent{ag_v1, ag_v2, ag_v3}

	for _, ag_v := range agts {
		wg.Add(1)
		go func(ag_v voteragent.RestVoterAgent) {
			_ = ag_v.DoVote()
			defer wg.Done()
		}(*ag_v)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if int(res.Winner) != 1 {
		t.Fatalf("Incorrect result, winner should be 1 instead of %d", int(res.Winner))
	}
	correct := []comsoc.Alternative{1, 2, 3}
	for idx, val := range res.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result, ranking should be [1, 2, 3] instead of %d", res.Ranking)
		}
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

	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"stv",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		agtsNames,
		4,
		[]comsoc.Alternative{4, 2, 3, 1},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	for idx := range agtsNames {
		wg.Add(1)
		ag_v := voteragent.NewRestVoterAgent(
			agtsNames[idx],
			servAddr,
			ag_b.BallotId,
			prefs[idx],
			nil,
		)
		go func(ag_v voteragent.RestVoterAgent) {
			_ = ag_v.DoVote()
			defer wg.Done()
		}(*ag_v)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if int(res.Winner) != 1 {
		t.Fatalf("Incorrect result, winner should be 1 instead of %d", int(res.Winner))
	}
	correct := []comsoc.Alternative{1, 2, 4, 3}
	for idx, val := range res.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result, ranking should be [1, 2, 4, 3] instead of %d", res.Ranking)
		}
	}
}

func TestCondorcet(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"condorcet",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		3,
		[]comsoc.Alternative{1, 2, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	ag_v1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		nil,
	)
	ag_v2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		nil,
	)
	ag_v3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{3, 2, 1},
		nil,
	)

	agts := []*voteragent.RestVoterAgent{ag_v1, ag_v2, ag_v3}

	for _, ag_v := range agts {
		wg.Add(1)
		go func(ag_v voteragent.RestVoterAgent) {
			_ = ag_v.DoVote()
			defer wg.Done()
		}(*ag_v)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if res.Ranking != nil {
		t.Fatalf("No ranking with Condorcet instead of %d", res.Ranking)
	}
	if int(res.Winner) != 1 {
		t.Fatalf("Incorrect result, winner should be 1 instead of %d", int(res.Winner))
	}
}

func TestCondorcetNoWinner(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"condorcet",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		3,
		[]comsoc.Alternative{1, 2, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	ag_v1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		nil,
	)
	ag_v2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 1},
		nil,
	)
	ag_v3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{3, 1, 2},
		nil,
	)

	agts := []*voteragent.RestVoterAgent{ag_v1, ag_v2, ag_v3}

	for _, ag_v := range agts {
		wg.Add(1)
		go func(ag_v voteragent.RestVoterAgent) {
			_ = ag_v.DoVote()
			defer wg.Done()
		}(*ag_v)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if int(res.Winner) != 0 {
		t.Fatalf("Incorrect result, winner should be 0 instead of %d", int(res.Winner))
	}
	if res.Ranking != nil {
		t.Fatalf("No ranking with Condorcet instead of %d", res.Ranking)
	}
}

func TestCopeland(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"copeland",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		3,
		[]comsoc.Alternative{1, 2, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	ag_v1 := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		nil,
	)
	ag_v2 := voteragent.NewRestVoterAgent(
		"ag_id2",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3},
		nil,
	)
	ag_v3 := voteragent.NewRestVoterAgent(
		"ag_id3",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{3, 2, 1},
		nil,
	)

	agts := []*voteragent.RestVoterAgent{ag_v1, ag_v2, ag_v3}

	for _, ag_v := range agts {
		wg.Add(1)
		go func(ag_v voteragent.RestVoterAgent) {
			_ = ag_v.DoVote()
			defer wg.Done()
		}(*ag_v)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if int(res.Winner) != 1 {
		t.Fatalf("Incorrect result, winner should be 1 instead of %d", int(res.Winner))
	}
	correct := []comsoc.Alternative{1, 2, 3}
	for idx, val := range res.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result, ranking should be [1, 2, 3] instead of %d", res.Ranking)
		}
	}
}
