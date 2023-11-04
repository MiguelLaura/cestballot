package main

import (
	"fmt"
	"ia04-tp/agt/ballotagent"
	"ia04-tp/comsoc"
	"time"
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
