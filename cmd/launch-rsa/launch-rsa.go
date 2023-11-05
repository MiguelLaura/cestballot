package main

import (
	"fmt"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
)

func main() {
	server := agt.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
