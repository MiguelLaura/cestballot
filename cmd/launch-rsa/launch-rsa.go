/*
Lance un serveur REST qui gère les agents.

Utilisation :

	launch-rsa [flags]

Les flags peuvent être :

	-h, --host nomHôte
		Indique le nom de l'hôte.
		Défaut : localhost

	-p, --port numeroPort
		Indique le port du serveur.
		Défaut : 8080
*/
package main

import (
	"flag"
	"fmt"
	"log"

	"gitlab.utc.fr/mennynat/ia04-tp/agt"
)

func main() {

	// Traitement des flags

	var host string
	var port int

	flag.StringVar(&host, "host", "localhost", "Hôte du serveur")
	flag.StringVar(&host, "h", "localhost", "Hôte du serveur (raccourci)")

	flag.IntVar(&port, "port", 8080, "Port du serveur")
	flag.IntVar(&port, "p", 8080, "Port du serveur (raccourci)")

	flag.Parse()

	if port < 0 {
		log.Fatalf("Le numéro de port ne peut être négatif (donné %d)", port)
	}

	// Execution du script

	server := agt.NewRestServerAgent(fmt.Sprintf("%s:%d", host, port))
	server.Start()
	fmt.Scanln()
}
