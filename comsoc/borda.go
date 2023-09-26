package comsoc

func BordaSWF(p Profile) (count Count, err error) {
	if err = checkProfile(p); err != nil {
		return
	}

	count = make(Count)

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

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)

	if err != nil {
		return nil, err
	}

	bestAlts = maxCount(count)
	return
}
