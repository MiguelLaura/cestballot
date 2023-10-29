package sequential

import "gitlab.utc.fr/mennynat/ia04-tp/utils"

// ---------------------------
//           FILL
// ---------------------------

func Fill[T any](tab []T, v T) {
	for idx := range tab {
		tab[idx] = v
	}
}

// ---------------------------
//          FOREACH
// ---------------------------

func ForEach[T any](tab []T, f func(T) T) {
	for idx, value := range tab {
		tab[idx] = f(value)
	}
}

// ---------------------------
//           COPY
// ---------------------------

func Copy[T any](src []T, dest []T) {
	copy(dest, src)
}

// ---------------------------
//           EQUAL
// ---------------------------

func Equal[T comparable](tab1 []T, tab2 []T) bool {
	if len(tab1) != len(tab2) {
		return false
	}

	for idx := range tab1 {
		if tab1[idx] != tab2[idx] {
			return false
		}
	}

	return true
}

// ---------------------------
//           FIND
// ---------------------------

func Find[T comparable](tab []T, f func(T) bool) (index int, val T) {
	index = -1

	for idx, value := range tab {
		if f(value) {
			index = idx
			val = value
			break
		}
	}

	return
}

// ---------------------------
//            MAP
// ---------------------------

func Map[T any](tab []T, f func(T) T) []T {
	mappedTab := make([]T, len(tab))

	for idx, value := range tab {
		mappedTab[idx] = f(value)
	}

	return mappedTab
}

// ---------------------------
//          REDUCE
// ---------------------------

func Reduce[T any](tab []T, init T, f func(T, T) T) T {
	acc := init

	for _, value := range tab {
		acc = f(acc, value)
	}

	return acc
}

// ---------------------------
//           EVERY
// ---------------------------

func Every[T comparable](tab []T, f func(T) bool) bool {
	for _, value := range tab {
		if !f(value) {
			return false
		}
	}

	return true
}

// ---------------------------
//            ANY
// ---------------------------

func Any[T comparable](tab []T, f func(T) bool) bool {
	for _, value := range tab {
		if f(value) {
			return true
		}
	}

	return false
}

// ---------------------------
//           SORT
// ---------------------------

// comp(v1, v2) function is true when v1 < v2
func Sort[T comparable](data []T, comp func(T, T) bool) {
	if len(data) > 1 {
		middleIdx := len(data) / 2
		Sort(data[:middleIdx], comp)
		Sort(data[middleIdx:], comp)
		Copy(utils.Merge(data[:middleIdx], data[middleIdx:], comp), data)
	}
}
