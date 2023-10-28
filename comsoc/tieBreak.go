package comsoc

import (
	"errors"
	"sort"
)

func TieBreakFactory(orderedAlts []Alternative) func([]Alternative) (Alternative, error) {
	return func(alts []Alternative) (bestAlt Alternative, err error) {

		if alts == nil {
			err := "erreur : le slice d'alternatives est vide"
			return bestAlt, errors.New(err)
		}

		if len(alts) == 1 {
			return alts[0], nil
		}

		bestAlt = alts[0]
		for _, alt := range alts {
			if isPref(alt, bestAlt, orderedAlts) {
				bestAlt = alt
			}
		}
		return bestAlt, nil
	}
}

// [A FAIRE]
// Changement dans signature
func SWFFactory(swf func(p Profile) (Count, error), tieBreaker func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	return func(p Profile) ([]Alternative, error) {
		count, err := swf(p)
		if err != nil {
			return nil, err
		}

		exAequo := make(map[int][]Alternative)
		res := make([]Alternative, len(count))

		for alt, votes := range count {
			if _, err := exAequo[votes]; !err {
				exAequo[votes] = make([]Alternative, 1)
				exAequo[votes][0] = alt
			} else {
				exAequo[votes] = append(exAequo[votes], alt)
			}
		}

		keys := make([]int, 0, len(exAequo))
		for k := range exAequo {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		idxRes := len(count) - 1

		for i := 0; i < len(keys); i++ {
			alts := exAequo[keys[i]]
			nbAlts := len(alts)
			for j := 0; j < nbAlts; j++ {
				bestAlt, err := tieBreaker(alts)
				if err != nil {
					return nil, err
				}
				res[idxRes-(nbAlts-1-j)] = bestAlt
				for idx, alt := range alts {
					if alt == bestAlt {
						alts = append(alts[:idx], alts[idx+1:]...)
						break
					}
				}
			}
			idxRes -= nbAlts
		}
		return res, nil
	}
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tieBreaker func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	return func(p Profile) (Alternative, error) {
		bestAlts, err := scf(p)
		if err != nil {
			return 0, err
		}
		return tieBreaker(bestAlts)
	}

}

// [RAJOUT]
func SWFFactoryApproval(swf func(p Profile, thresholds []int) (Count, error), tieBreaker func([]Alternative) (Alternative, error)) func(Profile, []int) ([]Alternative, error) {
	return func(p Profile, thresholds []int) ([]Alternative, error) {
		count, err := swf(p, thresholds)
		if err != nil {
			return nil, err
		}

		exAequo := make(map[int][]Alternative)
		res := make([]Alternative, len(count))

		for alt, votes := range count {
			if _, err := exAequo[votes]; !err {
				exAequo[votes] = make([]Alternative, 1)
				exAequo[votes][0] = alt
			} else {
				exAequo[votes] = append(exAequo[votes], alt)
			}
		}

		keys := make([]int, 0, len(exAequo))
		for k := range exAequo {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		idxRes := len(count) - 1

		for i := 0; i < len(keys); i++ {
			alts := exAequo[keys[i]]
			nbAlts := len(alts)
			for j := 0; j < nbAlts; j++ {
				bestAlt, err := tieBreaker(alts)
				if err != nil {
					return nil, err
				}
				res[idxRes-(nbAlts-1-j)] = bestAlt
				for idx, alt := range alts {
					if alt == bestAlt {
						alts = append(alts[:idx], alts[idx+1:]...)
						break
					}
				}
			}
			idxRes -= nbAlts
		}
		return res, nil
	}
}

// [RAJOUT]
func SCFFactoryApproval(scf func(p Profile, thresholds []int) ([]Alternative, error), tieBreaker func([]Alternative) (Alternative, error)) func(Profile, []int) (Alternative, error) {
	return func(p Profile, thresholds []int) (Alternative, error) {
		bestAlts, err := scf(p, thresholds)
		if err != nil {
			return 0, err
		}
		return tieBreaker(bestAlts)
	}

}
