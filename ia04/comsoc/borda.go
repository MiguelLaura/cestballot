package comsoc

func BordaSWF(p Profile) (count Count, err error) {
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
		for i := 0; i < len(alts); i++ {
			count[alts[len(alts)-1-i]] += i
		}
	}

	return count, err

}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	bestAlts = maxCount(count)
	return bestAlts, err
}
