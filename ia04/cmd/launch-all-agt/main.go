package main

import (
	"fmt"
	"ia04/agt/ballotagent"
	"ia04/agt/voteragent"
	"ia04/comsoc"
	"log"
)

func main() {
	const n = 10
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	clAgts := make([]voteragent.RestVoterAgent, 0, n)
	servAgt := ballotagent.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients...")
	for i := 0; i < n; i++ {
		agentId := fmt.Sprintf("id%02d", i)
		// [A FAIRE]
		// Générer aléatoirement les prefs
		prefs := []comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10}
		agt := voteragent.NewRestVoterAgent(agentId, url2, "scrutin12", prefs, nil)
		clAgts = append(clAgts, *agt)
	}

	for _, agt := range clAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt voteragent.RestVoterAgent) {
			go agt.Start()
		}(agt)
	}

	fmt.Scanln()
}
