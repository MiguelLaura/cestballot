/*
Crée un serveur REST et tout un ensemble de bureau de votes et d'agents qui vont voter dessus.

Utilisation :

	launch-all-agt [flags]

Les flags peuvent être :

	-h, --host nomHôte
		Indique le nom de l'hôte.
		Défaut : localhost

	-p, --port numeroPort
		Indique le port du serveur.
		Défaut : 8080

	-b, --n-ballots nombreBallot
		Indique le nombre de bureaux de vote à créer.
		Défaut : 3

	-v, --n-voters nombreVotants
		Indique le nombre de voters à créer.
		Défaut : 10

	-t, --tiebreak tiebreak
		Indique le tiebreak à utiliser lors des votes.
		Format : alt1,alt2,alt3
		Défaut : 4,2,3,5,9,8,7,1,6,11,12,10
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/cmd"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func main() {

	// Traitement des flags

	var host string
	var port, nBallot, nVoter int
	var tbFlag cmd.AltFlag

	flag.StringVar(&host, "host", "localhost", "Hôte du serveur")
	flag.StringVar(&host, "h", "localhost", "Hôte du serveur (raccourci)")

	flag.IntVar(&port, "port", 8080, "Port du serveur")
	flag.IntVar(&port, "p", 8080, "Port du serveur (raccourci)")

	flag.IntVar(&nBallot, "n-ballots", 3, "Nombre de bureaux de votes à créer")
	flag.IntVar(&nBallot, "b", 3, "Nombre de bureaux de votes à créer (raccourci)")

	flag.IntVar(&nVoter, "n-voters", 10, "Nombre d'agent votants à créer")
	flag.IntVar(&nVoter, "v", 10, "Nombre d'agent votants à créer (raccourci)")

	flag.Var(&tbFlag, "tiebreak", "Tiebreak utilisée dans les votes")
	flag.Var(&tbFlag, "t", "Tiebreak utilisée dans les votes")

	flag.Parse()

	if port < 0 {
		log.Fatalf("Le numéro de port ne peut être négatif (donné %d)", port)
	}

	if nVoter < 1 {
		log.Fatalf("Ne peut pas avoir moins de 1 voter (donné %d)", nVoter)
	}

	if nBallot < 1 {
		log.Fatalf("Ne peut pas avoir moins de 1 bureau de vote (donné %d)", nBallot)
	}

	// Execution du script

	var wg sync.WaitGroup
	var url1 = fmt.Sprintf("%s:%d", host, port)
	var url2 = fmt.Sprintf("http://%s:%d", host, port)
	rules := [...]string{"majority", "borda", "approval", "condorcet", "copeland", "STV"}
	tieBreak := tbFlag.GetAlts()
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
		voterIds := make([]string, nVoter)
		for j := 0; j < nVoter; j++ {
			voterIds[j] = "ag_id" + fmt.Sprint(j+1)
		}
		agt := ballotagent.NewRestBallotAgent(url2, rule, deadline, voterIds, alts, tieBreak)
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
