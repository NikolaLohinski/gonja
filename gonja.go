package gonja

import (
	"github.com/nikolalohinski/gonja/builtins"
	"github.com/nikolalohinski/gonja/config"
	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/loaders"
)

var (
	version = "0.0.0+trunk"

	DefaultLoader      loaders.Loader
	DefaultConfig      *config.Config
	DefaultEnvironment *exec.Environment
)

func init() {
	builtins.Globals.Set("gonja", map[string]interface{}{
		"version": version,
	})

	DefaultLoader = loaders.MustNewFileSystemLoader("")
	DefaultConfig = config.New()
	DefaultEnvironment = &exec.Environment{
		Context:    builtins.Globals,
		Filters:    builtins.Filters,
		Tests:      builtins.Tests,
		Statements: builtins.Statements,
	}
}
