package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func main() {
	var wg sync.WaitGroup
	const nBallot = 3
	const nVoter = 10
	const url1 = ":8080"
	const url2 = "http://localhost:8080"
	rules := [...]string{"majority", "borda", "approval", "condorcet", "copeland", "STV"}
	tieBreak := []comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10}
	alts := len(tieBreak)

	baAgts := make([]ballotagent.RestBallotAgent, 0, nBallot)
	voAgts := make([]voteragent.RestVoterAgent, 0, nVoter)
	servAgt := agt.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients ballot...")
	for i := 0; i < nBallot; i++ {
		wg.Add(1)
		rule := rules[rand.Intn(len(rules))]
		deadline := time.Now().Add(time.Duration(15 * time.Second)).Format(time.RFC3339)
		var voterIds [nVoter]string
		for j := 0; j < nVoter; j++ {
			voterIds[j] = "ag_id" + fmt.Sprint(j+1)
		}
		agt := ballotagent.NewRestBallotAgent(url2, rule, deadline, voterIds[:], alts, tieBreak)
		baAgts = append(baAgts, *agt)
	}

	for _, agt := range baAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		go func(agt ballotagent.RestBallotAgent) {
			agt.Start()
			defer wg.Done()
		}(agt)
	}
	wg.Wait()

	log.Println("démarrage des clients voter...")
	prefs := make([]comsoc.Alternative, len(tieBreak))
	copy(prefs, tieBreak)
	for i := 0; i < nVoter; i++ {
		wg.Add(1)
		agentId := fmt.Sprintf("id%02d", i+1)
		ballotId := "scrutin0" + fmt.Sprint(rand.Intn(nBallot)+1)
		rand.Shuffle(len(prefs), func(i, j int) { prefs[i], prefs[j] = prefs[j], prefs[i] })
		prefsAgt := make([]comsoc.Alternative, len(prefs))
		copy(prefsAgt, prefs)
		options := []int{rand.Intn(alts) + 1}
		agt := voteragent.NewRestVoterAgent(agentId, url2, ballotId, prefsAgt, options)
		voAgts = append(voAgts, *agt)
	}

	for _, agt := range voAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		go func(agt voteragent.RestVoterAgent) {
			agt.Start()
			defer wg.Done()
		}(agt)
	}
	wg.Wait()

	fmt.Scanln()
}
