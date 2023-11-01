// Package comsoc handles the voting methods.
//
// Handles the Majority voting method.
package comsoc

// MajoritySWF provides the Social Welfare Function of majority method.
func MajoritySWF(p Profile) (count Count, err error) {
	if err = checkProfileAlternative(p, p[0]); err != nil {
		return
	}

	count = make(Count)

	for _, alternatives := range p {
		if _, ok := count[alternatives[0]]; !ok {
			count[alternatives[0]] = 1
		} else {
			count[alternatives[0]]++
		}
	}

	return
}

// MajoritySCF provides the Social Choice Function of majority method.
func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)

	if err != nil {
		return nil, err
	}

	bestAlts = maxCount(count)
	return
}
