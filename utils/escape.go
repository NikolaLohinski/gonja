// Package utils provides utility functions for template processing.
package utils

import (
	"html"
	"strings"
)

func Escape(in string) string {
	output := html.EscapeString(in)
	output = strings.ReplaceAll(output, "'", "&#39;")
	return output
}
