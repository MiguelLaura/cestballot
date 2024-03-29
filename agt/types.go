package agt

import "gitlab.utc.fr/mennynat/ia04-tp/comsoc"

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
	Options  []int                `json:"options,omitempty"`
}

type RequestResult struct {
	BallotId string `json:"ballot-id"`
}

type ResponseResult struct {
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking,omitempty"`
}
