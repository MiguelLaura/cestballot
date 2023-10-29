package concurrent

import (
	"sync"

	"gitlab.utc.fr/mennynat/ia04-tp/utils"
	"gitlab.utc.fr/mennynat/ia04-tp/utils/sequential"
)

// ---------------------------
//           UTIL
// ---------------------------

type Slice[T any] struct {
	slice    []T // The subslice
	startIdx int // The startindex of the subslice in the slice
	endIdx   int // the endindex of the subslice in the slice
}

// Creates subslices of a slice and send it one by one through a channel
func slicer[T any](tab []T) <-chan Slice[T] {
	sizeSlices := len(tab) / 10
	if sizeSlices == 0 {
		sizeSlices = 1
	}
	chnl := make(chan Slice[T])

	go func() {
		defer close(chnl)

		for startIdx := 0; startIdx < len(tab); startIdx += sizeSlices {
			endIdx := min(len(tab), startIdx+sizeSlices)
			chnl <- Slice[T]{tab[startIdx:endIdx], startIdx, endIdx}
		}
	}()

	return chnl
}

// Creates subslices of the same size from two slices and send them one by one through a channel
// Returns : channel of []Slice containing [subsliceOfTab1, subsliceOfTab2]
func doubleSlicer[T any](tab1 []T, tab2 []T) <-chan []Slice[T] {
	var slices []Slice[T]
	minLength := min(len(tab1), len(tab2))
	sizeSlices := minLength / 10
	if sizeSlices == 0 {
		sizeSlices = 1
	}
	chnl := make(chan []Slice[T])

	go func() {
		defer close(chnl)

		for startIdx := 0; startIdx < minLength; startIdx += sizeSlices {
			endIdx := min(minLength, startIdx+sizeSlices)
			slices = make([]Slice[T], 2)
			slices[0] = Slice[T]{tab1[startIdx:endIdx], startIdx, endIdx}
			slices[1] = Slice[T]{tab2[startIdx:endIdx], startIdx, endIdx}
			chnl <- slices
		}
	}()

	return chnl
}

// Reads a channel in a go routine until it closes
func splitItOut[T any](chnl chan T) {
	go func() {
		for range chnl {
		}
	}()
}

// ---------------------------
//           FILL
// ---------------------------

func Fill[T any](tab []T, v T) {
	var wg sync.WaitGroup

	for sl := range slicer(tab) {
		wg.Add(1)
		go func(sl []T) {
			defer wg.Done()
			sequential.Fill(sl, v)
		}(sl.slice)
	}

	wg.Wait()
}

// ---------------------------
//          FOREACH
// ---------------------------

func ForEach[T any](tab []T, f func(T) T) {
	var wg sync.WaitGroup

	for sl := range slicer(tab) {
		wg.Add(1)
		go func(sl []T) {
			defer wg.Done()
			sequential.ForEach(sl, f)
		}(sl.slice)
	}

	wg.Wait()
}

// ---------------------------
//           COPY
// ---------------------------

func Copy[T any](src []T, dest []T) {
	var wg sync.WaitGroup

	for sls := range doubleSlicer(src, dest) {
		wg.Add(1)
		go func(sliceSrc []T, sliceDest []T) {
			defer wg.Done()
			sequential.Copy(sliceSrc, sliceDest)
		}(sls[0].slice, sls[1].slice)
	}

	wg.Wait()
}

// ---------------------------
//           EQUAL
// ---------------------------

func Equal[T comparable](tab1 []T, tab2 []T) bool {
	if len(tab1) != len(tab2) {
		return false
	}

	chnl := make(chan bool)

	// Go routine that checks if the subslices are equal;
	//  if one of them is not equal, sends the value false through the channel
	go func() {
		var wg sync.WaitGroup
		defer close(chnl)

		for sls := range doubleSlicer(tab1, tab2) {
			wg.Add(1)
			go func(sliceTab1 []T, sliceTab2 []T) {
				if !sequential.Equal(sliceTab1, sliceTab2) {
					chnl <- false
				}
				wg.Done()
			}(sls[0].slice, sls[1].slice)
		}

		wg.Wait()
	}()

	// If the above goroutine sends something through the channel, it means that one
	//  of the subslice is not equal to the other
	for range chnl {
		splitItOut(chnl)
		return false
	}

	return true

}

// ---------------------------
//           FIND
// ---------------------------

func Find[T comparable](tab []T, f func(T) bool) (index int, val T) {
	index = -1

	// Create a channel transmitting the find indexes of the tab slice (when found)
	chnl := make(chan int)

	go func() {
		var wg sync.WaitGroup
		defer close(chnl)

		for sl := range slicer(tab) {
			wg.Add(1)
			go func(slice Slice[T]) {
				idx, _ := sequential.Find(slice.slice, f)
				// If the value if found in the slice, its index is sent through the channel
				if idx != -1 {
					chnl <- slice.startIdx + idx
				}
				wg.Done()
			}(sl)
		}

		wg.Wait()
	}()

	// When a new value if found, the channel tells us its index in the slice
	for idx := range chnl {
		if idx < index || index == -1 {
			index = idx
			val = tab[index]
		}
	}

	return
}

// ---------------------------
//            MAP
// ---------------------------

func Map[T any](tab []T, f func(T) T) []T {
	mappedTab := make([]T, len(tab))
	Copy(tab, mappedTab)
	ForEach(mappedTab, f)
	return mappedTab
}

// ---------------------------
//          REDUCE
// ---------------------------

func Reduce[T any](tab []T, init T, f func(T, T) T) T {
	acc := init
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for sl := range slicer(tab) {
		if len(sl.slice) == 1 {
			mutex.Lock()
			acc = f(acc, sl.slice[0])
			mutex.Unlock()
		} else {
			wg.Add(1)
			go func(slice []T) {
				defer wg.Done()
				seqAcc := sequential.Reduce(slice[1:], slice[0], f)
				mutex.Lock()
				acc = f(acc, seqAcc)
				mutex.Unlock()
			}(sl.slice)
		}
	}

	wg.Wait()

	return acc
}

// ---------------------------
//           EVERY
// ---------------------------

func Every[T comparable](tab []T, f func(T) bool) bool {
	chnl := make(chan bool)

	go func() {
		var wg sync.WaitGroup
		defer close(chnl)

		for sl := range slicer(tab) {
			wg.Add(1)
			go func(slice []T) {
				defer wg.Done()
				if !sequential.Every(slice, f) {
					chnl <- false
				}
			}(sl.slice)
		}

		wg.Wait()
	}()

	// If a go routine sends something through the channel, it means that it has found at least one element
	//  not satisfying f
	for range chnl {
		splitItOut(chnl)
		return false
	}
	return true
}

// ---------------------------
//            ANY
// ---------------------------

func Any[T comparable](tab []T, f func(T) bool) bool {
	chnl := make(chan bool)

	go func() {
		var wg sync.WaitGroup
		defer close(chnl)

		for sl := range slicer(tab) {
			wg.Add(1)
			go func(slice []T) {
				defer wg.Done()
				if sequential.Any(slice, f) {
					chnl <- true
				}
			}(sl.slice)
		}

		wg.Wait()
	}()

	// If a go routine sends something through the channel, it means that it has found at least one element
	//  satisfying f
	for range chnl {
		splitItOut(chnl)
		return true
	}
	return false
}

// ---------------------------
//           SORT
// ---------------------------

// Sort recursively the tab
// comp(v1, v2) function is true when v1 < v2
func sortRec[T comparable](tab []T, sem chan uint8, comp func(T, T) bool) {
	if len(tab) < 2 {
		return
	}

	var wg sync.WaitGroup
	middleIdx := len(tab) / 2

	// Sorts the left tab (i.e. after middleIdx)
	select {
	case sem <- 0:
		wg.Add(1)
		go func() {
			defer wg.Done()
			sortRec(tab[:middleIdx], sem, comp)
			<-sem
		}()
	default:
		sequential.Sort(tab[:middleIdx], comp)
	}

	// Sorts the right tab (i.e. before middleIdx)
	select {
	case sem <- 0:
		wg.Add(1)
		go func() {
			defer wg.Done()
			sortRec(tab[middleIdx:], sem, comp)
			<-sem
		}()
	default:
		sequential.Sort(tab[middleIdx:], comp)
	}

	wg.Wait()
	Copy(utils.Merge(tab[:middleIdx], tab[middleIdx:], comp), tab)
}

// Merge sort on the tab
// comp(v1, v2) function is true when v1 < v2
func Sort[T comparable](tab []T, comp func(T, T) bool) {
	// Creates a semaphore that allow up to 8 concurrent operations
	//  Used to limit the subslices sorted by go routines or sequentially
	//  otherwise causes the program to be killed having to many threads
	sem := make(chan uint8, 8)
	sortRec(tab, sem, comp)
}
