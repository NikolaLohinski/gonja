package pystring

import "unicode"

// Return True if all characters in the string are digits and there is at least one character, False otherwise. Digits include decimal characters and digits that need special handling, such as the compatibility superscript digits. This covers digits which cannot be used to form numbers in base 10, like the Kharosthi numbers. Formally, a digit is a character that has the property value Numeric_Type=Digit or Numeric_Type=Decimal.
func IsLower(s string) bool {
	hasCased := false
	for _, char := range s {
		if unicode.IsUpper(char) {
			return false
		} else if unicode.IsLower(char) {
			hasCased = true
		}
	}
	return hasCased
}

// Return True if all characters in the string are digits and there is at least one character, False otherwise. Digits include decimal characters and digits that need special handling, such as the compatibility superscript digits. This covers digits which cannot be used to form numbers in base 10, like the Kharosthi numbers. Formally, a digit is a character that has the property value Numeric_Type=Digit or Numeric_Type=Decimal.
func (pys PyString) IsLower() bool {
	return IsLower(string(pys))
}
