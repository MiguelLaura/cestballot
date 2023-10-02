package comsoc

import (
	"errors"
)

// Creates a simple tie break that returns the first best alternative
func TieBreakFactory(allAlts []Alternative) func([]Alternative) (alt Alternative, err error) {

	return func(currAlts []Alternative) (alt Alternative, err error) {
		if len(currAlts) == 0 || len(allAlts) == 0 {
			return 0, errors.New("no alternatives")
		}
		alt = currAlts[0]

		for altIdx := 1; altIdx < len(currAlts); altIdx++ {
			// Chooses the best alternative (i.e. the one with the minimal rank)
			if isPref(currAlts[altIdx], alt, allAlts) {
				alt = currAlts[altIdx]
			}
		}
		return
	}

}

func SWFFactory(
	swf func(Profile) (Count, error),
	tiebreak func([]Alternative) (Alternative, error),
) func(Profile) (Count, error) {

	return func(p Profile) (Count, error) {

		countGl, err := swf(p)
		if err != nil {
			return nil, err
		}

		// Maps the count of the elements in countGl to the corresponding alternatives
		mapExequo := make(map[int][]Alternative)
		countReduced := make(Count)

		// Gets the exequo alternatives
		for alt, cnt := range countGl {
			// If no alternatives has this occurrence
			if _, ok := mapExequo[cnt]; !ok {
				mapExequo[cnt] = make([]Alternative, 1)
				mapExequo[cnt][0] = alt
			} else {
				mapExequo[cnt] = append(mapExequo[cnt], alt)
			}
		}

		// Creates a map without the exequo alternatives
		for cnt, alts := range mapExequo {
			bestAlt, _ := tiebreak(alts)
			countReduced[bestAlt] = cnt
		}

		return countReduced, err

	}
}

func SCFFactory(
	scf func(Profile) ([]Alternative, error),
	tiebreak func([]Alternative) (Alternative, error),
) func(Profile) (Alternative, error) {

	return func(p Profile) (Alternative, error) {
		bestAlts, err := scf(p)
		if err != nil {
			return 0, err
		}

		return tiebreak(bestAlts)
	}

}