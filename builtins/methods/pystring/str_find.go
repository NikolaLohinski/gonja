package pystring

import (
	"strings"
)

// Return True if the string ends with the specified suffix, otherwise return False.
// suffix can also be a tuple of suffixes to look for. With optional start, test
// beginning at that position. With optional end, stop comparing at that position.
func Find(s, subStr string, start, end *int) int {
	s, _ = Idx(s, start, end)
	if s == "" {
		return 0
	}
	return strings.Index(s, subStr)
}

// Return True if the string ends with the specified suffix, otherwise return False.
// suffix can also be a tuple of suffixes to look for. With optional start, test
// beginning at that position. With optional end, stop comparing at that position.
func (pys PyString) Find(substr PyString, start, end *int) int {
	return Find(string(pys), string(substr), start, end)
}
