// Package comsoc handles the voting methods.
//
// Contains the types used in comsoc.
package comsoc

// An alternative.
type Alternative int

// A profile.
type Profile [][]Alternative

// The number of votes received for each alternative.
type Count map[Alternative]int
