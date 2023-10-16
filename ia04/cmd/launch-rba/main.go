package main

import (
	"fmt"
	"ia04/agt/ballotagent"
)

func main() {
	server := ballotagent.NewRestBallotAgent(":8080")
	server.Start()
	fmt.Scanln()
}
