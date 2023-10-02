package comsoc

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	resDuel, err := scoring(p)

	if err != nil {
		return nil, err
	}

	alts := p[0]

	for alt, score := range resDuel {
		if score == len(alts)-1 {
			bestAlts = append(bestAlts, alt)
			return
		}
	}

	return
}
