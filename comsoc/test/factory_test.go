package comsoc_test

import (
	"testing"

	"ia04-tp/comsoc"
)

func TestSWFFactory(t *testing.T) {
	profile := comsoc.Profile{
		{1, 2, 3, 4},
		{2, 1, 3, 4},
		{1, 2, 3, 4},
		{2, 1, 3, 4},
		{3, 1, 2, 4},
		{4, 1, 3, 2},
	}
	preferences := []comsoc.Alternative{3, 1, 2, 4}

	tiebreak := comsoc.TieBreakFactory(preferences)

	swf := comsoc.SWFFactory(comsoc.MajoritySWF, tiebreak)

	res, err := swf(profile)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 4 {
		t.Errorf("Should contain 4 elements; but found %d", len(res))
	}

	if res[0] != 1 || res[1] != 2 || res[2] != 3 || res[3] != 4 {
		t.Errorf("The result should be [1, 2, 3, 4]; but is %v", res)
	}
}

func TestCSFFactory(t *testing.T) {
	profile := comsoc.Profile{
		{1, 2},
		{2, 1},
		{1, 2},
		{2, 1},
	}
	preferences := []comsoc.Alternative{2, 1}

	tiebreak := comsoc.TieBreakFactory(preferences)

	scf := comsoc.SCFFactory(comsoc.MajoritySCF, tiebreak)

	res, err := scf(profile)

	if err != nil {
		t.Error(err)
	}

	if res != 2 {
		t.Errorf("The result should be 2, found %d", res)
	}
}
