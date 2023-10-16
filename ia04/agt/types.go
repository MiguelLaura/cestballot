package agt

import "ia04/comsoc"

type RequestBallot struct {
	Rule     string               `json:"rule"`
	Deadline string               `json:"deadline"`
	VoterIds []string             `json:"voter-ids"`
	Alts     int                  `json:"#alts"`
	TieBreak []comsoc.Alternative `json:"tie-break"`
}

type ResponseBallot struct {
	BallotId string `json:"ballot-id"`
}

// options est facultatif et permet de passer des renseignements supplémentaires
// (par exemple le seuil d'acceptation en approval)
type RequestVoter struct {
	AgentId  string               `json:"agent-id"`
	BallotId string               `json:"ballot-id"`
	Prefs    []comsoc.Alternative `json:"prefs"`
	Options  []int                `json:"options"`
}

type RequestResult struct {
	BallotId string `json:"ballot-id"`
}

// [A FAIRE] changer type ranking
type ResponseResult struct {
	Winner  comsoc.Alternative `json:"winner"`
	Ranking comsoc.Count       `json:"ranking"`
}
