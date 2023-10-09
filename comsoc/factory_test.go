package comsoc

import (
	"testing"
)

func TestSWFFactory(t *testing.T) {
	profile := Profile{
		{1, 2, 3, 4},
		{2, 1, 3, 4},
		{1, 2, 3, 4},
		{2, 1, 3, 4},
		{3, 1, 2, 4},
		{4, 1, 3, 2},
	}
	preferences := []Alternative{3, 1, 2, 4}

	tiebreak := TieBreakFactory(preferences)

	swf := SWFFactory(MajoritySWF, tiebreak)

	res, err := swf(profile)

	if err != nil {
		t.Error(err)
	}

	if len(res) > 2 {
		t.Errorf("Should not have more than two element in the count")
	}

	if res[1] != 2 || res[3] != 1 {
		t.Errorf("Should have 2 for alt 1, and 1 for alt 3")
	}
}

func TestCSFFactory(t *testing.T) {
	profile := Profile{
		{1, 2},
		{2, 1},
		{1, 2},
		{2, 1},
	}
	preferences := []Alternative{2, 1}

	tiebreak := TieBreakFactory(preferences)

	scf := SCFFactory(MajoritySCF, tiebreak)

	res, err := scf(profile)

	if err != nil {
		t.Error(err)
	}

	if res != 2 {
		t.Errorf("The result should be 2, found %d", res)
	}
}
