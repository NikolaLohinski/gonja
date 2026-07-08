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

func (mcs *MacroControlStructure) String() string {
	t := mcs.Position()
	return fmt.Sprintf("MacroControlStructure(Macro=%s Line=%d Col=%d)", mcs.Macro, t.Line, t.Col)
}

func (mcs *MacroControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	macro, err := exec.MacroNodeToFunc(mcs.Macro, r)
	if err != nil {
		return errors.Wrapf(err, `Unable to parse macro '%s'`, mcs.Name)
	}
	r.Environment.Context.Set(mcs.Name, macro)
	return nil
}

func macroParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	macro := &nodes.Macro{
		Location: p.Current(),
		Kwargs:   []*nodes.Pair{},
	}

	name := args.Match(tokens.Name)
	if name == nil {
		return nil, args.Error("Macro-tag needs at least an identifier as name.", nil)
	}
	macro.Name = name.Val

	if args.Match(tokens.LeftParenthesis) == nil {
		return nil, args.Error("Expected '('.", nil)
	}

	for args.Match(tokens.RightParenthesis) == nil {
		// Detect Jinja2 *args / **kwargs syntax. The lexer emits these as
		// Multiply ("*") and Power ("**") tokens respectively.
		if args.Match(tokens.Power) != nil {
			argName := args.Match(tokens.Name)
			if argName == nil {
				return nil, args.Error("Expected identifier after '**'.", nil)
			}
			if macro.KwArgsName != "" {
				return nil, args.Error("Only one '**kwargs' parameter is allowed.", nil)
			}
			macro.KwArgsName = argName.Val
		} else if args.Match(tokens.Multiply) != nil {
			argName := args.Match(tokens.Name)
			if argName == nil {
				return nil, args.Error("Expected identifier after '*'.", nil)
			}
			if macro.VarArgsName != "" {
				return nil, args.Error("Only one '*args' parameter is allowed.", nil)
			}
			macro.VarArgsName = argName.Val
		} else {
			argName := args.Match(tokens.Name)
			if argName == nil {
				return nil, args.Error("Expected argument name as identifier.", nil)
			}

			if args.Match(tokens.Assign) != nil {
				expr, err := args.ParseExpression()
				if err != nil {
					return nil, err
				}
				macro.Kwargs = append(macro.Kwargs, &nodes.Pair{
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
				macro.Kwargs = append(macro.Kwargs, arg)
			}
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
	macro.Wrapper = wrapper

	if !endargs.End() {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	p.Template.Macros[macro.Name] = macro

	return &MacroControlStructure{macro}, nil
}
