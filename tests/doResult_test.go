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
	req := agt.RequestResult{
		BallotId: ag_b.BallotId,
	}
	url := servAddr + "/result"
	data, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("The error code should be %s but received [%d] %s", "[200] 200 OK", resp.StatusCode, resp.Status)
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

	req := agt.RequestResult{
		BallotId: ag_b.BallotId,
	}
	url := servAddr + "/result"
	data, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooEarly {
		t.Fatalf("The error code should be %s but received [%d] %s. Deadline not over", "[425] 425 Too Earky", resp.StatusCode, resp.Status)
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
	req := agt.RequestResult{
		BallotId: ag_b.BallotId,
	}
	url := servAddr + "/result"
	data, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("The error code should be %s but received [%d] %s", "[200] 200 OK", resp.StatusCode, resp.Status)
	}

	req_rep := agt.ResponseResult{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err = json.Unmarshal(buf.Bytes(), &req_rep)
	if err != nil {
		t.Fatal(err)
	}

	if int(req_rep.Winner) != 2 {
		t.Fatalf("Incorrect result, winner should be 2 instead of %d", int(req_rep.Winner))
	}
	correct := []comsoc.Alternative{2, 1, 4, 3}
	for idx, val := range req_rep.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result, ranking should be [2, 1, 4, 3] instead of %d", req_rep.Ranking)
		}
	}
}
