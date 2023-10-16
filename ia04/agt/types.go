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
