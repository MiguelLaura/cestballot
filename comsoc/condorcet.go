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

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	if err = checkProfile(p); err != nil {
		return
	}

	alts := p[0]
	resDuel := make(Count)

	for _, alt := range alts {
		resDuel[alt] = 0
	}

	for combi := range generateCombinations(alts) {
		score1, score2 := 0, 0

		for _, pref := range p {
			if isPref(combi[0], combi[1], pref) {
				score1 += 1
			} else {
				score2 += 1
			}
		}

		if score1 > score2 {
			resDuel[combi[0]]++
		} else if score2 > score1 {
			resDuel[combi[1]]++
		}
		// If equal, do nothing...
		// But could do a tie break to choose (randomly or not) one alt or the other
	}

	for alt, score := range resDuel {
		if score == len(alts)-1 {
			bestAlts = append(bestAlts, alt)
			return
		}
	}

	return
}
