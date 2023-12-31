package gonja

import (
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
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
		Context:           builtins.Globals,
		Filters:           builtins.Filters,
		Tests:             builtins.Tests,
		ControlStructures: builtins.ControlStructures,
	}
}
