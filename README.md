# Application de vote

Auteur·rices : Nathan Menny, Laura Miguel   
Groupe : TD1-I

Application de vote en GO réalisée dans le cadre des cours d'IA04.

# Hypothèses

* On a au moins 2 candidats/alternatives
* On a au moins 1 voteur déclaré auprès du ballot
* Les voteurs peuvent s'abstenir
* Quand tous les agents s'abstiennent, le résultat est déterminé entièrement par le tieBreak
* Il n'y a pas d'erreur quand on donne des options alors que la règle n'est pas `approval`
* Il y a une erreur quand on utilise la méthode `approval` et qu'on donne en option autre chose qu'un seul nombre dans un slice
* Quand il n'y a pas de SWF (`condorcet`), on ne renvoie rien pour l'attribut ranking du résultat
* Avec `condorcet`, on peut ne pas avoir de gagnant, dans ce cas l'attribut gagnant du résultat est à 0 (si aucun voteur inscrit auprès du ballot n'a voté, le gagnant est le préféré dans tieBreak)
* Un gagnant de `condorcet` gagne tous ses matchs, il ne peut y avoir d'égalité

# Utilisation

Les méthodes de vote suivantes ont été implémentées :
* Condorcet
* Borda
* Majority
* Approval
* STV
* Copeland

## Commandes

Il y a 4 commandes se trouvent dans `cmd/` :
* [launch-rsa](#launch-rsa) : lance le serveur
* [launch-rba](#launch-rba) : lance un ballot
* [launch-rva](#launch-rva) : lance un voteur et vote
* [launch-all-agt](#launch-all-agt) : lance le serveur, des ballots et des voteurs et attend pour le résultat (les ballots ont tous les mêmes candidats, tieBreak et voteurs, mais ont des règles et deadlines différentes ; les voteurs votent dans un seul ballot décidé aléatoirement)

## Installation

Le projet peut être récupéré avec go en utilisant la commande suivante :   
`go get gitlab.utc.fr/mennynat/ia04-tp`   
Ceci va permettre de récupérer la dernière version du package et de le mettre dans `$GOPATH/pkg/mod/gitlab.utc.fr/mennynat`

Ensuite, une fois le projet récupéré, les scripts ci-dessus peuvent être installé à l'aide d'un `go install gitlab.utc.fr/mennynat/ia04-tp/cmd/<nom_launcher>` afin d'en faire des exécutables utilisables de la même manière que décrit ci-dessous. Ces exécutables seront trouvables dans `$GOPATH/bin`.

### Détail des commandes

#### launch-rsa

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

#### launch-rba

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

#### launch-rva

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

#### launch-all-agt

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
		Indique le nombre de voteurs à créer.
		Défaut : 10

	-t, --tiebreak tiebreak
		Indique le tiebreak à utiliser lors des votes.
		Format : alt1,alt2,alt3
		Défaut : 4,2,3,5,9,8,7,1,6,11,12,10

```

## Tests

Des fichiers tests ont été créés afin de tester les méthodes de vote et les agents. Pour les lancer :

❗❗ Pour les tests avec les agents, il faut lancer au préalable un serveur REST à l'adresse "localhost:8080" à l'aide de [launch-rsa](#launch-rsa) 

```bash
go test ./tests/
```
