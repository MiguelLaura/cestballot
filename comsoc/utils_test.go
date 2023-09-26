package comsoc

import (
	"slices"
	"testing"
)

func TestRank(t *testing.T) {
	alts := []Alternative{1, 2, 3, 4, 5}

	if rank(2, alts) != 1 {
		t.Errorf("The rank of 2 should be 1, but found %d", rank(2, alts))
	}

	if rank(0, alts) != -1 {
		t.Errorf("The rank of 0 should be -1, but found %d", rank(0, alts))
	}
}

func TestIsPref(t *testing.T) {
	alts := []Alternative{1, 2, 3, 4, 5}

	if !isPref(3, 5, alts) {
		t.Errorf("3 should be prefered to 5")
	}

	if isPref(3, 0, alts) || isPref(0, 3, alts) {
		t.Errorf("It's not supposed to be true since 0 is not in the alternatives")
	}

	if isPref(4, 2, alts) {
		t.Errorf("4 should not be prefered to 2")
	}
}

func TestMaxCount(t *testing.T) {
	count := Count{
		1: 2,
		2: 1,
		3: 2,
	}

	bestAlts := maxCount(count)

	if len(bestAlts) != 2 {
		t.Errorf("There should be 2 best alternatives")
	}

	if !slices.Contains(bestAlts, 1) || !slices.Contains(bestAlts, 3) {
		t.Errorf("The two best alternatives should be 1 and 3, but found %d and %d", bestAlts[0], bestAlts[1])
	}
}

func TestCheckProfile(t *testing.T) {
	profile1 := Profile{
		[]Alternative{1, 2, 3},
		[]Alternative{3, 1, 2},
		[]Alternative{2, 3, 1},
	}

	profile2 := Profile{
		[]Alternative{1, 2, 3},
		[]Alternative{3, 1, 2, 1},
		[]Alternative{2, 3, 1},
	}

	profile3 := Profile{
		[]Alternative{1, 2, 3},
		[]Alternative{3, 1, 2, 7},
		[]Alternative{2, 3, 1},
	}

	profile4 := Profile{
		[]Alternative{1, 2, 3},
		[]Alternative{3, 1, 2},
		[]Alternative{2, 1},
	}

	if checkProfile(profile1) != nil {
		t.Errorf("Should be a correct profile")
	}

	if checkProfile(profile2) == nil {
		t.Errorf("Should not be correct (two times the same alternative)")
	}

	if checkProfile(profile3) == nil {
		t.Errorf("Should not be correct (one more alternative)")
	}

	if checkProfile(profile4) == nil {
		t.Errorf("Should not be correct (one alternative is missing)")
	}
}
