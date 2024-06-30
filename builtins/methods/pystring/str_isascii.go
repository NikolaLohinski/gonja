package pystring

// Return True if all characters in the string are decimal characters and there is at least one character, False otherwise. Decimal characters are those that can be used to form numbers in base 10, e.g. U+0660, ARABIC-INDIC DIGIT ZERO. Formally a decimal character is a character in the Unicode General Category “Nd”.
func IsASCII(s string) bool {
	if len(s) == 0 {
		return true
	}

	for _, char := range s {
		if char > 127 {
			return false
		}
	}
	return true
}

// Return True if all characters in the string are decimal characters and there is at least one character, False otherwise. Decimal characters are those that can be used to form numbers in base 10, e.g. U+0660, ARABIC-INDIC DIGIT ZERO. Formally a decimal character is a character in the Unicode General Category “Nd”.
func (pys PyString) IsASCII() bool {
	return IsASCII(string(pys))
}
