package utils

// Merges two slices in a certain order (determined by comp function)
// comp(v1, v2) function is true when v1 < v2
func Merge[T comparable](leftTab []T, rightTab []T, comp func(T, T) bool) (result []T) {
	result = make([]T, len(leftTab)+len(rightTab))
	leftIndex, rightIndex := 0, 0

	for index := 0; index < len(result); index++ {
		switch {
		// If we're out of bound on the left tab, just go through the right one
		case leftIndex >= len(leftTab):
			result[index] = rightTab[rightIndex]
			rightIndex++
		// If we're out of bound on the right tab, just go through the left one
		case rightIndex >= len(rightTab):
			result[index] = leftTab[leftIndex]
			leftIndex++
		// If the value on the left tab is < than the value on the right tab, it's added to the result
		case comp(leftTab[leftIndex], rightTab[rightIndex]):
			result[index] = leftTab[leftIndex]
			leftIndex++
		// Else, adds the element in the right tab
		default:
			result[index] = rightTab[rightIndex]
			rightIndex++
		}
	}
	return
}
