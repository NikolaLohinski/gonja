package gonja

import (
	"crypto/sha256"
	"fmt"
	"path"

	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

var (
	DefaultLoader      = loaders.MustNewFileSystemLoader("")
	DefaultConfig      = config.New()
	DefaultContext     = exec.EmptyContext().Update(builtins.GlobalFunctions).Update(builtins.GlobalVariables)
	DefaultEnvironment = &exec.Environment{
		Context:           DefaultContext,
		Filters:           builtins.Filters,
		Tests:             builtins.Tests,
		ControlStructures: builtins.ControlStructures,
	}
)

func FromString(source string) (*exec.Template, error) {
	return FromBytes([]byte(source))
}

func FromBytes(source []byte) (*exec.Template, error) {
	path := fmt.Sprintf("/%s", string(sha256.New().Sum(source)))

	loader, err := loaders.NewMemoryLoader(map[string]string{path: string(source)})
	if err != nil {
		return nil, err
	}

	return exec.NewTemplate(path, DefaultConfig, loader, DefaultEnvironment)
}

func FromFile(filepath string) (*exec.Template, error) {
	loader, err := loaders.NewFileSystemLoader(path.Dir(filepath))
	if err != nil {
		return nil, err
	}

	return exec.NewTemplate(path.Base(filepath), DefaultConfig, loader, DefaultEnvironment)
}
