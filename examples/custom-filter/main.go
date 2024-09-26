package main

import (
	"bytes"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/pkg/errors"
)

// Define a custom config. The StartString and EndString are different from default
var CustomConfig = config.Config{
	BlockStartString:    "{%",
	BlockEndString:      "%}",
	VariableStartString: "{$",
	VariableEndString:   "$}",
	CommentStartString:  "{#",
	CommentEndString:    "#}",
	AutoEscape:          false,
	StrictUndefined:     false,
	TrimBlocks:          true,
	LeftStripBlocks:     true,
}

// Define custom filter functions. Take a look at the original one as inspiration https://github.com/NikolaLohinski/gonja/blob/master/builtins/filters.go
var filterFuncB64Encode exec.FilterFunction = func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{{Name: "wrap", Default: nil}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'to_yaml'"))
	}
	// wrap is unsupported in golang, try to implement it later on
	o := b64.StdEncoding.EncodeToString([]byte(in.String()))
	return exec.AsValue(o)
}

var filterFuncB64Decode exec.FilterFunction = func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{{Name: "wrap", Default: nil}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'to_yaml'"))
	}
	// wrap is unsupported in golang b64, try to implement it later on
	o, err := b64.StdEncoding.DecodeString(in.String())
	if err != nil {
		panic(err)
	}
	return exec.AsValue(string(o))
}

// Get a custom environment with new filter we defined based on gonja.DefaultEnvironment
func CustomEnvironment() *exec.Environment {
	e := gonja.DefaultEnvironment
	if !e.Filters.Exists("b64encode") {
		e.Filters.Register("b64encode", filterFuncB64Encode)
	}
	if !e.Filters.Exists("b64decode") {
		e.Filters.Register("b64decode", filterFuncB64Decode)
	}
	return e
}

// This examples show how to template bytes slice with custom config and custom environment using our own new filters (b64encode and b64decode)
func TemplateFromBytesWithConfig(source []byte, config *config.Config) (*exec.Template, error) {
	rootID := fmt.Sprintf("root-%s", string(sha256.New().Sum(source)))

	loader, err := loaders.NewFileSystemLoader("")
	if err != nil {
		return nil, err
	}
	shiftedLoader, err := loaders.NewShiftedLoader(rootID, bytes.NewReader(source), loader)
	if err != nil {
		return nil, err
	}

	return exec.NewTemplate(rootID, config, shiftedLoader, CustomEnvironment())
}

func TemplateFromStringWithConfig(source string, config *config.Config) (*exec.Template, error) {
	return TemplateFromBytesWithConfig([]byte(source), config)
}

// Template a jinja2 string using our own custom environment with custom filter and jinja2 config.
func TemplateString(srcString string, data map[string]interface{}) string {
	tmpl, err := TemplateFromStringWithConfig(srcString, &CustomConfig)
	if err != nil {
		panic(err)
	}
	execContext := exec.NewContext(data)
	out, err := tmpl.ExecuteToString(execContext)
	if err != nil {
		panic(err)
	}
	return out
}

func main() {
	// Testing some stuff
	testEncodeString := `{% set output = mystring | b64encode %}{$ output $}`
	data := map[string]interface{}{
		"mystring": "abcd",
	}
	testEncodeDataOut := TemplateString(testEncodeString, data)
	fmt.Printf("b64encode output: '%s'\n", testEncodeDataOut)

	data = map[string]interface{}{
		"mystring": testEncodeDataOut,
	}
	testDecodeString := `{% set output = mystring | b64decode %}{$ output $}`
	testDecodeDataOut := TemplateString(testDecodeString, data)
	fmt.Printf("b64decode output: '%s'\n", testDecodeDataOut)
	if testDecodeDataOut == "abcd" {
		println("matched")
	}
}
