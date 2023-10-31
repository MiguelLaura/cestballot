package main

import (
	"fmt"
	"ia04-tp/agt/ballotagent"
	"ia04-tp/comsoc"
)

func main() {
	ag := ballotagent.NewRestBallotAgent("http://localhost:8080", "majority", "2023-10-31T11:12:08+01:00", []string{"ag_id1"}, 12, []comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10})
	ag.Start()
	fmt.Scanln()
}
