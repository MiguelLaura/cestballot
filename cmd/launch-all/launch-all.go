/*
launch-all starts a REST server and creates a given number of ballots where a
a given number of voters can vote randomly.

Usage :

	launch-all [flags]

The flags are:

	--host hostName
		Specify the name of the host.
		Default : localhost
	--port portNumber
		Specify the number of the port to which the server listens.
		Default : 8080
	--nbBallots nbBallots
		Specify the number of ballots created
		Default : 2
	--nbVoters nbVoters
		Specify the number of voters that can vote to the ballots
		Default : 20
	--nbAlts nbAlts
		Specify the number of alternatives in the ballots
		Default : 15
	-v
		Starts the server in verbose mode.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"ia04-tp/agt/voteragent"
	"ia04-tp/comsoc"
	"ia04-tp/restserver"
)

func main() {

	host := flag.String("host", "localhost", "url of the host")
	port := flag.Int("port", 8080, "port of the host")
	nbBallots := flag.Int("nbBallot", 2, "number of ballots created")
	nbVoters := flag.Int("nbVoters", 20, "total number of voters")
	nbAlts := flag.Int("nbAlts", 15, "total number of alternatives")
	verbose := flag.Bool("v", false, "is the server verbose")

	flag.Parse()

	// Init
	var servURL = fmt.Sprintf("%s:%d", *host, *port)
	var rules = [...]string{"majority", "borda", "stv", "copeland", "approval"}
	wg := sync.WaitGroup{}

	ballots := make([]string, *nbBallots)

	probNoVote := 1.0 / 10.0
	voters := make([]string, *nbVoters)
	for voterIdx := range voters {
		voters[voterIdx] = "ag_id" + fmt.Sprint(voterIdx+1)
	}

	tiebreak := computePreferences(*nbAlts)

	// REST Server
	log.Println("-----------------------")
	log.Println("Starting REST server...")
	log.Println("-----------------------")
	serv := restserver.NewRestServerAgent("serv", servURL)
	go serv.Start()

	time.Sleep(time.Second) // Waits a bit to be sure the server starts

	// Ballot
	log.Println("-------------------")
	log.Println("Starting ballots...")
	log.Println("-------------------")
	for idxBallot := 0; idxBallot < *nbBallots; idxBallot++ {
		rule := rules[rand.Intn(len(rules))]
		ballot, err := voteragent.DoNewBallot(
			"http://"+servURL,
			rule,
			time.Now().Add(10*time.Second).Format(time.RFC3339),
			voters,
			tiebreak,
		)
		if err != nil {
			log.Fatal(err)
		}

		ballots[idxBallot] = ballot.Id
		say(*verbose, "Ballot %q created with rule %s and tiebreak %v\n", ballot.Id, rule, tiebreak)
	}

	// Voters
	log.Println("------------------")
	log.Println("Starting voters...")
	log.Println("------------------")
	for idxVoter := 0; idxVoter < *nbVoters; idxVoter++ {

		if rand.Float64() < probNoVote {
			say(*verbose, "Voter %q does not want to vote\n", voters[idxVoter])
			continue
		}

		wg.Add(1)
		go func(idxVoter int) {
			defer wg.Done()

			prefs := computePreferences(*nbAlts)

			voter := voteragent.NewRestVoterAgent(
				voters[idxVoter],
				voters[idxVoter],
				prefs,
				[]int{rand.Intn(*nbAlts) + 1},
				"http://"+servURL,
			)

			ballot := ballots[rand.Intn(len(ballots))]

			_, err := voter.DoVote(ballot)
			if err != nil {
				log.Fatal(err)
			}

			say(*verbose, "Voter %q vote on ballot %q with prefs %v\n", voters[idxVoter], ballot, prefs)
		}(idxVoter)
	}

	wg.Wait()
	time.Sleep(10 * time.Second)

	// Get the results
	log.Println("--------------")
	log.Println("Get results...")
	log.Println("--------------")
	for idxBallot := 0; idxBallot < *nbBallots; idxBallot++ {

		res, err := voteragent.DoResult("http://"+servURL, ballots[idxBallot])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Ballot %q : result %d - %v\n", ballots[idxBallot], res.Winner, res.Ranking)
	}

	fmt.Scanln()
}

func computePreferences(nbPrefs int) []comsoc.Alternative {
	res := make([]comsoc.Alternative, nbPrefs)
	for altIdx := range res {
		res[altIdx] = comsoc.Alternative(altIdx + 1)
	}
	rand.Shuffle(len(res), func(i, j int) { res[i], res[j] = res[j], res[i] })
	return res
}

func say(verbose bool, format string, replace ...any) {
	if verbose {
		log.Printf(format, replace...)
	}
}
