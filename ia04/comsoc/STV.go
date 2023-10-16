package comsoc

// Changement signature

// Chaque individu indique donne son ordre de préférence $>_i$
// Pour n candidats, on fait n-1 tours
// (à moins d’avoir avant une majorité stricte pour un candidat)
// On suppose qu’à chaque tour chaque individu “vote” pour son candidat préféré
// (parmi ceux encore en course)
// À chaque tour on élimine le plus mauvais candidat
// (celui qui a le moins de voix)
func STV_SWF(p Profile) (count Count, err error) {

	alts := p[0]
	err = checkProfileAlternative(p, alts)
	if err != nil {
		return count, err
	}

	count = make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}

	var eliminated []Alternative

	for round := 0; round < len(p); round++ {

		if len(eliminated) == len(alts) {
			break
		}

		countRound := make(Count)
		for _, alt := range alts {
			skip := false
			for _, eli := range eliminated {
				if alt == eli {
					skip = true
					break
				}
			}
			if !skip {
				countRound[alt] = 0
			}
		}

		for _, prefs := range p {
			for _, alt := range prefs {
				skip := false
				for _, eli := range eliminated {
					if alt == eli {
						skip = true
						break
					}
				}
				if !skip {
					countRound[alt] += 1
					break
				}
			}

		}

		worsts := minCount(countRound)
		for _, alt := range worsts {
			count[alt] = round
		}
		eliminated = append(eliminated, worsts...)
	}
	return count, err
}

func STV_SCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := STV_SWF(p)
	bestAlts = maxCount(count)
	return bestAlts, err
}
