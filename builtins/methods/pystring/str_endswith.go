package pystring

import (
	"strings"
)

// Return True if the string ends with the specified suffix, otherwise return False.
// suffix can also be a tuple of suffixes to look for. With optional start, test
// beginning at that position. With optional end, stop comparing at that position.
func EndsWith(s, subStr string, start, end *int) bool {
	s, err := Idx(s, start, end)
	if err != nil {
		return false
	}
	if s == "" {
		return true
	}
	return strings.HasSuffix(s, subStr)
}

// Return True if the string ends with the specified suffix, otherwise return False.
// suffix can also be a tuple of suffixes to look for. With optional start, test
// beginning at that position. With optional end, stop comparing at that position.
func (pys PyString) EndsWith(substr PyString, start, end *int) bool {
	return EndsWith(string(pys), string(substr), start, end)
}
