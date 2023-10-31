# A corriger
* changer les parcours de listes pour vérifier la présence d'un élément avec slices.Contains
* implémenter des tests sur le serveur

# Hypothèses
* On a au moins 1 voter
* On a au moins 2 candidats/alternatives
* Quand il n'y a pas de SWF (Condorcet), on renvoie nil/null
* Pas d'erreur quand on donne des options alors que la règle n'est pas approval
* Erreur quand on donne en option autre chose qu'un seul nombre dans une slice et qu'on va utiliser l'option
* Les agtid peuvent s'abstenir
* Quand tous les agt s'abstiennent, on renvoie le résultat à partir de tiebreak

# Launch all
* Créé des ballot avec les mêmes candidats, même tieBreak, même voter, mais des règles et deadlines différentes.
* Les voter votent dans un seul ballot décidé aléatoirement.
