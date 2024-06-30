package pystring

import "unicode"

// Return True if all characters in the string are digits and there is at least
// one character, False otherwise. Digits include decimal characters and digits that
// need special handling, such as the compatibility superscript digits. This covers
// digits which cannot be used to form numbers in base 10, like the Kharosthi numbers.
// Formally, a digit is a character that has the property value Numeric_Type=Digit
// or Numeric_Type=Decimal.
func IsNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, char := range s {
		if !unicode.IsDigit(char) && !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}

// Return True if all characters in the string are digits and there is at least
// one character, False otherwise. Digits include decimal characters and digits that
// need special handling, such as the compatibility superscript digits. This covers
// digits which cannot be used to form numbers in base 10, like the Kharosthi numbers.
// Formally, a digit is a character that has the property value Numeric_Type=Digit
// or Numeric_Type=Decimal.
func (pys PyString) IsNumeric() bool {
	return IsNumeric(string(pys))
}
