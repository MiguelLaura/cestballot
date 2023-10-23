package comsoc

func MajoritySWF(p Profile) (count Count, err error) {
	alts := p[0]
	err = checkProfileAlternative(p, alts)
	if err != nil {
		return count, err
	}

	count = make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}

	for _, alts := range p {
		count[alts[0]] += 1
	}

	return count, err
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)
	bestAlts = maxCount(count)
	return bestAlts, err
}
