# Application de vote

Auteur·rices : Laura Miguel, Nathan Menny

Application de vote en GO réalisée dans le cadre des cours d'IA04.

## Utilisation

Tous les scripts suivants se trouvent dans le répertoire cmd/, il est ainsi possible de les compiler avec un `go build` ou `go install`.

### Serveur REST

Le serveur REST doit être lancé en premier, c'est sur lui que les requêtes se font et c'est lui qui les dirigent vers les bons bureaux de vote.

```
	launch-rest-server [flags]

Avec les flags suivants : 

	--host hôte
		Indique l'hôte sur lequel le serveur tourne.
		Défaut : localhost
	--port portNumber
		Indique le port sur lequel le serveur écoute.
		Défaut : 8080
	-v
		Lance le serveur en mode verbeux.
```

### Création de bureaux de vote

Les bureaux de vote supportent les méthodes suivantes : majorité, borda, stv, copeland, approbation

```
	create-ballot [flags]

Avec les flags suivants : 

	--host hôte
		Indique l'hôte sur lequel le serveur tourne.
		Défaut : localhost
	--port portNumber
		Indique le port sur lequel le serveur écoute.
		Défaut : 8080
	--rule nomRègle
		Spécifie la méthode de vote utilisée.
		Une des valeurs suivantes : [majority, borda, stv, copeland, approval]
		Défaut : majority
	--deadline dateLimite
		Indique la date de fermeture du vote au format RFC3339.
		Défaut : Temps courant + 1 minute
	--voters listeVotants
		Spécifie tous les votants autorisés à voter dans ce bureau de vote.
		Format : voter1,voter2,voter3,...
		Défaut : ag_id1,ag_id2,ag_id3
	--tiebreak listeAlternatives
		Spécifie le tiebreak utilisé pour départager les alternatives en cas d'égalités.
		Format : alt1,alt2,alt3,...
		Défaut : 4,2,3,1
```

### Création d'agents voteurs

```
	create-voter [flags]

Avec les flags suivants : 

	--host hôte
		Indique l'hôte sur lequel le serveur tourne.
		Défaut : localhost
	--port portNumber
		Indique le port sur lequel le serveur écoute.
		Défaut : 8080
	--id idAgent
		L'ID de l'agent votant.
		Défaut : ag_id1
	--name nomAgent
		Spécifie le nom de l'agent votant.
		Défaut : ag_id1
	--prefs listeAlternatives
		Spécifie les préférences de l'agent.
		Format : alt1,alt2,alt3,...
		Défaut : 1,2,4,3
	--opts listeOptions
		[Optionnel] Spécifie une liste d'options influençant la méthode de vote.
		Format : opt1,opt2,...
		Défaut :
	--ballot idBureauVote
		Spécifie le bureau de vote où les agents votent.
		Défaut : scrutin1
```

### Programme de test

Lance un serveur REST ainsi que des Bureaux de votes ou un certain nombre d'agent votants
vont aller voter de manière aléatoire.

```
	launch-all [flags]

Avec les flags suivants : 

	--host hôte
		Indique l'hôte sur lequel le serveur tourne.
		Défaut : localhost
	--port portNumber
		Indique le port sur lequel le serveur écoute.
		Défaut : 8080
	--nbBallots nbBallots
		Indique le nombre de bureaux de votes à créer.
		Défaut : 2
	--nbVoters nbVoters
		Indique le nombre d'agents votants.
		Défaut : 20
	--nbAlts nbAlts
		Indique le nombre d'alternatives.
		Défaut : 15
	-v
		Lance le serveur en mode verbeux.
```

## Hypothèses d'utilisation

- Un ballot ne peux avoir moins de deux alternatives
- Un ballot doit avoir au moins un votant
- Si deux alternatives ou plus se retrouvent ex æquo c'est celle qui a le plus petit rang dans le tiebreak donné au bureau de vote qui est choisi
- Si un bureau de vote ferme sans avoir reçu de votes c'est le choix préféré du tiebreak qui est choisi
- Pour un vote par approbation, si l'agent votant ne donne pas la limte à laquelle il vote, on considère qu'il vote seulement pour son préféré
