package voteragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ia04/agt"
	"ia04/comsoc"
	"log"
	"net/http"
)

type RestClientAgent struct {
	agentId  string
	url      string
	ballotId string
	prefs    []comsoc.Alternative
	options  []int
}

func NewRestClientAgent(agentId string, url string, ballotId string, prefs []comsoc.Alternative, options []int) *RestClientAgent {
	return &RestClientAgent{agentId, url, ballotId, prefs, options}
}

func (rca *RestClientAgent) doVote() (err error) {
	req := agt.RequestVoter{
		AgentId:  rca.agentId,
		BallotId: rca.ballotId,
		Prefs:    rca.prefs,
		Options:  rca.options,
	}

	// sérialisation de la requête
	url := rca.url + "/vote"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	return
}

func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s", rca.agentId)
	err := rca.doVote()

	if err != nil {
		log.Fatal(rca.agentId, "error:", err.Error())
	} else {
		log.Printf("[%s] %s %d %d\n", rca.agentId, rca.ballotId, rca.prefs, rca.options)
	}
}
