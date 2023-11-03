package restserver_test

import (
	"testing"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
	"gitlab.utc.fr/mennynat/ia04-tp/restserver"
)

func createBallot() (restserver.NewBallotResponse, error) {
	return voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Add(3*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1, 2, 4, 3},
	)
}

func TestAgentVoter(t *testing.T) {

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
}

func TestAgentIncorrectBallot(t *testing.T) {

	_, err := createBallot()
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

	res, _ := agt.DoVote("aBallot")

	if res != restserver.BAD_REQUEST {
		t.Fatalf("The ballot doesn't exist but it causes no problem")
	}

}

func TestAgentIncorrectID(t *testing.T) {

	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"agent",
		"ag_id1",
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
		servAddr,
	)

	res, _ := agt.DoVote(rep.Id)

	if res != restserver.BAD_REQUEST {
		t.Fatal("The agent cannot vote here but it causes no problem")
	}
}

func TestAgentVoterAlreadyVoted(t *testing.T) {

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

	res, _ = agt.DoVote(rep.Id)

	if res != restserver.VOTE_ALREADY_DONE {
		t.Fatalf("The vote should already be done")
	}
}

func TestIncorrectAmountOfPrefs(t *testing.T) {

	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{2, 3},
		nil,
		servAddr,
	)

	res, _ := agt.DoVote(rep.Id)

	if res != restserver.BAD_REQUEST {
		t.Fatalf("The agent does not have the right amount of preferences")
	}

}

func TestNegativePrefs(t *testing.T) {
	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{-1, 2, 3, 4},
		nil,
		servAddr,
	)

	res, _ := agt.DoVote(rep.Id)

	if res != restserver.BAD_REQUEST {
		t.Fatalf("The agent does not have the right preferences")
	}

}

func TestNullPrefs(t *testing.T) {
	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{0, 2, 3, 4},
		nil,
		servAddr,
	)

	res, _ := agt.DoVote(rep.Id)

	if res != restserver.BAD_REQUEST {
		t.Fatalf("The agent does not have the right preferences")
	}

}

func TestIncorrectPrefs(t *testing.T) {
	rep, err := createBallot()
	if err != nil {
		t.Fatal(err)
	}

	agt := voteragent.NewRestVoterAgent(
		"ag_id1",
		"ag_id1",
		[]comsoc.Alternative{1, 2, 3, 5},
		nil,
		servAddr,
	)

	res, _ := agt.DoVote(rep.Id)

	if res != restserver.BAD_REQUEST {
		t.Fatalf("The agent does not have the right preferences")
	}
}

func TestDeadlineOver(t *testing.T) {
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

	time.Sleep(3 * time.Second)

	res, _ := agt.DoVote(rep.Id)

	if res != restserver.DEADLINE_OVER {
		t.Fatalf("The deadline is over")
	}
}
