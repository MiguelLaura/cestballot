package ballotagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

type RestBallotAgent struct {
	BallotId string
	url      string
	rule     string
	deadline string
	voterIds []string
	alts     int
	tieBreak []comsoc.Alternative
}

func NewRestBallotAgent(url string, rule string, deadline string, voterIds []string, alts int, tieBreak []comsoc.Alternative) *RestBallotAgent {
	return &RestBallotAgent{url: url, rule: rule, deadline: deadline, voterIds: voterIds, alts: alts, tieBreak: tieBreak}
}

func (*RestBallotAgent) decodeResponseBallot(r *http.Response) (res agt.ResponseBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &res)
	return
}

func (rba *RestBallotAgent) DoNewBallot() (err error) {
	req := agt.RequestBallot{
		Rule:     rba.rule,
		Deadline: rba.deadline,
		VoterIds: rba.voterIds,
		Alts:     rba.alts,
		TieBreak: rba.tieBreak,
	}

	// sérialisation de la requête
	url := rba.url + "/new_ballot"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	res, _ := rba.decodeResponseBallot(resp)
	rba.BallotId = string(res.BallotId)
	return
}

func (rba *RestBallotAgent) Start() {
	log.Printf("démarrage du ballot")
	err := rba.DoNewBallot()

	if err != nil {
		log.Fatal("Ballot ", err.Error())
	} else {
		log.Printf("[%s] %s %s %s %d %d\n", rba.BallotId, rba.rule, rba.deadline, rba.voterIds, rba.alts, rba.tieBreak)
	}
}
