// Package controlstructures provides built-in control structure implementations.
package controlstructures

import (
	"github.com/ardanlabs/gonja/exec"
	"github.com/ardanlabs/gonja/parser"
)

var All = exec.NewControlStructureSet(map[string]parser.ControlStructureParser{
	"autoescape": autoescapeParser,
	"block":      blockParser,
	"extends":    extendsParser,
	"filter":     filterParser,
	"for":        forParser,
	"from":       fromParser,
	"if":         ifParser,
	"import":     importParser,
	"include":    includeParser,
	"macro":      macroParser,
	"raw":        rawParser,
	"set":        setParser,
	"with":       withParser,
})
