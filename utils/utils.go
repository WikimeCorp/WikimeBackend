package utils

import (
	"strings"

	crand "crypto/rand"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

const (
	// Length is 64, in order to diversify the strings as much as possible without compromising performance.
	// Since the mask can cover 64 characters, we do not need to check that after applying the mask,
	// the index will be larger than the size of the array.
	randomAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
)

const (
	maxBufSize int = 128
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

func FastRandomString(size int) string {
	ansStr := strings.Builder{}
	ansStr.Grow(size)

	bufSize := maxBufSize
	if size < bufSize {
		bufSize = size
	}
	buffer := make([]byte, bufSize)

	alphabetLen := len(randomAlphabet)
	maxByte := 255 - (256 % alphabetLen)
	i := 0
outer:
	for {
		if _, err := crand.Read(buffer[:bufSize]); err != nil {
			panic("WTF")
		}

		for _, el := range buffer[:bufSize] {
			intEl := int(el)
			if intEl > maxByte {
				continue
			}
			ansStr.WriteByte(randomAlphabet[intEl%alphabetLen])
			i++
			if i == size {
				break outer
			}
		}
		if ansStr.Len() < bufSize {
			bufSize = ansStr.Len()
		}
	}

	return ansStr.String()
}
