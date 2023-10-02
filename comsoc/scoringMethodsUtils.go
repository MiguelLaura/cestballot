package comsoc

/* func generateCombinations(alts []Alternative) chan []Alternative {
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
} */

// Generates all the possible alternative combinations at a certain level (i.e. number of elements per combination)
func generateCombinations(alts []Alternative, level int) chan []Alternative {
	if level > len(alts) {
		return nil
	}

	chnl := make(chan []Alternative)

	go func() {
		defer close(chnl)
		generateCombinationsRec(alts, level, chnl)
	}()

	return chnl
}

func generateCombinationsRec(alts []Alternative, level int, chnl chan []Alternative, combiAlts ...Alternative) {
	switch level {
	case 1:
		for idxAlt := 0; idxAlt < len(alts); idxAlt++ {
			combi := make([]Alternative, 0)
			combi = append(combi, combiAlts...)
			combi = append(combi, alts[idxAlt])
			chnl <- combi
		}
	default:
		c := make([]Alternative, 0)
		c = append(c, combiAlts...)

		if len(alts) > 1 {
			c = append(c, 0)
			for idxAlt := 0; idxAlt < len(alts)-1; idxAlt++ {
				c[len(c)-1] = alts[idxAlt]
				generateCombinationsRec(alts[idxAlt+1:], level-1, chnl, c...)
			}
		}
	}
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

	for combi := range generateCombinations(alts, 2) {
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
