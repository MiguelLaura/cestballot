package test

import (
	"testing"

	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

func TestCheckProfileAlternativeError(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
		{3, 2, 2},
	}

	_, err := comsoc.BordaSWF(prefs)
	if err == nil {
		t.Errorf("erreur, devrait lever une erreur : erreur : l'Alternative 1 apparait 0 fois dans prefs; pour prefs de profile[3]")
	}
}

func TestCheckProfileAlternativeSizeError(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 1},
		{3, 2, 1},
	}

	_, err := comsoc.BordaSWF(prefs)
	if err == nil {
		t.Errorf("erreur, devrait lever une erreur : erreur : prefs de taille différente de alts (len(prefs)==2 et len(alts)==3); pour prefs de profile[2]")
	}
}

func TestBordaSWF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := comsoc.BordaSWF(prefs)

	if res[1] != 4 {
		t.Errorf("erreur, résultat pour 1 devrait être 4, %d calculé", res[1])
	}
	if res[2] != 3 {
		t.Errorf("erreur, résultat pour 2 devrait être 3, %d calculé", res[2])
	}
	if res[3] != 2 {
		t.Errorf("erreur, résultat pour 3 devrait être 2, %d calculé", res[3])
	}
}

func TestBordaSCF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.BordaSCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("erreur, 1 devrait être la seule meilleure Alternative")
	}
}

func TestBordaSCFMultiple(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.BordaSCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 3 {
		t.Errorf("erreur, les 3 Alternative devraient être à égalité")
	}
}

func TestMajoritySWF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := comsoc.MajoritySWF(prefs)

	if res[1] != 2 {
		t.Errorf("erreur, resultat pour 1 devrait être 2, %d calculé", res[1])
	}
	if res[2] != 0 {
		t.Errorf("erreur, résultat pour 2 devrait être 0, %d calculé", res[2])
	}
	if res[3] != 1 {
		t.Errorf("erreur, résultat pour 3 devrait être 1, %d calculé", res[3])
	}
}

func TestMajoritySCF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.MajoritySCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("erreur, 1 devrait être la seule meilleure Alternative")
	}
}

func TestMajoritySCFMultiple(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.MajoritySCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 2 {
		t.Errorf("erreur, les 2 Alternative (1, 3) devraient être à égalité")
	}
}

func TestApprovalSWF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 3, 2},
		{2, 3, 1},
	}
	thresholds := []int{2, 1, 2}

	res, _ := comsoc.ApprovalSWF(prefs, thresholds)

	if res[1] != 2 {
		t.Errorf("erreur, résultat pour 1 devrait être 2, %d calculé", res[1])
	}
	if res[2] != 2 {
		t.Errorf("erreur, résultat pour 2 devrait être 2, %d calculé", res[2])
	}
	if res[3] != 1 {
		t.Errorf("erreur, résultat pour 3 devrait être 1, %d calculé", res[3])
	}
}

func TestApprovalSWFThresholdError(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 3, 2},
		{2, 3, 1},
	}
	thresholds := []int{2, 1}

	_, err := comsoc.ApprovalSWF(prefs, thresholds)
	if err == nil {
		t.Errorf("erreur, devrait lever une erreur : erreur : profile de taille différente de threshold (len(profile)==3 et len(threshold)==2)")
	}

}

func TestApprovalSCF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
	}
	thresholds := []int{2, 1, 2}

	res, err := comsoc.ApprovalSCF(prefs, thresholds)

	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("erreur, 1 devrait être la seule meilleure Alternative")
	}
}

func TestApprovalSCFMultiple(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 3, 2},
		{3, 2, 1},
		{2, 1, 3},
	}
	thresholds := []int{1, 1, 1}

	res, _ := comsoc.ApprovalSCF(prefs, thresholds)
	if len(res) != 3 {
		t.Errorf("erreur, les 3 Alternative devraient être à égalité")
	}
}

func TestCondorcetWinner(t *testing.T) {
	prefs1 := comsoc.Profile{
		{3, 2, 1},
		{1, 2, 3},
		{1, 2, 3},
	}

	prefs2 := comsoc.Profile{
		{1, 2, 3},
		{2, 3, 1},
		{3, 1, 2},
	}

	res1, _ := comsoc.CondorcetWinner(prefs1)
	res2, _ := comsoc.CondorcetWinner(prefs2)

	if len(res1) == 0 || res1[0] != 1 {
		t.Errorf("erreur, 1 devrait être la seule meilleure Alternative pour prefs1")
	}
	if len(res2) != 0 {
		t.Errorf("pas de meilleure Alternative pour prefs2")
	}
}

func TestCopelandSWF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := comsoc.CopelandSWF(prefs)

	if res[1] != 2 {
		t.Errorf("erreur, resultat pour 1 devrait être 2, %d calculé", res[1])
	}
	if res[2] != 0 {
		t.Errorf("erreur, résultat pour 2 devrait être 0, %d calculé", res[2])
	}
	if res[3] != -2 {
		t.Errorf("erreur, résultat pour 3 devrait être -2, %d calculé", res[3])
	}
}

func TestCopelandSCF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.CopelandSCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("erreur, 1 devrait être la seule meilleure Alternative")
	}
}

func TestSTV_SWF(t *testing.T) {
	prefs := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := comsoc.STV_SWF(prefs)

	if res[1] != 2 {
		t.Errorf("erreur, resultat pour 1 devrait être 2, %d calculé", res[1])
	}
	if res[2] != 0 {
		t.Errorf("erreur, résultat pour 2 devrait être 0, %d calculé", res[2])
	}
	if res[3] != 1 {
		t.Errorf("erreur, résultat pour 3 devrait être 1, %d calculé", res[3])
	}
}

func TestSTV_SCF(t *testing.T) {
	prefs1 := comsoc.Profile{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	prefs2 := comsoc.Profile{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},

		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},

		{4, 3, 1, 2},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
	}

	res1, err := comsoc.STV_SCF(prefs1)
	if err != nil {
		t.Error(err)
	}

	res2, err := comsoc.STV_SCF(prefs2)
	if err != nil {
		t.Error(err)
	}

	if len(res1) != 1 || res1[0] != 1 {
		t.Errorf("erreur, res1, 1 devrait être la seule meilleure Alternative")
	}
	if len(res2) != 1 || res2[0] != 1 {
		t.Errorf("erreur res2, 1 devrait être la seule meilleure Alternative")
	}
}

func TestTieBreakFactory(t *testing.T) {
	var alt comsoc.Alternative = 2
	orderedAlts := []comsoc.Alternative{1, 2, 3}
	alts := []comsoc.Alternative{3, 2}

	tieBreak := comsoc.TieBreakFactory(orderedAlts)
	res, err := tieBreak(alts)

	if err != nil {
		t.Error(err)
	}
	if res != alt {
		t.Errorf("erreur, résultat devrait être 2, mais est %d", res)
	}
}

func TestSWFFactory(t *testing.T) {
	orderedAlts1 := []comsoc.Alternative{1, 2, 3, 4}
	orderedAlts2 := []comsoc.Alternative{3, 1, 2, 4}
	prefs1 := comsoc.Profile{
		{1, 3, 2, 4},
		{1, 2, 3, 4},
		{2, 1, 3, 4},
		{2, 1, 3, 4},
	}
	prefs2 := comsoc.Profile{
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 1, 2, 4},
		{3, 1, 2, 4},
	}
	prefs3 := comsoc.Profile{
		{1, 2, 3, 4},
		{2, 1, 3, 4},
		{1, 2, 3, 4},
		{2, 1, 3, 4},
		{3, 1, 2, 4},
		{4, 1, 3, 2},
	}

	tieBreak1 := comsoc.TieBreakFactory(orderedAlts1)
	f1 := comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak1)
	res1, err1 := f1(prefs1)
	res2, err2 := f1(prefs2)

	tieBreak2 := comsoc.TieBreakFactory(orderedAlts2)
	f2 := comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak2)
	res3, err3 := f2(prefs3)

	if err1 != nil {
		t.Error(err1)
	}
	if res1[0] != 1 {
		t.Errorf("erreur, le premier resultat devrait être 1, %d calculé", res1[0])
	}
	if res1[1] != 2 {
		t.Errorf("erreur, le deuxième résultat devrait être 2, %d calculé", res1[1])
	}
	if res1[2] != 3 {
		t.Errorf("erreur, le troisième résultat devrait être 3, %d calculé", res1[2])
	}
	if res1[3] != 4 {
		t.Errorf("erreur, le troisième résultat devrait être 4, %d calculé", res1[3])
	}

	if err2 != nil {
		t.Error(err2)
	}
	if res2[0] != 3 {
		t.Errorf("erreur, le premier resultat devrait être 3, %d calculé", res2[0])
	}
	if res2[1] != 1 {
		t.Errorf("erreur, le deuxième résultat devrait être 1, %d calculé", res2[1])
	}
	if res2[2] != 2 {
		t.Errorf("erreur, le troisième résultat devrait être 2, %d calculé", res2[2])
	}
	if res2[3] != 4 {
		t.Errorf("erreur, le troisième résultat devrait être 4, %d calculé", res2[3])
	}

	if err3 != nil {
		t.Error(err3)
	}
	if res3[0] != 1 {
		t.Errorf("erreur, le premier resultat devrait être 3, %d calculé", res3[0])
	}
	if res3[1] != 2 {
		t.Errorf("erreur, le deuxième résultat devrait être 1, %d calculé", res3[1])
	}
	if res3[2] != 3 {
		t.Errorf("erreur, le troisième résultat devrait être 2, %d calculé", res3[2])
	}
	if res3[3] != 4 {
		t.Errorf("erreur, le troisième résultat devrait être 4, %d calculé", res3[3])
	}
}

func TestSCFFactory(t *testing.T) {
	orderedAlts1 := []comsoc.Alternative{1, 2, 3}
	orderedAlts2 := []comsoc.Alternative{2, 1}
	prefs1 := comsoc.Profile{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
		{2, 1, 3},
	}
	prefs2 := comsoc.Profile{
		{1, 2},
		{2, 1},
		{1, 2},
		{2, 1},
	}

	tieBreak1 := comsoc.TieBreakFactory(orderedAlts1)
	f1 := comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak1)
	res1, err1 := f1(prefs1)

	tieBreak2 := comsoc.TieBreakFactory(orderedAlts2)
	f2 := comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak2)
	res2, err2 := f2(prefs2)

	if err1 != nil {
		t.Error(err1)
	}
	if res1 != 1 {
		t.Errorf("erreur, resultat devrait être 1, %d calculé", res1)
	}

	if err2 != nil {
		t.Error(err2)
	}
	if res2 != 2 {
		t.Errorf("The result should be 2, found %d", res2)
	}
}
