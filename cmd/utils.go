package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gitlab.utc.fr/mennynat/ia04-tp/comsoc"
)

// Permet d'acquérir les votants en ligne de commande

type VotersFlag struct {
	voters []string
}

func (vf *VotersFlag) String() string {
	return strings.Join(vf.voters, ",")
}

func (vf *VotersFlag) Set(s string) error {
	if s[len(s)-1] == ',' {
		log.Fatalf("Format du flag des voters incorrect")
	}
	vf.voters = strings.Split(s, ",")
	return nil
}

func (vf *VotersFlag) GetVoters() []string {
	if vf.voters == nil {
		return []string{"ag_id1"}
	}
	return vf.voters
}

// Permet d'acquérir les options en ligne de commande

type OptFlag struct {
	opts []int
}

func (of *OptFlag) String() string {
	return fmt.Sprintf("%#v", of.opts)
}

func (of *OptFlag) Set(s string) error {
	optsStr := strings.Split(s, ",")
	opts := make([]int, len(optsStr))

	for optIdx, optStr := range optsStr {
		optConv, err := strconv.Atoi(optStr)

		if err != nil {
			log.Fatal("Une des option donnée n'est pas un entier")
		}

		opts[optIdx] = optConv
	}

	of.opts = opts
	return nil
}

func (of *OptFlag) GetOpts() []int {
	return of.opts
}

// Permet d'acquérir les alternatives en ligne de commande

type AltFlag struct {
	alternatives []comsoc.Alternative
}

func (af *AltFlag) String() string {
	return fmt.Sprintf("%#v", af.alternatives)
}

func (af *AltFlag) Set(s string) error {
	altsStr := strings.Split(s, ",")
	alts := make([]comsoc.Alternative, len(altsStr))

	for altIdx, altStr := range altsStr {
		altConv, err := strconv.Atoi(altStr)

		if err != nil {
			log.Fatal("Une des alternative donnée n'est pas un entier")
		}

		alts[altIdx] = comsoc.Alternative(altConv)
	}

	af.alternatives = alts
	return nil
}

func (af *AltFlag) GetAlts() []comsoc.Alternative {
	if af.alternatives == nil {
		return []comsoc.Alternative{4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10}
	}
	return af.alternatives
}
