package main

import (
	"fmt"
	"ia04-tp/agt/voteragent"
	"ia04-tp/comsoc"
)

func main() {
	ag := voteragent.NewRestVoterAgent("ag_id1", "http://localhost:8080", "scrutin01", []comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10}, nil)
	ag.Start()
	fmt.Scanln()
}
