package controlStructures

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
	"github.com/pkg/errors"
)

type MacroControlStructure struct {
	*nodes.Macro
}

func (controlStructure *MacroControlStructure) String() string {
	t := controlStructure.Position()
	return fmt.Sprintf("MacroControlStructure(Macro=%s Line=%d Col=%d)", controlStructure.Macro, t.Line, t.Col)
}

func (controlStructure *MacroControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	macro, err := exec.MacroNodeToFunc(controlStructure.Macro, r)
	if err != nil {
		return errors.Wrapf(err, `Unable to parse marco '%s'`, controlStructure.Name)
	}
	r.Environment.Context.Set(controlStructure.Name, macro)
	return nil
}

func macroParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	controlStructure := &nodes.Macro{
		Location: p.Current(),
		Kwargs:   []*nodes.Pair{},
	}

	name := args.Match(tokens.Name)
	if name == nil {
		return nil, args.Error("Macro-tag needs at least an identifier as name.", nil)
	}
	controlStructure.Name = name.Val

	if args.Match(tokens.LeftParenthesis) == nil {
		return nil, args.Error("Expected '('.", nil)
	}

	for args.Match(tokens.RightParenthesis) == nil {
		argName := args.Match(tokens.Name)
		if argName == nil {
			return nil, args.Error("Expected argument name as identifier.", nil)
		}

		if args.Match(tokens.Assign) != nil {
			expr, err := args.ParseExpression()
			if err != nil {
				return nil, err
			}
			controlStructure.Kwargs = append(controlStructure.Kwargs, &nodes.Pair{
				Key: &nodes.String{
					Location: argName,
					Val:      argName.Val,
				},
				Value: expr,
			})
		} else {
			arg := &nodes.Pair{
				Key: &nodes.String{
					Location: argName,
					Val:      argName.Val,
				},
			}
			if p.Config.StrictUndefined {
				arg.Value = &nodes.Error{
					Location: argName,
					Error:    fmt.Errorf("parameter \"%s\" was not provided", argName.Val),
				}
			} else {
				arg.Value = &nodes.None{
					Location: argName,
				}
			}
			controlStructure.Kwargs = append(controlStructure.Kwargs, arg)
		}

		if args.Match(tokens.RightParenthesis) != nil {
			break
		}
		if args.Match(tokens.Comma) == nil {
			return nil, args.Error("Expected ',' or ')'.", nil)
		}
	}

	if !args.End() {
		return nil, args.Error("Malformed macro-tag.", nil)
	}

	wrapper, endargs, err := p.WrapUntil("endmacro")
	if err != nil {
		return nil, err
	}
	controlStructure.Wrapper = wrapper

	if !endargs.End() {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	p.Template.Macros[controlStructure.Name] = controlStructure

	return &MacroControlStructure{controlStructure}, nil
}
