/*
Lance un nouveau bureau de vote.

Utilisation :

	launch-rba [flags]

Les flags peuvent être :

	-h, --host nomHôte
		Indique le nom de l'hôte.
		Défaut : localhost

	-p, --port numeroPort
		Indique le port du serveur.
		Défaut : 8080

	--rule méthodeDeVote
		Méthode de vote utilisée dans le nouveau bureau de vote.
		Défaut : majority

	-d, --deadline deadline
		Donne la deadline après laquelle le bureau de vote ferme
		Format : RFC3339
		Défaut : temps actuel + 5 secondes

	-a, --agents liste des agents voters
		Donne la liste des agents autorisés à voter.
		Format : id1,id2,id3
		Défaut : ag_id1

	-t, --tiebreak tiebreak
		Indique le tiebreak à utiliser lors des votes.
		Format : alt1,alt2,alt3
		Défaut : 4,2,3,5,9,8,7,1,6,11,12,10
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gitlab.utc.fr/mennynat/ia04-tp/agt/ballotagent"
	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func main() {

	// Traitement des flags

	var host, rule, deadline string
	var port int
	var voters VotersFlag
	var alts AltFlag

	flag.StringVar(&host, "host", "localhost", "Hôte du serveur")
	flag.StringVar(&host, "h", "localhost", "Hôte du serveur (raccourci)")

	flag.IntVar(&port, "port", 8080, "Port du serveur")
	flag.IntVar(&port, "p", 8080, "Port du serveur (raccourci)")

	flag.StringVar(&rule, "rule", "majority", "Règle de vote utilisée")

	flag.StringVar(&deadline, "deadline", time.Now().Add(5*time.Second).Format(time.RFC3339), "Deadline au format RFC3339")
	flag.StringVar(&deadline, "d", time.Now().Add(5*time.Second).Format(time.RFC3339), "Deadline au format RFC3339 (raccourci)")

	flag.Var(&voters, "agents", "Liste des agents autorisés à voter")
	flag.Var(&voters, "a", "Liste des agents autorisés à voter (raccourci)")

	flag.Var(&alts, "tiebreak", "Tiebreak utilisée dans le vote")
	flag.Var(&alts, "t", "Tiebreak utilisée dans le vote")

	flag.Parse()

	if port < 0 {
		log.Fatalf("Le numéro de port ne peut être négatif (donné %d)", port)
	}

	// Execution du script

	ag := ballotagent.NewRestBallotAgent(
		fmt.Sprintf("http://%s:%d", host, port),
		rule,
		deadline,
		voters.GetVoters(),
		len(alts.GetAlts()),
		alts.GetAlts(),
	)
	ag.Start()
	fmt.Scanln()
}

// -----------------------------
// 	  Structures utilitaires
// -----------------------------

// Permet d'acquérir les votants en ligne de commande

type VotersFlag struct {
	voters []string
}

func (vf *VotersFlag) String() string {
	return strings.Join(vf.voters, ",")
}

func (vf *VotersFlag) Set(s string) error {
	if s[len(s)-1] == ',' {
		log.Fatalf("Format du flag des voters incorrect")
	}
	vf.voters = strings.Split(s, ",")
	return nil
}

func (vf *VotersFlag) GetVoters() []string {
	if vf.voters == nil {
		return []string{"ag_id1"}
	}
	return vf.voters
}

// Permet d'acquérir les alternatives en ligne de commande

type AltFlag struct {
	alternatives []comsoc.Alternative
}

func (vf *AltFlag) String() string {
	return fmt.Sprintf("%#v", vf.alternatives)
}

func (vf *AltFlag) Set(s string) error {
	altsStr := strings.Split(s, ",")
	alts := make([]comsoc.Alternative, len(altsStr))

	for altIdx, altStr := range altsStr {
		altConv, err := strconv.Atoi(altStr)

		if err != nil {
			log.Fatal("Une des alternative donnée n'est pas un entier")
		}

		alts[altIdx] = comsoc.Alternative(altConv)
	}

	vf.alternatives = alts
	return nil
}

func (vf *AltFlag) GetAlts() []comsoc.Alternative {
	if vf.alternatives == nil {
		return []comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10}
	}
	return vf.alternatives
}
