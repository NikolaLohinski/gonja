package controlStructures

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type AutoescapeControlStructure struct {
	Wrapper    *nodes.Wrapper
	Autoescape bool
}

func (controlStructure *AutoescapeControlStructure) Position() *tokens.Token {
	return controlStructure.Wrapper.Position()
}
func (controlStructure *AutoescapeControlStructure) String() string {
	t := controlStructure.Position()
	return fmt.Sprintf("AutoescapeControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (controlStructure *AutoescapeControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	sub := r.Inherit()
	sub.Config.AutoEscape = controlStructure.Autoescape

	err := sub.ExecuteWrapper(controlStructure.Wrapper)
	if err != nil {
		return err
	}

	return nil
}

func autoescapeParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	controlStructure := &AutoescapeControlStructure{}

	wrapper, _, err := p.WrapUntil("endautoescape")
	if err != nil {
		return nil, err
	}
	controlStructure.Wrapper = wrapper

	modeToken := args.Match(tokens.Name)
	if modeToken == nil {
		return nil, args.Error("A mode is required for autoescape controlStructure.", nil)
	}
	if modeToken.Val == "true" {
		controlStructure.Autoescape = true
	} else if modeToken.Val == "false" {
		controlStructure.Autoescape = false
	} else {
		return nil, args.Error("Only 'true' or 'false' is valid as an autoescape controlStructure.", nil)
	}

	if !args.Stream().End() {
		return nil, args.Error("Malformed autoescape controlStructure args.", nil)
	}

	return controlStructure, nil
}
