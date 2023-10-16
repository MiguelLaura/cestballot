package ballotagent

import (
	"errors"
	"fmt"
	"ia04/comsoc"
)

func checkTieBreak(alts int, tieBreak []comsoc.Alternative) error {
	if len(tieBreak) != alts {
		err := "erreur : pas le bon nombre d'alternatives dans tieBreak"
		return errors.New(err)
	}
	check := make(map[int]int)
	for i := 1; i < alts+1; i++ {
		check[i] = 0
	}
	for _, val := range tieBreak {
		if val < 1 || int(val) > alts {
			err := "erreur : valeur incorrecte dans tieBreak (pas dans le bon range) :" + fmt.Sprint(val)
			return errors.New(err)
		}
		check[int(val)] += 1
	}
	for key, val := range check {
		if val != 1 {
			err := "erreur : valeur incorrecte dans tieBreak : " + fmt.Sprint(key) + " apparaît " + fmt.Sprint(val) + " fois."
			return errors.New(err)
		}
	}
	return nil
}
