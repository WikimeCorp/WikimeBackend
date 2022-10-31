package utils

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

const (
	// Length is 63, in order to diversify the strings as much as possible without compromising performance.
	// Since the mask can cover 63 characters, we do not need to check that after applying the mask,
	// the index will be larger than the size of the array.
	randomAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!"
	letterIdxBits  = 6 // it takes 6 bits to encode the alphabet
	letterIdxMask  = 1<<letterIdxBits - 1
	letterIdxMax   = letterIdxMask / letterIdxBits
)

func InsertInSorted[T constraints.Ordered](to []T, el T) []T {
	var tmp T
	to = append(to, tmp)

	pos, _ := BinarySearch(to[:len(to)-1], el)
	copy(to[pos+1:], to[pos:])
	to[pos] = el
	return to

}

func BinarySearch[T constraints.Ordered](sl []T, el T) (int, bool) {
	return slices.BinarySearch(sl, el)
}
