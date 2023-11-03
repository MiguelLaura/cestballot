// Package comsoc handles the voting methods.
//
// Handles the Borda voting method.
package comsoc

// BordaSWF provides the Social Welfare Function of borda method.
func BordaSWF(p Profile) (count Count, err error) {
	if err = CheckProfileAlternative(p, p[0]); err != nil {
		return
	}

	count = make(Count)

	// Process the vote for each agent
	for _, prefs := range p {
		for rankPref, pref := range prefs {
			if _, ok := count[pref]; !ok {
				count[pref] = 0
			}
			count[pref] += len(prefs) - rankPref - 1
		}
	}

	return
}

// BordaSCF provides the Social Choice Function of borda method.
func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)

	if err != nil {
		return nil, err
	}

	bestAlts = MaxCount(count)
	return
}
