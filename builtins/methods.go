package builtins

import (
	methods "github.com/ardanlabs/gonja/builtins/methods"
	"github.com/ardanlabs/gonja/exec"
)

// Methods exports all builtins methods.
var Methods exec.Methods = methods.All
