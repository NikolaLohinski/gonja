package methods

import "github.com/nikolalohinski/gonja/v2/exec"

var All = exec.Methods{
	Bool:  boolMethods,
	Str:   strMethods,
	Int:   intMethods,
	Float: floatMethods,
	Dict:  dictMethods,
	List:  listMethods,
}
