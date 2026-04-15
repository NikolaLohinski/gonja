// Package methods provides built-in method implementations for template types.
package methods

import "github.com/ardanlabs/gonja/exec"

var All = exec.Methods{
	Bool:  boolMethods,
	Str:   strMethods,
	Int:   intMethods,
	Float: floatMethods,
	Dict:  dictMethods,
	List:  listMethods,
}
