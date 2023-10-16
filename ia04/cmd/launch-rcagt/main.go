package main

import (
	"fmt"

	"ia04/agt/voteragent"
	"ia04/comsoc"
)

func main() {
	candidates := [...]comsoc.Alternative{1, 2, 3}
	ag := voteragent.NewRestClientAgent("id1", "http://localhost:8080", candidates[:])
	ag.Start()
	fmt.Scanln()
}
