// Package comsoc handles the voting methods.
//
// Handles the Approval voting method.
package comsoc

import "errors"

// ApprovalSWF provides the Social Welfare Function of approval method.
// The thresholds of each voter are provided to know where they stop voting.
func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	if err = checkProfileAlternative(p, p[0]); err != nil {
		return nil, err
	}

	if len(thresholds) != len(p) {
		return nil, errors.New("error, treshold not good length")
	}

	count = make(Count)

	// Process the vote for each agent
	for voter, prefs := range p {
		for rankPref, pref := range prefs {
			if rankPref >= thresholds[voter] {
				break
			}

			if _, ok := count[pref]; !ok {
				count[pref] = 1
			} else {
				count[pref]++
			}
		}
	}
	return
}

// ApprovalSCF provides the Social Choice Function of approval method.
// The thresholds of each voter are provided to know where they stop voting.
func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)

	if err != nil {
		return nil, err
	}

	bestAlts = maxCount(count)
	return
}
