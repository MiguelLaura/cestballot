package comsoc

import (
	"errors"
	"slices"
)

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	for index, value := range prefs {
		if value == alt {
			return index
		}
	}

	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	rk1, rk2 := rank(alt1, prefs), rank(alt2, prefs)
	return rk1 != -1 && rk2 != -1 && rk1 < rk2
}

// renvoie les meilleures alternatives pour un décomtpe donné
func maxCount(count Count) (bestAlts []Alternative) {
	bestAlts = make([]Alternative, 0)
	bestCnt := 0
	for alt, cnt := range count {
		if bestCnt < cnt {
			bestAlts = make([]Alternative, 1)
			bestAlts[0] = alt
			bestCnt = cnt
		} else if bestCnt == cnt {
			bestAlts = append(bestAlts, alt)
		}
	}
	return
}

func checkDoublons(pref []Alternative) bool {
	cpPref := make([]Alternative, len(pref))
	copy(cpPref, pref)
	slices.Sort(cpPref)

	for idx := 1; idx < len(cpPref); idx++ {
		if cpPref[idx-1] == cpPref[idx] {
			return true
		}
	}

	return false
}

func equals(pref1 []Alternative, pref2 []Alternative) bool {
	if len(pref1) != len(pref2) {
		return false
	}

	cpPref1, cpPref2 := make([]Alternative, len(pref1)), make([]Alternative, len(pref2))
	copy(cpPref1, pref1)
	copy(cpPref2, pref2)
	slices.Sort(cpPref1)
	slices.Sort(cpPref2)

	for idx := range cpPref1 {
		if cpPref1[idx] != cpPref2[idx] {
			return false
		}
	}

	return true
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
/* func checkProfile(prefs Profile) error {
	if len(prefs) == 0 {
		return nil
	}

	return checkProfileAlternative(prefs[1:], prefs[0])
} */

// vérifie les préférences d'un agent, par ex. qu'ils sont tous complets et que chaque alternative n'apparait qu'une seule fois
func checkProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) == 0 && len(alts) == 0 {
		return nil
	}

	// Checks if the current individual has two identical individuals
	if checkDoublons(prefs) {
		return errors.New("two same alternatives found for the same individual")
	}

	// Checks if the current individual does not have the same alternatives as the others
	if !equals(alts, prefs) {
		return errors.New("at least one alternative does not exist in the preferences")
	}

	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error {

	// Checks if there are two times the same alternative in the alts slice
	if checkDoublons(alts) {
		return errors.New("two same alternatives found in the alts slice")
	}

	// Checks for the other individuals
	for indiv := 0; indiv < len(prefs); indiv++ {
		if err := checkProfile(prefs[indiv], alts); err != nil {
			return err
		}
	}

	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
/* func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	if len(prefs) == 0 {
		return nil
	}

	// Checks if the first individual has two times the same alternative
	if checkDoublons(alts) {
		return errors.New("two same alternatives found for the same individual")
	}

	alternatives := alts

	// Checks for the other individuals
	for indiv := 0; indiv < len(prefs); indiv++ {
		prefsIndiv := prefs[indiv]

		// Checks if the current individual has two identical individuals
		if checkDoublons(prefsIndiv) {
			return errors.New("two same alternatives found for the same individual")
		}

		// Checks if the current individual does not have the same alternatives as the others
		if !equals(alternatives, prefsIndiv) {
			return errors.New("at least one alternative does not exist in the other profiles")
		}
	}

	return nil
} */
