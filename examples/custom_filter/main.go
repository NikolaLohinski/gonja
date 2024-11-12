package main

import (
	"encoding/base64"
	"os"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

var filterFuncB64Encode exec.FilterFunction = func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	return exec.AsValue(base64.StdEncoding.EncodeToString([]byte(in.String())))
}

func main() {
	gonja.DefaultEnvironment.Filters.Register("b64encode", filterFuncB64Encode)

	template, err := gonja.FromString("{{ var | b64encode }}")
	if err != nil {
		panic(err)
	}

	ctx := exec.NewContext(map[string]interface{}{
		"var": "gonja",
	})
	if err = template.Execute(os.Stdout, ctx); err != nil { // Prints Z29uamE= to stdout
		panic(err)
	}
}
