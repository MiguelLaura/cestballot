# Implémenté
* Condorcet, Borda, majority, approval, stv, Copeland

## Launch
* launch-rsa : lance le serveur
* launch-rba : lance un ballot
* launch-rva : lance un voter
* launch-all-agt : lance le serveur, des ballots et des voters (les ballotes ont tous les mêmes candidats, tieBreak et voter mais ont des règles et deadlines différentes; les voter votent dans un seul ballot décidé aléatoirement)

# A faire
* Implémenter des tests sur le serveur

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
