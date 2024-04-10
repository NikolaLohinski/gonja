package pystring

// Alias for FormatWithDialect
func FormatMapWithDialect(d Dialect, s string, vargs []any, kwarg map[string]any) (string, error) {
	return FormatWithDialect(d, s, vargs, kwarg)
}

// Alias for Format
func FormatMap(s string, vargs []any, kwarg map[string]any) (string, error) {
	return Format(s, vargs, kwarg)
}

// Alias for Format
func (s PyString) FormatMap(vargs []any, kwarg map[string]any) (PyString, error) {
	return s.FormatMap(vargs, kwarg)
}

// Alias for FormatMapWithDialect
func (s PyString) FormatMapWithDialect(d Dialect, vargs []any, kwarg map[string]any) (PyString, error) {
	return s.FormatWithDialect(d, vargs, kwarg)
}
