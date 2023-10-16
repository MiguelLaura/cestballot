package restserver

import "ia04/comsoc"

const BAD_REQUEST = 400
const NOT_IMPL = 501
const METH_NOT_IMPL = 405

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

const VOTE_TAKEN = 200
const VOTE_ALREADY_DONE = 403
const DEADLINE_OVER = 503

type VoteRequest struct {
	Agent   string               `json:"agent-id"`
	Ballot  string               `json:"ballot-id"`
	Prefs   []comsoc.Alternative `json:"prefs"`
	Options []int                `json:"options,omitempty"`
}

const TOO_EARLY = 425

type ResultRequest struct {
	Ballot string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking,omitempty"`
}
