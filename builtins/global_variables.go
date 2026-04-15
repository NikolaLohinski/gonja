package builtins

import "github.com/ardanlabs/gonja/exec"

var GlobalVariables = exec.NewContext(map[string]any{
	"gonja": map[string]any{
		"version": "v0.0.0+trunk",
	},
})
