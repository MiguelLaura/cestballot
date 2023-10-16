// Manque le choix de la méthode de vote, l'init de ballot
// + random pour prefs des voters
package main

import (
	"fmt"
	"log"

	"ia04/agt/ballotagent"
	"ia04/agt/voteragent"
	"ia04/comsoc"
)

func main() {
	const n = 3
	const url1 = ":8080"
	const url2 = "http://localhost:8080"
	candidates := [...]comsoc.Alternative{1, 2, 3}

	clAgts := make([]voteragent.RestClientAgent, 0, n)
	servAgt := ballotagent.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("id%02d", i)
		agt := voteragent.NewRestClientAgent(id, url2, candidates[:])
		clAgts = append(clAgts, *agt)
	}

	for _, agt := range clAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt voteragent.RestClientAgent) {
			go agt.Start()
		}(agt)
	}

	fmt.Scanln()
}
