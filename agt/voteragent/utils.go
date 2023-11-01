// Package voteragent contains an agent that can vote.
package voteragent

import "gitlab.utc.fr/mennynat/ia04-tp/comsoc"

// rank sends the index of the given alternative in the preferences.
func rank(alt comsoc.Alternative, prefs []comsoc.Alternative) int {
	for index, value := range prefs {
		if value == alt {
			return index
		}
	}

	return -1
}

// isPref indicates if alt1 is preferred to alt2.
func isPref(alt1, alt2 comsoc.Alternative, prefs []comsoc.Alternative) bool {
	rk1, rk2 := rank(alt1, prefs), rank(alt2, prefs)
	return rk1 != -1 && rk2 != -1 && rk1 < rk2
}
