package comsoc

// Les alternatives seront représentées par des entiers.
type Alternative int

// Les profils de préférences sont telles que si profile est un profil,
// profile[12] représentera les préférences du votant 12.
// Les alternatives sont classées de la préférée à la moins préférée :
// profile[12][0] represente l'alternative préférée du votant 12.
type Profile [][]Alternative

// Enfin, les méthodes de vote renvoient un décompte sous forme d'une map
// qui associe à chaque alternative un entier :
// plus cet entier est élevé, plus l'alternative a de points
// et plus elle est préférée pour le groupe compte tenu de la méthode considérée.
type Count map[Alternative]int
