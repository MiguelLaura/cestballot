package main

import (
	"fmt"
	rst "gitlab.utc.fr/mennynat/ia04-tp/restserver"
)

func main() {
	server := rst.NewRestServerAgent("1", "localhost:8080")
	server.Start()
	fmt.Scanln()
}
