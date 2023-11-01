// Package comsoc handles the voting methods.
//
// Handles the Copeland voting method.
package comsoc

// CopelandSWF provides the Social Welfare Function of copeland method.
func CopelandSWF(p Profile) (count Count, err error) {
	if err = checkProfileAlternative(p, p[0]); err != nil {
		return
	}

	count, err = duel(p)

	return
}

// CopelandSCF provides the Social Choice Function of copeland method.
func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)

	if err != nil {
		return nil, err
	}

	bestAlts = maxCount(count)
	return
}
