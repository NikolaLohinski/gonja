package statements

import (
	"fmt"

	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	"github.com/nikolalohinski/gonja/tokens"
	"github.com/pkg/errors"
)

type MacroStmt struct {
	*nodes.Macro
}

func (stmt *MacroStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("MacroStmt(Macro=%s Line=%d Col=%d)", stmt.Macro, t.Line, t.Col)
}

func (stmt *MacroStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	macro, err := exec.MacroNodeToFunc(stmt.Macro, r)
	if err != nil {
		return errors.Wrapf(err, `Unable to parse marco '%s'`, stmt.Name)
	}
	r.Ctx.Set(stmt.Name, macro)
	return nil
}

func macroParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &nodes.Macro{
		Location: p.Current(),
		Kwargs:   []*nodes.Pair{},
	}

	name := args.Match(tokens.Name)
	if name == nil {
		return nil, args.Error("Macro-tag needs at least an identifier as name.", nil)
	}
	stmt.Name = name.Val

	if args.Match(tokens.Lparen) == nil {
		return nil, args.Error("Expected '('.", nil)
	}

	for args.Match(tokens.Rparen) == nil {
		argName := args.Match(tokens.Name)
		if argName == nil {
			return nil, args.Error("Expected argument name as identifier.", nil)
		}

		if args.Match(tokens.Assign) != nil {
			expr, err := args.ParseExpression()
			if err != nil {
				return nil, err
			}
			stmt.Kwargs = append(stmt.Kwargs, &nodes.Pair{
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
			stmt.Kwargs = append(stmt.Kwargs, arg)
		}

		if args.Match(tokens.Rparen) != nil {
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
	stmt.Wrapper = wrapper

	if !endargs.End() {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	p.Template.Macros[stmt.Name] = stmt

	return &MacroStmt{stmt}, nil
}

func init() {
	All.Register("macro", macroParser)
}
