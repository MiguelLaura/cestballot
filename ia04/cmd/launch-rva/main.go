package main

import (
	"fmt"
	"ia04/agt/voteragent"
	"ia04/comsoc"
)

func main() {
	ag := voteragent.NewRestVoterAgent("ag_id1", "http://localhost:8080", "scrutin12", []comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10}, nil)
	ag.Start()
	fmt.Scanln()
}
