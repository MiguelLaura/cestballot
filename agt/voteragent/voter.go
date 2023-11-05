package voteragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

type RestVoterAgent struct {
	agentId  string
	url      string
	ballotId string
	prefs    []comsoc.Alternative
	options  []int
}

func NewRestVoterAgent(agentId string, url string, ballotId string, prefs []comsoc.Alternative, options []int) *RestVoterAgent {
	return &RestVoterAgent{agentId, url, ballotId, prefs, options}
}

func (rva *RestVoterAgent) DoVote() (err error) {
	req := agt.RequestVoter{
		AgentId:  rva.agentId,
		BallotId: rva.ballotId,
		Prefs:    rva.prefs,
		Options:  rva.options,
	}

	// sérialisation de la requête
	url := rva.url + "/vote"
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

func (rva *RestVoterAgent) Start() {
	log.Printf("démarrage de %s", rva.agentId)
	err := rva.DoVote()

	if err != nil {
		log.Fatal(rva.agentId, " error:", err.Error())
	} else {
		log.Printf("[%s] %s %d %d\n", rva.agentId, rva.ballotId, rva.prefs, rva.options)
	}
}
