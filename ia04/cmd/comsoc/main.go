package main

import (
	"fmt"
	"ia04/agt"
	"ia04/comsoc"
	"time"
)

func main() {
	// var alt1 comsoc.Alternative = 2
	// var alt2 comsoc.Alternative = 7
	// var alt3 comsoc.Alternative = 1
	// count := make(comsoc.Count)
	// count[alt1] = 456
	// count[alt2] = 0
	// profile := [...]comsoc.Alternative{alt1, alt2}
	// profile1 := [...]comsoc.Alternative{alt2, alt1}
	// profile2 := [...]comsoc.Alternative{alt2, alt1}
	// var prefs comsoc.Profile
	// prefs = append(prefs, profile[:])
	// prefs = append(prefs, profile1[:])
	// prefs = append(prefs, profile2[:])
	// thresholds := [...]int{1, 1, 7}
	orderedAlts := []comsoc.Alternative{1, 2, 3}
	// prefs := [][]comsoc.Alternative{
	// 	{1, 2, 3},
	// 	{1, 2, 3},
	// 	{3, 2, 1},
	// }
	// thresholds := [...]int{1, 1, 7, 6}

	// fmt.Println(comsoc.CondorcetWinner(prefs))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// tieBreak := comsoc.TieBreakFactory(orderedAlts)
	// f := comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)
	// res, _ := f(prefs)
	// fmt.Println(res)
	agt.C = make(chan []comsoc.Alternative)
	ag1 := agt.Agent{ID: 1, Name: "Anna", Prefs: []comsoc.Alternative{1, 2, 3}}
	ag2 := agt.Agent{ID: 2, Name: "Bastien", Prefs: []comsoc.Alternative{2, 1, 3}}
	go ag1.Start()
	go ag2.Start()
	var profile [][]comsoc.Alternative
	profile = append(profile, <-agt.C)
	profile = append(profile, <-agt.C)
	time.Sleep(3 * time.Second)
	fmt.Print(profile, "\n")
	tieBreak := comsoc.TieBreakFactory(orderedAlts)
	f := comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)
	res, _ := f(profile)
	fmt.Println(res)

}
