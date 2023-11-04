/*
create-ballot creates a new ballot.

Usage :

	create-ballot [flags]

The flags are:

	--host hostName
		Specify the url of the host.
		Default : localhost
	--port portNumber
		Specify the number of the port to which the server listens.
		Default : 8080
	--rule ruleName
		Specify the name of the rule to use to vote.
		Can be one of the values : [majority, borda, stv, copeland, approval]
		Default : majority
	--deadline deadline
		Specify the deadline in RFC3339 format.
		Default : current time + 1 minute
	--voters listOfVoters
		Specify all the voters authorized to vote on this ballot.
		Format : voter1,voter2,voter3,...
		Default : ag_id1,ag_id2,ag_id3
	--tiebreak listOfAlternatives
		Specify the tiebreak to use when voting.
		Format : alt1,alt2,alt3,...
		Default : 4,2,3,1
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func extractAlternativesFromString(altStringList string) []comsoc.Alternative {
	strSplit := strings.Split(altStringList, ",")

	res := make([]comsoc.Alternative, len(strSplit))

	for idx, tbString := range strSplit {
		tbConv, err := strconv.Atoi(tbString)
		if err != nil {
			log.Fatal("The given alternatives does not contain int values")
		}
		res[idx] = comsoc.Alternative(tbConv)
	}
	return res
}

func main() {

	host := flag.String("host", "localhost", "url of the host")
	port := flag.Int("port", 8080, "port of the host")
	rule := flag.String("rule", "majority", "the voting rule")
	deadline := flag.String("deadline", time.Now().Add(5*time.Second).Format(time.RFC3339), "the deadline of the voting process")
	voters := flag.String("voters", "ag_id1,ag_id2,ag_id3", "list of all agents")
	tiebreak := flag.String("tiebreak", "4,2,3,1", "list of the tiebreak")

	flag.Parse()

	theVoters := strings.Split(*voters, ",")
	fmt.Println(theVoters)
	theTiebreak := extractAlternativesFromString(*tiebreak)

	res, err := voteragent.DoNewBallot(
		fmt.Sprintf("http://%s:%d", *host, *port),
		*rule,
		*deadline,
		theVoters,
		theTiebreak,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Ballot %q successfully created", res.Id)

}
