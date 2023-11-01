// Package comsoc handles the voting methods.
//
// Contains useful functions to help creating votes.
package comsoc

import (
	"errors"
	"slices"
)

// rank returns the index in the preferences of the given alternative.
func rank(alt Alternative, prefs []Alternative) int {
	for index, value := range prefs {
		if value == alt {
			return index
		}
	}

	return -1
}

// isPref indicates is alt1 is preferred to alt2 given a specific preference.
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	rk1, rk2 := rank(alt1, prefs), rank(alt2, prefs)
	return rk1 != -1 && rk2 != -1 && rk1 < rk2
}

// maxCount provides the best alternatives.
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

// checkDuplicate indicates if there are duplicates in the alternatives.
func checkDuplicate(pref []Alternative) bool {
	if len(pref) <= 1 {
		return false
	}

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

// equals indicates if two preferences have the same elements (in the same order or not).
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

// CheckProfile checks if the given preference is correct.
func CheckProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) == 0 && len(alts) == 0 {
		return nil
	}

	// Checks if the current individual has two identical individuals
	if checkDuplicate(prefs) {
		return errors.New("two same alternatives found for the same individual")
	}

	// Checks if the current individual does not have the same alternatives as the others
	if !equals(alts, prefs) {
		return errors.New("at least one alternative does not exist in the preferences")
	}

	return nil
}

// checkProfileAlternative checks if every profiles are correct.
func checkProfileAlternative(prefs Profile, alts []Alternative) error {

	// Checks if there are two times the same alternative in the alts slice
	if checkDuplicate(alts) {
		return errors.New("two same alternatives found in the alts slice")
	}

	// Checks for the other individuals
	for indiv := 0; indiv < len(prefs); indiv++ {
		if err := CheckProfile(prefs[indiv], alts); err != nil {
			return err
		}
	}

	return nil
}
