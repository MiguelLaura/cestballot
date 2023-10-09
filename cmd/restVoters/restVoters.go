package main

import (
	"fmt"
	rst "ia04/restserver"
)

func main() {
	server := rst.NewRestServerAgent("1", "localhost:8080")
	server.Start()
	fmt.Scanln()
}
