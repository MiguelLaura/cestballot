package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func extractAlternativesFromString[T int | comsoc.Alternative](altStringList string) []T {
	strSplit := strings.Split(altStringList, ",")

	res := make([]T, len(strSplit))

	if len(altStringList) == 0 {
		return res
	}

	for idx, tbString := range strSplit {
		tbConv, err := strconv.Atoi(tbString)
		if err != nil {
			log.Fatal("The given alternatives does not contain int values")
		}
		res[idx] = T(tbConv)
	}
	return res
}

func main() {

	var wg sync.WaitGroup

	host := flag.String("host", "http://localhost", "url of the host")
	port := flag.Int("port", 8080, "port of the host")
	agentId := flag.String("id", "ag_id1", "id of the agent")
	agentName := flag.String("name", "ag_id1", "name of the agent")
	preferences := flag.String("prefs", "1,2,4,3", "preferences of the agent")
	opts := flag.String("opts", "", "opts of the agent for specific vote methods")
	ballot := flag.String("ballot", "vote0", "the ID of the ballot to which the voter will vote")

	flag.Parse()

	thePrefs := extractAlternativesFromString[comsoc.Alternative](*preferences)
	theOpts := extractAlternativesFromString[int](*opts)

	agent := voteragent.NewRestVoterAgent(
		*agentId,
		*agentName,
		thePrefs,
		theOpts,
		fmt.Sprintf("%s:%d", *host, *port),
	)

	wg.Add(1)
	go func() {
		agent.Start(*ballot)
		wg.Done()
	}()

	wg.Wait()
}
