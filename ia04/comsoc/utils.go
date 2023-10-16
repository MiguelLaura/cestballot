package comsoc

import (
	"errors"
	"fmt"
)

// renvoie l'indice où se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	for idx, val := range prefs {
		if val == alt {
			return idx
		}
	}
	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	for _, val := range prefs {
		if val == alt1 {
			return true
		} else if val == alt2 {
			return false
		}
	}
	return false
	// rankAlt1 := rank(alt1, prefs)
	// rankAlt2 := rank(alt2, prefs)
	// if rankAlt1 != -1 && rankAlt2 != -1 && rankAlt1 < rankAlt2 {
	// 	return true
	// }
	// return false
}

// renvoie les meilleures alternatives pour un décompte donné
func maxCount(count Count) (bestAlts []Alternative) {
	for key, val := range count {
		if bestAlts == nil {
			bestAlts = append(bestAlts, key)
		} else if val > count[bestAlts[0]] {
			bestAlts = []Alternative{key}
		} else if val == count[bestAlts[0]] {
			bestAlts = append(bestAlts, key)
		}
	}
	return bestAlts
}

// renvoie les meilleures alternatives pour un décompte donné
func minCount(count Count) (worstAlts []Alternative) {
	for key, val := range count {
		if worstAlts == nil {
			worstAlts = append(worstAlts, key)
		} else if val < count[worstAlts[0]] {
			worstAlts = []Alternative{key}
		} else if val == count[worstAlts[0]] {
			worstAlts = append(worstAlts, key)
		}
	}
	return worstAlts
}

// vérifie les préférences d'un agent,
// par ex. qu'ils sont tous complets et
// que chaque alternative n'apparaît qu'une seule fois
func checkProfile(prefs []Alternative, alts []Alternative) error {
	nbAlts := len(alts)
	if len(prefs) != nbAlts {
		err := "erreur : prefs de taille différente de alts (len(prefs)==" + fmt.Sprint(len(prefs)) + " et len(alts)==" + fmt.Sprint(nbAlts) + ")"
		return errors.New(err)
	}
	for _, alt1 := range alts {
		nbAlt1 := 0
		for _, alt2 := range prefs {
			if alt1 == alt2 {
				nbAlt1 += 1
			}
			if nbAlt1 > 1 {
				err := "erreur : au moins deux fois la même alternative (" + fmt.Sprint(alt1) + ") dans prefs"
				return errors.New(err)
			}
		}
		if nbAlt1 == 0 {
			err := "erreur : l'alternative " + fmt.Sprint(alt1) + " apparait 0 fois dans prefs"
			return errors.New(err)
		}
	}
	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative
// de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	for idx, profile := range prefs {
		err := checkProfile(profile, alts)
		if err != nil {
			return fmt.Errorf("%w; pour prefs de profile["+fmt.Sprint(idx)+"]", err)
		}
	}
	return nil
}
