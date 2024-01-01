package controlStructures

import "github.com/nikolalohinski/gonja/v2/exec"

var All = exec.ControlStructureSet{
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
}
