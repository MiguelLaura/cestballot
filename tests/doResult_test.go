package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func TestIncorrectBallot(t *testing.T) {
	req := agt.RequestResult{
		BallotId: "Inexistant",
	}
	url := servAddr + "/result"
	data, _ := json.Marshal(req)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("The error code should be %s but received [%d] %s", "[404] 404 Not Found", resp.StatusCode, resp.Status)
	}
}

func TestTooEarly(t *testing.T) {
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
		t.Fatal(err)
	}

	_, err = ag_b.DoResult()
	if err == nil {
		t.Fatal("The error code should be [425] 425 Too Early")
	} else if err.Error() != "[425] 425 Too Early" {
		t.Fatalf("The error code should be %s but received %s. Deadline not over", "[425] 425 Too Early", err.Error())
	}
}

func TestResult(t *testing.T) {
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
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	res, err := ag_b.DoResult()
	if err != nil {
		t.Fatalf("The request should work")
	}

	if int(res.Winner) != 2 {
		t.Fatalf("Incorrect result, winner should be 2 instead of %d", int(res.Winner))
	}
	correct := []comsoc.Alternative{2, 1, 4, 3}
	for idx, val := range res.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result, ranking should be [2, 1, 4, 3] instead of %d", res.Ranking)
		}
	}
}
