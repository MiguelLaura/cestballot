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

// Changement dans signature A CORRIGER
func SWFFactory(swf func(p Profile) (Count, error), tieBreaker func([]Alternative) (Alternative, error)) func(Profile) (Count, error) {
	return func(p Profile) (Count, error) {
		count, err := swf(p)
		if err != nil {
			return nil, err
		}

		exAequo := make(map[int][]Alternative)
		res := make(Count)

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

		add := 0
		for _, k := range keys {
			alts := exAequo[k]
			newAdd := len(alts) - 1
			for i := newAdd; i > -1; i-- {
				bestAlt, err := tieBreaker(alts)
				if err != nil {
					return nil, err
				}
				res[bestAlt] = count[bestAlt] + add + i
				for idx, alt := range alts {
					if alt == bestAlt {
						alts = append(alts[:idx], alts[idx+1:]...)
						break
					}
				}
			}
			add += newAdd
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
