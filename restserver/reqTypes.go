// Package restserver handling the REST server allowing to vote
package restserver

import "gitlab.utc.fr/mennynat/ia04-tp/comsoc"

/*
	General status code
*/

const HTTP_VERB_INCORRECT = 405

/*
	Request status codes
*/

const VOTE_CREATED = 201
const BAD_REQUEST = 400
const NOT_IMPL = 501

const VOTE_TAKEN = 200
const VOTE_ALREADY_DONE = 403
const DEADLINE_OVER = 503

const RESULT_OBTAINED = 200
const NOT_FOUND = 404
const TOO_EARLY = 425

/*
	JSON messages
*/

// The new-ballot request.
type NewBallotRequest struct {
	Rule     string               `json:"rule"`      // The rule used to vote
	Deadline string               `json:"deadline"`  // The deadline in RFC3339 format
	Voters   []string             `json:"voter-ids"` // List of the voters ID allowed to vote
	Alts     int                  `json:"#alts"`     // Number of alternatives
	TieBreak []comsoc.Alternative `json:"tie-break"` // The tiebreak use when scf cannot decide the best alternative
}

// The The new-ballot response.
type NewBallotResponse struct {
	Id string `json:"ballot-id"` // The ID of the created ballot
}

// The vote request.
type VoteRequest struct {
	Agent   string               `json:"agent-id"`          // The voter's ID
	Ballot  string               `json:"ballot-id"`         // The ballot's ID
	Prefs   []comsoc.Alternative `json:"prefs"`             // The preferences of the voter
	Options []int                `json:"options,omitempty"` // Optional values used for specific voting methods
}

// The result request.
type ResultRequest struct {
	Ballot string `json:"ballot-id"` // The ballot's ID
}

// The The new-ballot response.
type ResultResponse struct {
	Winner  comsoc.Alternative   `json:"winner"`            // The winner of the vote
	Ranking []comsoc.Alternative `json:"ranking,omitempty"` // A ranking
}
