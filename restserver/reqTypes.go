package restserver

import "ia04/comsoc"

type NewBallotRequest struct {
	Rule     string               `json:"rule"`
	Deadline string               `json:"deadline"`
	Voters   []string             `json:"voter-ids"`
	Alts     int                  `json:"#alts"`
	TieBreak []comsoc.Alternative `json:"tie-break"`
}

type NewBallotResponse struct {
	Id string `json:"ballot-id"`
}

type VoteRequest struct {
	Agent   string               `json:"agent-id"`
	Vote    string               `json:"vote-id"`
	Prefs   []comsoc.Alternative `json:"prefs"`
	Options []int                `json:"options"`
}

type ResultRequest struct {
	Ballot string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking"`
}
