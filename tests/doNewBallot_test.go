package test

import (
	"ia04-tp/agt/ballotagent"
	"ia04-tp/comsoc"
	"testing"
	"time"
)

const servAddr string = "http://localhost:8080"

/*
-------------------------------------------------------------------

	Start the REST server before running the following functions

-------------------------------------------------------------------
*/

func TestCorrectVotingMethod(t *testing.T) {
	ag := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag.DoNewBallot()
	if err != nil {
		t.Fatalf("The request should work")
	}
}

func TestIncorrectVotingMethod(t *testing.T) {
	ag := ballotagent.NewRestBallotAgent(
		servAddr,
		"copernic",
		time.Now().Add(1*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag.DoNewBallot()
	if err == nil {
		t.Fatal("The error code should be [501] 501 Not Implemented")
	} else if err.Error() != "[501] 501 Not Implemented" {
		t.Fatalf("The error code should be %s but received %s", "[501] 501 Not Implemented", err.Error())
	}
}

func TestWrongDeadline(t *testing.T) {
	ag := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag.DoNewBallot()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. The deadline is in the past", "[400] 400 Bad Request", err.Error())
	}
}

func TestNotEnoughAgent(t *testing.T) {
	ag := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Format(time.RFC3339),
		[]string{},
		4,
		[]comsoc.Alternative{1, 2, 4, 3},
	)
	err := ag.DoNewBallot()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. Missing voter agent", "[400] 400 Bad Request", err.Error())
	}
}

func TestIncorrectTiebreak(t *testing.T) {
	ag := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		1,
		[]comsoc.Alternative{1},
	)
	err := ag.DoNewBallot()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. Missing candidates", "[400] 400 Bad Request", err.Error())
	}

	ag = ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		4,
		[]comsoc.Alternative{1, 2, 3, 3},
	)
	err = ag.DoNewBallot()
	if err == nil {
		t.Fatal("The error code should be [400] 400 Bad Request")
	} else if err.Error() != "[400] 400 Bad Request" {
		t.Fatalf("The error code should be %s but received %s. Incorrect tieBreak", "[400] 400 Bad Request", err.Error())
	}
}
