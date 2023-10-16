package agt

import (
	"fmt"
	"ia04/comsoc"
)

var C chan []comsoc.Alternative

type AgentID int

type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a comsoc.Alternative, b comsoc.Alternative) bool
	Start()
}

type Agent struct {
	ID    AgentID
	Name  string
	Prefs []comsoc.Alternative
}

func (ag Agent) Equal(ag2 Agent) bool {
	return (ag.ID == ag2.ID && ag.Name == ag2.Name)
}

func (ag Agent) DeepEqual(ag2 Agent) bool {
	if !ag.Equal(ag2) {
		return false
	}
	for idx, val := range ag.Prefs {
		if ag2.Prefs[idx] != val {
			return false
		}
	}
	return true
}

func (ag Agent) Clone() Agent {
	res := Agent{ID: ag.ID, Name: ag.Name}
	var prefs []comsoc.Alternative
	copy(prefs, ag.Prefs)
	res.Prefs = prefs
	return res
}

func (ag Agent) String() string {
	return fmt.Sprintf("l'agent à l'ID %d, au nom %s, et avec les préférences %d}", ag.ID, ag.Name, ag.Prefs)
}

func (ag Agent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	for _, val := range ag.Prefs {
		if val == a {
			return true
		} else if val == b {
			return false
		}
	}
	return false
}

func (ag Agent) Start() {
	C <- ag.Prefs
	fmt.Print("Un agent vient de voter :", ag.String(), "\n")
}
