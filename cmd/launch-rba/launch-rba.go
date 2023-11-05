package main

import (
	"fmt"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func main() {
	ag := ballotagent.NewRestBallotAgent(
		"http://localhost:8080",
		"majority",
		time.Now().Add(5*time.Second).Format(time.RFC3339),
		[]string{"ag_id1"},
		12,
		[]comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10},
	)
	ag.Start()
	fmt.Scanln()
}
