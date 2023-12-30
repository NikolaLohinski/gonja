package statements

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
	"github.com/pkg/errors"
)

type SetStmt struct {
	location   *tokens.Token
	target     nodes.Expression
	expression nodes.Expression
}

func (stmt *SetStmt) Position() *tokens.Token { return stmt.location }
func (stmt *SetStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("SetStmt(Line=%d Col=%d)", t.Line, t.Col)
}

func (stmt *SetStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	// Evaluate expression
	value := r.Eval(stmt.expression)
	if value.IsError() {
		return value
	}

	switch n := stmt.target.(type) {
	case *nodes.Name:
		r.Environment.Context.Set(n.Name.Val, value.Interface())
	case *nodes.GetAttribute:
		target := r.Eval(n.Node)
		if target.IsError() {
			return errors.Wrapf(target, `Unable to evaluate target %s`, n)
		}
		if err := target.Set(exec.AsValue(n.Attr), value.Interface()); err != nil {
			return errors.Wrapf(err, `Unable to set value on "%s"`, n.Attr)
		}
	case *nodes.GetItem:
		target := r.Eval(n.Node)
		if target.IsError() {
			return errors.Wrapf(target, `Unable to evaluate target %s`, n)
		}
		arg := r.Eval(n.Arg)
		if arg.IsError() {
			return errors.Wrapf(target, `Unable to evaluate argument %s`, n.Arg)
		}
		if err := target.Set(arg, value.Interface()); err != nil {
			return errors.Wrapf(err, `Unable to set value on "%s"`, n.Arg)
		}
	default:
		return errors.Errorf(`Illegal set target node %s`, n)
	}

	return nil
}

func setParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &SetStmt{
		location: p.Current(),
	}

	// Parse variable name
	ident, err := args.ParseVariable()
	if err != nil {
		return nil, errors.Wrap(err, `Unable to parse identifier`)
	}
	switch n := ident.(type) {
	case *nodes.Name, *nodes.Call, *nodes.GetItem, *nodes.GetAttribute:
		stmt.target = n
	default:
		return nil, errors.Errorf(`Unexpected set target %s`, n)
	}

	if args.Match(tokens.Assign) == nil {
		return nil, args.Error("Expected '='.", args.Current())
	}

	// Variable expression
	expr, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	stmt.expression = expr

	// Remaining arguments
	if !args.End() {
		return nil, args.Error("Malformed 'set'-tag args.", args.Current())
	}

	return stmt, nil
}

func init() {
	All.Register("set", setParser)
}
