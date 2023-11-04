package ballotagent

import (
	"bytes"
	"encoding/json"
	"ia04-tp/agt"
	"ia04-tp/agt/ballotagent"
	"ia04-tp/agt/voteragent"
	"ia04-tp/comsoc"
	"net/http"
	"testing"
	"time"
)

func TestIncorrectBallot(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
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

	time.Sleep(5 * time.Second)
	req := agt.RequestResult{
		BallotId: ag_b.BallotId,
	}

	// sérialisation de la requête
	url := servAddr + "/result"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("[%d] %s", resp.StatusCode, resp.Status)
	}
}

func TestTooEarly(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
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

	// sérialisation de la requête
	url := servAddr + "/result"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooEarly {
		t.Fatalf("[%d] %s", resp.StatusCode, resp.Status)
	}
}

func TestResult(t *testing.T) {
	ag_b := ballotagent.NewRestBallotAgent(
		servAddr,
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
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

	time.Sleep(5 * time.Second)
	req := agt.RequestResult{
		BallotId: ag_b.BallotId,
	}

	// sérialisation de la requête
	url := servAddr + "/result"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("[%d] %s", resp.StatusCode, resp.Status)
	}

	req_rep := agt.ResponseResult{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err = json.Unmarshal(buf.Bytes(), &req_rep)
	if err != nil {
		t.Fatal(err)
	}

	if int(req_rep.Winner) != 2 {
		t.Fatalf("Incorrect result")
	}
	correct := []comsoc.Alternative{2, 1, 4, 3}
	for idx, val := range req_rep.Ranking {
		if correct[idx] != val {
			t.Fatalf("Incorrect result")
		}
	}
}
