package comsoc

func CopelandSWF(p Profile) (count Count, err error) {
	if err = checkProfileAlternative(p, p[0]); err != nil {
		return
	}

	count, err = scoring(p)

	return
}

func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)

	if err != nil {
		return nil, err
	}

	bestAlts = maxCount(count)
	return
}
