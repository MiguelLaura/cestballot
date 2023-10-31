package voteragent

import "gitlab.utc.fr/mennynat/ia04-tp/comsoc"

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt comsoc.Alternative, prefs []comsoc.Alternative) int {
	for index, value := range prefs {
		if value == alt {
			return index
		}
	}

	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 comsoc.Alternative, prefs []comsoc.Alternative) bool {
	rk1, rk2 := rank(alt1, prefs), rank(alt2, prefs)
	return rk1 != -1 && rk2 != -1 && rk1 < rk2
}
