package pystring

import "strings"

// Split the string at the first occurrence of sep, and return a 3-tuple
// containing the part before the separator, the separator itself, and the part
// after the separator. If the separator is not found, return a 3-tuple containing
// the string itself, followed by two empty strings.
func Partition(s string, delim string) (string, string, string) {
	i := strings.Index(s, delim)
	if i == -1 || delim == "" {
		return s, "", ""
	}
	return s[:i], delim, s[i+len(delim):]
}

// Split the string at the first occurrence of sep, and return a 3-tuple
// containing the part before the separator, the separator itself, and the part
// after the separator. If the separator is not found, return a 3-tuple containing
// the string itself, followed by two empty strings.
func (pys PyString) Partition(delim string) (PyString, PyString, PyString) {
	r1, r2, r3 := Partition(string(pys), delim)
	return PyString(r1), PyString(r2), PyString(r3)
}
