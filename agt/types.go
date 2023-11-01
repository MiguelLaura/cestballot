// Package agt contains the different agents available in the project (voters and ballot)
package agt

import "gitlab.utc.fr/mennynat/ia04-tp/comsoc"

// The type of the Agent identifier.
type AgentID int

// A simple voting agent.
type Agent struct {
	ID    AgentID              // The ID of the agent; must be unique
	Name  string               // The Name of the agent
	Prefs []comsoc.Alternative // The preferences of the agent ordered from the most to the least preferred.
}

// An interface to describe the methods of a voting agent
type AgentI interface {
	Equal(ag AgentI) bool                                    // <tt>true</tt> if the agents have the save ID and Name, <tt>false</tt> else
	DeepEqual(ag AgentI) bool                                // <tt>true</tt> if the agents are equals even regarding their preferences, <tt>false</tt> else
	Clone() AgentI                                           // Creates a duplicate of the current agent
	String() string                                          // Gives a string representing the agent
	Prefers(a comsoc.Alternative, b comsoc.Alternative) bool // <tt>true</tt> if the agent prefers the alternative a to b, <tt>false</tt> else
	Start()                                                  // Starts the agent so as to it can vote
}
