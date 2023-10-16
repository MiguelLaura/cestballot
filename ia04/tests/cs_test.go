// version 2.0.0

//  RAJOUT DE COMSOC

package test_comsoc

import (
	"ia04/comsoc"
	"testing"
)

func TestCheckProfileAlternativeError(t *testing.T) {
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs1 := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	prefs2 := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
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
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.STV_SCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("erreur, 1 devrait être la seule meilleure Alternative")
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
	orderedAlts := []comsoc.Alternative{1, 2, 3}
	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
		{2, 1, 3},
	}

	tieBreak := comsoc.TieBreakFactory(orderedAlts)
	f := comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)
	res, err := f(prefs)

	if err != nil {
		t.Error(err)
	}
	if res[1] != 3 {
		t.Errorf("erreur, resultat pour 1 devrait être 3, %d calculé", res[1])
	}
	if res[2] != 2 {
		t.Errorf("erreur, résultat pour 2 devrait être 0, %d calculé", res[2])
	}
	if res[3] != 0 {
		t.Errorf("erreur, résultat pour 3 devrait être 1, %d calculé", res[3])
	}
}

func TestSCFFactory(t *testing.T) {
	orderedAlts := []comsoc.Alternative{1, 2, 3}
	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
		{2, 1, 3},
	}

	tieBreak := comsoc.TieBreakFactory(orderedAlts)
	f := comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)
	res, err := f(prefs)

	if err != nil {
		t.Error(err)
	}
	if res != 1 {
		t.Errorf("erreur, resultat devrait être 1, %d calculé", res)
	}
}
