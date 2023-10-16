// package agt

// import "ia04/comsoc"

// type RequestInit struct {
// 	Method     string               `json:"method"`
// 	Candidates []comsoc.Alternative `json:"candidates"`
// }

// type RequestVoter struct {
// 	Prefs []comsoc.Alternative `json:"prefs"`
// }

// type Response struct {
// 	Winner  comsoc.Alternative `json:"winner"`
// 	Ranking comsoc.Count       `json:"ranking"`
// }

package agt

import "ia04/comsoc"

type RequestBallot struct {
	Rule     string               `json:"rule"`
	Deadline string               `json:"deadline"`
	VoterIds []string             `json:"voters-ids"`
	Alt      int                  `json:"#alt"`
	TieBreak []comsoc.Alternative `json:"tie-break"`
}
