package comsoc

import "errors"

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	if err = checkProfileAlternative(p, p[0]); err != nil {
		return nil, err
	}

	if len(thresholds) != len(p) {
		return nil, errors.New("error, treshold not good length")
	}

	count = make(Count)

	for indiv, prefs := range p {
		for rankPref, pref := range prefs {
			if rankPref >= thresholds[indiv] {
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

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)

	if err != nil {
		return nil, err
	}

	bestAlts = maxCount(count)
	return
}
