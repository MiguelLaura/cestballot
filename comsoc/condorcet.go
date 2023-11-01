// Package comsoc handles the voting methods.
//
// Provides the methods to find the Condorcet winner if exists.
package comsoc

// CondorcetWinner gives the condorcet winner if exists
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	resDuel, err := duel(p)

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
