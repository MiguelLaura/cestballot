package comsoc

// Qui renvoie un slice éventuellement vide ou ne contenant qu'un seul élément
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	alts := p[0]
	err = checkProfileAlternative(p, alts)
	if err != nil {
		return bestAlts, err
	}

	count := make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}

	end := false
	for idx1, alt1 := range alts {
		for idx2 := idx1 + 1; idx2 < len(alts); idx2++ {
			alt2 := alts[idx2]
			// if alt1 == alt2 {
			// 	continue
			// }
			countAlt1WinsAlt2 := 0
			for _, prefs := range p {
				if isPref(alt1, alt2, prefs) {
					countAlt1WinsAlt2 += 1
				} else {
					countAlt1WinsAlt2 -= 1
				}
				if countAlt1WinsAlt2 > len(p)/2 {
					break
				} else if countAlt1WinsAlt2 < -len(p)/2 {
					break
				}
			}
			if countAlt1WinsAlt2 > 0 {
				// fmt.Println(alt1, "wins", alt2)
				count[alt1] += 1
				count[alt2] -= 1
			} else if countAlt1WinsAlt2 < 0 {
				// fmt.Println(alt2, "wins", alt1)
				count[alt1] -= 1
				count[alt2] += 1
			}
			// fmt.Println(count)
			if count[alt1] == len(alts)-1 {
				bestAlts = []Alternative{alt1}
				end = true
				break
			}
		}
		if end {
			break
		}
	}
	best := maxCount(count)
	if len(best) == 1 && count[best[0]] == len(alts)-1 {
		bestAlts = []Alternative{best[0]}
	}
	return bestAlts, nil
}
