# Application de vote

Auteur·rices : Laura Miguel, Nathan Menny   
Groupe : TD1-I

Application de vote en GO réalisée dans le cadre des cours d'IA04.

# Hypothèses
* On a au moins 1 voter
* On a au moins 2 candidats/alternatives
* Quand il n'y a pas de SWF (Condorcet), on renvoie nil/null
* Pas d'erreur quand on donne des options alors que la règle n'est pas approval
* Erreur quand on donne en option autre chose qu'un seul nombre dans une slice et qu'on va utiliser l'option
* Les ag_id peuvent s'abstenir
* Quand tous les agt s'abstiennent, on renvoie le résultat à partir de tiebreak
* Avec condorcet, on peut ne pas avoir de gagnant (si pas de voter ayant voté, le gagnant est le préféré dans tieBreak)
* Si le winner est 0, il n'y a pas de gagnant (par exemple, pour Condorcet)
* Un gagnant de Condorcet gagne tous ses matchs, il ne peut y avoir d'égalité.

# Utilisation

Les méthodes de vote suivantes ont été implémentées :
* Condorcet
* Borda
* Majority
* Approval
* STV
* Copeland

## Scripts

Les différents scripts suivants se trouvent dans cmd/ :
* [launch-rsa](#launch-rsa) : lance le serveur
* [launch-rba](#launch-rba) : lance un ballot
* [launch-rva](#launch-rva) : lance un voter et vote
* [launch-all-agt](#launch-all-agt) : lance le serveur, des ballots et des voters et attend pour le résultat (les ballotes ont tous les mêmes candidats, tieBreak et voter mais ont des règles et deadlines différentes; les voter votent dans un seul ballot décidé aléatoirement)

## Détail des commandes

### launch-rsa

```

	launch-rsa [flags]

Les flags peuvent être :

	-h, --host nomHôte
		Indique le nom de l'hôte.
		Défaut : localhost

	-p, --port numeroPort
		Indique le port du serveur.
		Défaut : 8080

```

### launch-rba

```

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
		Donne la deadline après laquelle le bureau de vote ferme.
		Format : RFC3339
		Défaut : temps actuel + 5 secondes

	-a, --agents listeAgents
		Donne la liste des agents autorisés à voter.
		Format : id1,id2,id3
		Défaut : ag_id1

	-t, --tiebreak tiebreak
		Indique le tiebreak à utiliser lors des votes.
		Format : alt1,alt2,alt3
		Défaut : 4,2,3,5,9,8,7,1,6,11,12,10

```

### launch-rva

```

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

```

### launch-all-agt

```

	launch-all-agt [flags]

Les flags peuvent être :

	-h, --host nomHôte
		Indique le nom de l'hôte.
		Défaut : localhost

	-p, --port numeroPort
		Indique le port du serveur.
		Défaut : 8080

	-b, --n-ballots nombreBallot
		Indique le nombre de bureaux de vote à créer.
		Défaut : 3

	-v, --n-voters nombreVotants
		Indique le nombre de voters à créer.
		Défaut : 10

	-t, --tiebreak tiebreak
		Indique le tiebreak à utiliser lors des votes.
		Format : alt1,alt2,alt3
		Défaut : 4,2,3,5,9,8,7,1,6,11,12,10

```