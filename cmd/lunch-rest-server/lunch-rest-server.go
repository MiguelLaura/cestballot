package main

import (
	"flag"
	"fmt"

	rst "gitlab.utc.fr/mennynat/ia04-tp/restserver"
)

func main() {

	host := flag.String("host", "localhost", "url of the host")
	port := flag.Int("port", 8080, "port of the host")

	flag.Parse()

	server := rst.NewRestServerAgent("1", fmt.Sprintf("%s:%d", *host, *port))
	server.Start()
	fmt.Scanln()
}
