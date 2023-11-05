/*
Lance un nouveau agent votant.

Utilisation :

	launch-rva [flags]

Les flags peuvent être :

	-h, --host nomHôte
		Indique le nom de l'hôte.
		Défaut : localhost

	-p, --port numeroPort
		Indique le port du serveur.
		Défaut : 8080

	-b, --ballot nomBureauVote
		Indique le nom du bureau de vote auquel l'agent vote.
		Défaut : scrutin01

	-a, --agent nomAgent
		Indique le nom de l'agent votant.
		Défaut : ag_id1

	--prefs preferences
		Indique les preferences de l'agent.
		Format : alt1,alt2,alt3
		Défaut : 4,2,3,5,9,8,7,1,6,11,12,10

	-o, --opts listeOptions
		Donne une liste optionnelle d'options pour les méthodes de vote.
		Format : opt1,opt2,opt3
		Défaut : nil
*/
package main

import (
	"flag"
	"fmt"
	"log"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/voteragent"
	"gitlab.utc.fr/mennynat/ia04-tp/cmd"
)

func main() {

	// Traitement des flags

	var host, ballot, agent string
	var port int
	var prefs cmd.AltFlag
	var opts cmd.OptFlag

	flag.StringVar(&host, "host", "localhost", "Hôte du serveur")
	flag.StringVar(&host, "h", "localhost", "Hôte du serveur (raccourci)")

	flag.IntVar(&port, "port", 8080, "Port du serveur")
	flag.IntVar(&port, "p", 8080, "Port du serveur (raccourci)")

	flag.StringVar(&ballot, "ballot", "scrutin01", "Nom du bureau de vote où l'agent vote")
	flag.StringVar(&ballot, "b", "scrutin01", "Nom du bureau de vote où l'agent vote (raccourci)")

	flag.StringVar(&agent, "agent", "ag_id1", "Nom de l'agent votant")
	flag.StringVar(&agent, "a", "ag_id1", "Nom de l'agent votant (raccourci)")

	flag.Var(&prefs, "prefs", "Préférences de l'agent")

	flag.Var(&opts, "opts", "Options de vote de l'agent")
	flag.Var(&opts, "o", "Options de vote de l'agent (raccourci)")

	flag.Parse()

	if port < 0 {
		log.Fatalf("Le numéro de port ne peut être négatif (donné %d)", port)
	}

	// Execution du script

	ag := voteragent.NewRestVoterAgent(
		agent,
		fmt.Sprintf("http://%s:%d", host, port),
		ballot,
		prefs.GetAlts(),
		opts.GetOpts(),
	)
	ag.Start()
	fmt.Scanln()
}
