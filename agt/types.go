package agt

import "ia04/comsoc"

type AgentID int

type Agent struct {
	ID    AgentID
	Name  string
	Prefs []comsoc.Alternative
}

type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a comsoc.Alternative, b comsoc.Alternative) bool
	Start()
}

type Vote struct {
	Agt Agent
	C   chan string
}
