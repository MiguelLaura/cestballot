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

	host := flag.String("host", "http://localhost", "url of the host")
	port := flag.Int("port", 8080, "port of the host")
	rule := flag.String("rule", "majority", "the voting rule")
	deadline := flag.String("deadline", time.Now().Add(10*time.Second).Format(time.RFC3339), "the deadline of the voting process")
	voters := flag.String("voters", "ag_id1,ag_id2,ag_id3", "list of all agents")
	tiebreak := flag.String("tiebreak", "4,2,3,1", "list of the tiebreak")

	flag.Parse()

	theVoters := strings.Split(*voters, ",")
	fmt.Println(theVoters)
	theTiebreak := extractAlternativesFromString(*tiebreak)

	res, err := voteragent.DoNewBallot(
		fmt.Sprintf("%s:%d", *host, *port),
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
