/*
lunch-rest-server starts the REST server to which the users connect to create Ballots and to vote

Usage :

	launch-rest-server [flags]

The flags are:

	--host hostName
		Specify the name of the host.
		Default : localhost
	--port portNumber
		Specify the number of the port to which the server listens.
		Default : 8080
	-v
		Starts the server in verbose mode.
*/
package main

import (
	"flag"
	"fmt"

	rst "gitlab.utc.fr/mennynat/ia04-tp/restserver"
)

func main() {

	host := flag.String("host", "localhost", "url of the host")
	port := flag.Int("port", 8080, "port of the host")
	verbose := flag.Bool("v", false, "is the server verbose")

	flag.Parse()

	server := rst.NewRestServerAgent("1", fmt.Sprintf("%s:%d", *host, *port))
	server.SetVerbose(*verbose)

	server.Start()
	fmt.Scanln()
}
