package restserver_test

import (
	"fmt"
	"testing"
	"time"

	"ia04-tp/agt/voteragent"
	"ia04-tp/comsoc"
	"ia04-tp/restserver"
)

func TestIncorrectBallot(t *testing.T) {

	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
		servAddr,
	)

	res, err := agt.DoVote(rep.Id)
	if err != nil {
		t.Fatal(err)
	}

	if res != restserver.VOTE_TAKEN {
		t.Fatal("The vote should be taken")
	}

	_, err = voteragent.DoResult(servAddr, "aBallot")

	if err.Error() != fmt.Sprintf("%d", restserver.NOT_FOUND) {
		t.Fatal("The ballot is supposed not to be found")
	}
}

func TestTooEarly(t *testing.T) {

	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
		servAddr,
	)

	res, err := agt.DoVote(rep.Id)
	if err != nil {
		t.Fatal(err)
	}

	if res != restserver.VOTE_TAKEN {
		t.Fatal("The vote should be taken")
	}

	_, err = voteragent.DoResult(servAddr, rep.Id)

	if err.Error() != fmt.Sprintf("%d", restserver.TOO_EARLY) {
		t.Fatal("Response is too early")
	}
}

func TestResult(t *testing.T) {

	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
		servAddr,
	)

	res, err := agt.DoVote(rep.Id)
	if err != nil {
		t.Fatal(err)
	}

	if res != restserver.VOTE_TAKEN {
		t.Fatalf("The vote should be taken")
	}

	time.Sleep(3 * time.Second)

	resVote, err := voteragent.DoResult(servAddr, rep.Id)
	if err != nil {
		t.Fatal(err)
	}

	if resVote.Winner != 2 {
		t.Fatalf("Incorrect result")
	}
}
