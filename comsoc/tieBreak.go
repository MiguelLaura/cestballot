// Package comsoc handles the voting methods.
//
// Handles tiebreaks with SCF and SWF.
package comsoc

import (
	"errors"
	"slices"
	"sort"

	"ia04-tp/utils/sequential"
)

// TieBreakFactory creates a simple tiebreak that returns the first best alternative.
func TieBreakFactory(allAlts []Alternative) func([]Alternative) (alt Alternative, err error) {

	return func(currAlts []Alternative) (alt Alternative, err error) {
		if len(currAlts) == 0 && len(allAlts) == 0 {
			return 0, errors.New("no alternatives")
		}

		if len(currAlts) == 0 {
			return allAlts[0], nil
		}

		alt = currAlts[0]

		for altIdx := 1; altIdx < len(currAlts); altIdx++ {
			// Chooses the best alternative (i.e. the one with the minimal rank)
			if IsPref(currAlts[altIdx], alt, allAlts) {
				alt = currAlts[altIdx]
			}
		}
		return
	}

}

// SWFFactory creates a SWF without doubles.
func SWFFactory(
	swf func(Profile) (Count, error),
	tiebreak func([]Alternative) (Alternative, error),
) func(Profile) ([]Alternative, error) {

	return func(p Profile) ([]Alternative, error) {
		if len(p) == 0 {
			return nil, errors.New("the given profile is empty")
		}

		count, err := swf(p)

		if err != nil {
			return nil, err
		}

		// Gets all the alternatives and map them to the number of votes tey received
		exAequo := make(map[int][]Alternative)

		for alt, nbVotes := range count {
			if _, err := exAequo[nbVotes]; !err {
				exAequo[nbVotes] = make([]Alternative, 1)
				exAequo[nbVotes][0] = alt
			} else {
				exAequo[nbVotes] = append(exAequo[nbVotes], alt)
			}
		}

		// Gets the number of vote and orders it in descending order
		nbVotes := make([]int, 0, len(exAequo))
		for nbVote := range exAequo {
			nbVotes = append(nbVotes, nbVote)
		}
		sort.Ints(nbVotes)
		slices.Reverse(nbVotes)

		// Constructs the list without exaequo
		res := make([]Alternative, len(count))
		resIdx := 0

		for _, nbVote := range nbVotes {

			altNbVote := exAequo[nbVote]
			nbAlts := len(altNbVote)

			for altIdx := 0; altIdx < nbAlts; altIdx++ {

				// Gets the best alternative available
				bestAlt, err := tiebreak(altNbVote)
				if err != nil {
					return nil, err
				}

				res[resIdx] = bestAlt

				// Removes the best alternative from all the previous alternatives
				idxBestAlt, _ := sequential.Find(altNbVote, func(alt Alternative) bool { return alt == bestAlt })
				altNbVote = append(altNbVote[:idxBestAlt], altNbVote[idxBestAlt+1:]...)

				resIdx++
			}
		}

		return res, nil
	}
}

// SCFFactory creates a SCF returning the one best alternative.
func SCFFactory(
	scf func(Profile) ([]Alternative, error),
	tiebreak func([]Alternative) (Alternative, error),
) func(Profile) (Alternative, error) {

	return func(p Profile) (Alternative, error) {
		if len(p) == 0 {
			return tiebreak([]Alternative{})
		}
		bestAlts, err := scf(p)
		if err != nil {
			return 0, err
		}

		return tiebreak(bestAlts)
	}

}
