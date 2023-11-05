package restserver_test

import (
	"fmt"
	"testing"
	"time"

	"ia04-tp/agt/voteragent"
	"ia04-tp/comsoc"
	"ia04-tp/restserver"
)

const servAddr string = "http://localhost:8080"

/*
-------------------------------------------------------------------

	Start the REST server before running the following functions

-------------------------------------------------------------------
*/

func TestCorrectVotingMethod(t *testing.T) {

	_, err := voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1, 2, 4, 3},
	)

	if err != nil {
		t.Fatalf("The request should work")
	}
}

func TestIncorrectVotingMethod(t *testing.T) {

	_, err := voteragent.DoNewBallot(
		servAddr,
		"copernic",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1, 2, 4, 3},
	)

	if err.Error() != fmt.Sprintf("%d", restserver.NOT_IMPL) {
		t.Log("Voting method not implemented")
		t.Fatalf("The error code should be %d but received %s", restserver.NOT_IMPL, err.Error())
	}
}

func TestWrongDeadline(t *testing.T) {

	_, err := voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1, 2, 4, 3},
	)

	if err.Error() != fmt.Sprintf("%d", restserver.BAD_REQUEST) {
		t.Log("Date in the past")
		t.Fatalf("The error code should be %d but received %s", restserver.BAD_REQUEST, err.Error())
	}
}

func TestNotEnoughAgent(t *testing.T) {

	_, err := voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{},
		[]comsoc.Alternative{1, 2, 4, 3},
	)

	if err.Error() != fmt.Sprintf("%d", restserver.BAD_REQUEST) {
		t.Log("Not enough voters")
		t.Fatalf("The error code should be %d but received %s", restserver.BAD_REQUEST, err.Error())
	}
}

func TestIncorrectTiebreak(t *testing.T) {

	_, err := voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1},
	)

	if err.Error() != fmt.Sprintf("%d", restserver.BAD_REQUEST) {
		t.Log("Not enough alternatives")
		t.Fatalf("The error code should be %d but received %s", restserver.BAD_REQUEST, err.Error())
	}

	_, err = voteragent.DoNewBallot(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1", "ag_id2", "ag_id3"},
		[]comsoc.Alternative{1, 2, 3, 3},
	)

	if err.Error() != fmt.Sprintf("%d", restserver.BAD_REQUEST) {
		t.Log("Duplicate in the tiebreak")
		t.Fatalf("The error code should be %d but received %s", restserver.BAD_REQUEST, err.Error())
	}
}
