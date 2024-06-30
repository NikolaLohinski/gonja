package pystring

import (
	"strings"
	"unicode"
)

// Return a copy of the string with all the cased characters
// converted to uppercase. Note that s.upper().isupper() might be False
// if s contains uncased characters or if the Unicode category of the resulting
// character(s) is not “Lu” (Letter, uppercase), but e.g. “Lt” (Letter, titlecase).
//
// The uppercasing algorithm used is described in section 3.13
// 'Default Case Folding' of the Unicode Standard.
// words on spaces only.
func Upper(s string) string {
	var res strings.Builder
	for _, char := range s {
		res.WriteRune(unicode.ToUpper(char))
	}
	return res.String()
}

// Return a copy of the string with all the cased characters
// converted to uppercase. Note that s.upper().isupper() might be False
// if s contains uncased characters or if the Unicode category of the resulting
// character(s) is not “Lu” (Letter, uppercase), but e.g. “Lt” (Letter, titlecase).
//
// The uppercasing algorithm used is described in section 3.13
// 'Default Case Folding' of the Unicode Standard.
// words on spaces only.
func (pys PyString) Upper() PyString {
	return PyString(Upper(string(pys)))
}
