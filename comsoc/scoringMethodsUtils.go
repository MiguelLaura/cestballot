package comsoc

func generateCombinations(alts []Alternative) chan []Alternative {
	chnl := make(chan []Alternative)

	go func() {
		defer close(chnl)

		for idxAlt1 := 0; idxAlt1 < len(alts)-1; idxAlt1++ {
			for idxAlt2 := idxAlt1 + 1; idxAlt2 < len(alts); idxAlt2++ {
				chnl <- []Alternative{alts[idxAlt1], alts[idxAlt2]}
			}
		}
	}()

	return chnl
}

func scoring(p Profile) (resDuel Count, err error) {
	if err = checkProfileAlternative(p, p[0]); err != nil {
		return
	}

	alts := p[0]
	resDuel = make(Count)

	for _, alt := range alts {
		resDuel[alt] = 0
	}

	for combi := range generateCombinations(alts) {
		scoreA, scoreB := 0, 0
		a, b := combi[0], combi[1]

		for _, pref := range p {
			if isPref(a, b, pref) {
				scoreA += 1
			} else {
				scoreB += 1
			}
		}

		if scoreA > scoreB {
			resDuel[a]++
		} else if scoreB > scoreA {
			resDuel[a]--
		}
		// If equal, do nothing...
		// But could do a tie break to choose (randomly or not) one alt or the other
	}

	return resDuel, nil
}
