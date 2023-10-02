package comsoc

func removeAlt(alt Alternative, alts []Alternative) []Alternative {
	idxAlt := rank(alt, alts)

	if idxAlt != -1 {
		return append(alts[:idxAlt], alts[idxAlt+1:]...)
	}

	return alts
}

func removeAltFromProfile(alt Alternative, p *Profile) {
	for indiv, prefs := range *p {
		(*p)[indiv] = removeAlt(alt, prefs)
	}
}

func STV_SWF(p Profile) (Count, error) {
	return MajoritySWF(p)
}

func STV_SCF(p Profile) (bestAlts []Alternative, err error) {
	for n := len(p[0]); n > 0; n-- {
		count, err := STV_SWF(p)
		if err != nil {
			return nil, err
		}

		var lessVoteAlt Alternative = -1
		var lessVoteCnt int = len(p) + 1

		// Check if one candidate has an absolute majority
		for alt, cnt := range count {
			if cnt > len(p)/2 {
				return []Alternative{alt}, nil
			}

			if cnt < lessVoteCnt {
				lessVoteAlt = alt
				lessVoteCnt = cnt
			}
		}

		// Removes the least liked candidate
		removeAltFromProfile(lessVoteAlt, &p)
	}

	return []Alternative{p[0][0]}, nil
}