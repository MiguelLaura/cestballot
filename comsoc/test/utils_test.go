package comsoc_test

import (
	"slices"
	"testing"

	"ia04-tp/comsoc"
)

func TestRank(t *testing.T) {
	alts := []comsoc.Alternative{1, 2, 3, 4, 5}

	if comsoc.Rank(2, alts) != 1 {
		t.Errorf("The rank of 2 should be 1, but found %d", comsoc.Rank(2, alts))
	}

	if comsoc.Rank(0, alts) != -1 {
		t.Errorf("The rank of 0 should be -1, but found %d", comsoc.Rank(0, alts))
	}
}

func TestIsPref(t *testing.T) {
	alts := []comsoc.Alternative{1, 2, 3, 4, 5}

	if !comsoc.IsPref(3, 5, alts) {
		t.Errorf("3 should be prefered to 5")
	}

	if comsoc.IsPref(3, 0, alts) || comsoc.IsPref(0, 3, alts) {
		t.Errorf("It's not supposed to be true since 0 is not in the alternatives")
	}

	if comsoc.IsPref(4, 2, alts) {
		t.Errorf("4 should not be prefered to 2")
	}
}

func TestMaxCount(t *testing.T) {
	count := comsoc.Count{
		1: 2,
		2: 1,
		3: 2,
	}

	bestAlts := comsoc.MaxCount(count)

	if len(bestAlts) != 2 {
		t.Errorf("There should be 2 best alternatives")
	}

	if !slices.Contains(bestAlts, 1) || !slices.Contains(bestAlts, 3) {
		t.Errorf("The two best alternatives should be 1 and 3, but found %d and %d", bestAlts[0], bestAlts[1])
	}
}

func TestCheckProfileAlternative(t *testing.T) {
	alts := []comsoc.Alternative{3, 2, 1}
	profile1 := comsoc.Profile{
		{1, 2, 3},
		{3, 1, 2},
		{2, 3, 1},
	}

	profile2 := comsoc.Profile{
		{1, 2, 3},
		{3, 1, 2, 1},
		{2, 3, 1},
	}

	profile3 := comsoc.Profile{
		{1, 2, 3},
		{3, 1, 2, 7},
		{2, 3, 1},
	}

	profile4 := comsoc.Profile{
		{1, 2, 3},
		{3, 1, 2},
		{2, 1},
	}

	if comsoc.CheckProfileAlternative(profile1, alts) != nil {
		t.Errorf("Should be a correct profile")
	}

	if comsoc.CheckProfileAlternative(profile2, alts) == nil {
		t.Errorf("Should not be correct (two times the same alternative)")
	}

	if comsoc.CheckProfileAlternative(profile3, alts) == nil {
		t.Errorf("Should not be correct (one more alternative)")
	}

	if comsoc.CheckProfileAlternative(profile4, alts) == nil {
		t.Errorf("Should not be correct (one alternative is missing)")
	}
}
