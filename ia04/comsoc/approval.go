package comsoc

import (
	"errors"
	"fmt"
)

// Attention, dans ce cas il faut ajouter un nombre représentant
// le seuil à partir duquel les alternatives ne sont plus approuvées
// thresholds -> ensemble des limites jusqu'où il fait descendre pour chacun

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	alts := p[0]
	nbAlts := len(alts)

	if len(thresholds) != len(p) {
		err := "erreur : profile de taille différente de threshold (len(profile)==" + fmt.Sprint(len(p)) + " et len(threshold)==" + fmt.Sprint(len(thresholds)) + ")"
		return nil, errors.New(err)
	}

	err = checkProfileAlternative(p, alts)
	if err != nil {
		return count, err
	}

	for idx, val := range thresholds {
		if val >= nbAlts {
			thresholds[idx] = nbAlts
		}
	}

	count = make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}

	for profile, alts := range p {
		lim := thresholds[profile]
		for i := 0; i < lim; i++ {
			count[alts[i]] += 1
		}
	}
	return count, err
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	bestAlts = maxCount(count)
	return bestAlts, err
}
