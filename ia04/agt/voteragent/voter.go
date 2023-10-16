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
	id    string
	url   string
	prefs []comsoc.Alternative
}

func NewRestClientAgent(id string, url string, prefs []comsoc.Alternative) *RestClientAgent {
	return &RestClientAgent{id, url, prefs}
}

func (rca *RestClientAgent) treatResponse(r *http.Response) comsoc.Alternative {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp agt.Response
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Winner
}

func (rca *RestClientAgent) doRequest() (res comsoc.Alternative, err error) {
	req := agt.RequestVoter{
		Prefs: rca.prefs,
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
	res = rca.treatResponse(resp)

	return
}

func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s", rca.id)
	res, err := rca.doRequest()

	if err != nil {
		log.Fatal(rca.id, "error:", err.Error())
	} else {
		log.Printf("[%s] %d = %d\n", rca.id, rca.prefs, res)
	}
}
