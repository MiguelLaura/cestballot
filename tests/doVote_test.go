package test

import (
	"testing"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func TestAgentVoter(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	err = ag_v.DoVote()
	if err != nil {
		t.Fatalf("The request should work")
	}
}

func TestAgentIncorrectBallot(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		"scrutin2000",
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. The ballot doesn't exist but it causes no problem", "[400] 400 Bad Request", err.Error())
	}
}

func TestAgentIncorrectID(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. The agent can't vote in this ballot", "[400] 400 Bad Request", err.Error())
	}

	ag_v = voteragent.NewRestVoterAgent(
		"agent_id329",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. The agent can't vote in this ballot", "[400] 400 Bad Request", err.Error())
	}
}

func TestAgentVoterAlreadyVoted(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	err = ag_v.DoVote()
	if err != nil {
		t.Fatalf("The request should work")
	}

	ag_v = voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [403] 403 Forbidden")
	} else if err.Error() != "[403] 403 Forbidden" {
		t.Fatalf("The error code should be %s but received %s. Voter already voted", "[403] 403 Forbidden", err.Error())
	}
}

func TestIncorrectAmountOfPrefs(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. Missing prefs", "[400] 400 Bad Request", err.Error())
	}
}

func TestNegativePrefs(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{-1, 2, 3, 4},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. Negative prefs", "[400] 400 Bad Request", err.Error())
	}
}

func TestNullPrefs(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{0, 1, 2, 3},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. Null prefs", "[400] 400 Bad Request", err.Error())
	}
}

func TestIncorrectPrefs(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{1, 2, 3, 5},
		nil,
	)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. Incorrect prefs", "[400] 400 Bad Request", err.Error())
	}
}

func TestDeadlineOver(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag_b.DoNewBallot()
	if err != nil {
		t.Fatal(err)
	}

	ag_v := voteragent.NewRestVoterAgent(
		"ag_id1",
		servAddr,
		ag_b.BallotId,
		[]comsoc.Alternative{2, 3, 4, 1},
		nil,
	)
	time.Sleep(1 * time.Second)
	err = ag_v.DoVote()
	if err == nil {
		t.Fatal("The error code should be [503] 503 Service Unavailable")
	} else if err.Error() != "[503] 503 Service Unavailable" {
		t.Fatalf("The error code should be %s but received %s. Deadline over", "[503] 503 Service Unavailable", err.Error())
	}
}
