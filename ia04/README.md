# A corriger
* changer les parcours de listes pour vérifier la présence d'un élément avec slices.Contains
* vérifier condorcet et threshold
* revoir les launch

# Hypothèses
* On a au moins 1 voter
* On a au moins 2 candidats/alternatives
* Quand il n'y a pas de SWF (Condorcet), on renvoie nil/null
* Pas d'erreur quand on donne des options alors que la règle n'est pas approval
* Erreur quand on donne en option autre chose qu'un seul nombre dans une slice et qu'on va utiliser l'option
