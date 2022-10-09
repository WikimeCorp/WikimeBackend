package utils

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func InsertInSorted[T constraints.Ordered](to []T, el T) []T {
	var tmp T
	to = append(to, tmp)

	pos, _ := BinarySearch(to, el)
	copy(to[pos+1:], to[pos:])
	to[pos] = el
	return to

}

func BinarySearch[T constraints.Ordered](sl []T, el T) (int, bool) {
	return slices.BinarySearch(sl, el)
}
