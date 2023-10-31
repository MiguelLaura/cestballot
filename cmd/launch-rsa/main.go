package main

import (
	"fmt"
	"ia04-tp/agt"
)

func main() {
	server := agt.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
