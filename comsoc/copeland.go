package comsoc

// Le meilleur candidat est celui qui bat le plus d’autres candidats
// On associe à chaque candidat a le score suivant :
// pour chaque autre candidat b≠a +1 si une majorité préfère a à b,
// -1 si une majorité préfère b à a et 0 sinon
// Le candidat élu est celui qui a le plus haut score de Copeland
func CopelandSWF(p Profile) (Count, error) {
	alts := p[0]
	err := checkProfileAlternative(p, alts)
	if err != nil {
		return make(Count), err
	}

	count := make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}

	for idx1, alt1 := range alts {
		for idx2 := idx1 + 1; idx2 < len(alts); idx2++ {
			alt2 := alts[idx2]
			if alt1 == alt2 {
				continue
			}
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
		}
	}
	return count, nil
}

func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)
	bestAlts = maxCount(count)
	return bestAlts, err
}
