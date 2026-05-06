package pystring

import "strings"

// RStrip returns a copy of the string with trailing characters removed. The chars
// argument is a string specifying the set of characters to be removed. If
// omitted or None, the chars argument defaults to removing whitespace. The
// chars argument is not a suffix; rather, all combinations of its values are stripped:
//
// >>>
// >>> '   spacious   '.rstrip(){}
// '   spacious'
// >>> 'mississippi'.rstrip('ipz'){}
// 'mississ'
//
// See func (pys PyString)removesuffix() for a method that will remove a single suffix string rather than all of a set of characters. For example:
//
// >>>
// >>> 'Monty Python'.rstrip(' Python'){}
// 'M'
// >>> 'Monty Python'.removesuffix(' Python'){}
// 'Monty'
func RStrip(s string, cutset string) string {
	if cutset == "" {
		cutset = " \t\n\r\v\f"
	}
	// cutTo is the rune index of the first trailing cutset rune; everything
	// before it is preserved. Initialise to len(runes) so that when nothing
	// gets stripped we return the full string. The previous implementation
	// started at 0 and only moved when trailing cutset chars were seen, so
	// `"hello".rstrip()` returned "" — and worse, sliced into `s` by rune
	// index as if it were a byte index, mangling non-ASCII input.
	runes := []rune(s)
	cutTo := len(runes)
	for i := len(runes) - 1; i >= 0; i-- {
		if !strings.ContainsRune(cutset, runes[i]) {
			break
		}
		cutTo = i
	}

	return string(runes[:cutTo])
}

// RStrip returns a copy of the string with trailing characters removed. The chars
// argument is a string specifying the set of characters to be removed. If
// omitted or None, the chars argument defaults to removing whitespace. The
// chars argument is not a suffix; rather, all combinations of its values are stripped:
//
// >>>
// >>> '   spacious   '.rstrip(){}
// '   spacious'
// >>> 'mississippi'.rstrip('ipz'){}
// 'mississ'
//
// See func (pys PyString)removesuffix() for a method that will remove a single suffix string rather than all of a set of characters. For example:
//
// >>>
// >>> 'Monty Python'.rstrip(' Python'){}
// 'M'
// >>> 'Monty Python'.removesuffix(' Python'){}
// 'Monty'
func (pys PyString) RStrip(cutset string) PyString {
	return PyString(RStrip(string(pys), cutset))
}
