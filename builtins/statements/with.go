package statements

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	"github.com/nikolalohinski/gonja/tokens"
)

type WithStmt struct {
	location *tokens.Token
	pairs    map[string]nodes.Expression
	wrapper  *nodes.Wrapper
}

func (stmt *WithStmt) Position() *tokens.Token { return stmt.location }
func (stmt *WithStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("WithStmt(Line=%d Col=%d)", t.Line, t.Col)
}

func (stmt *WithStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	sub := r.Inherit()

	for key, value := range stmt.pairs {
		val := r.Eval(value)
		if val.IsError() {
			return errors.Wrapf(val, `unable to evaluate parameter %s`, value)
		}
		sub.Environment.Context.Set(key, val)
	}

	return sub.ExecuteWrapper(stmt.wrapper)
}

func withParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &WithStmt{
		location: p.Current(),
		pairs:    map[string]nodes.Expression{},
	}

	wrapper, endargs, err := p.WrapUntil("endwith")
	if err != nil {
		return nil, err
	}
	stmt.wrapper = wrapper

	if !endargs.End() {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	for !args.End() {
		key := args.Match(tokens.Name)
		if key == nil {
			return nil, args.Error("Expected an identifier", args.Current())
		}
		if args.Match(tokens.Assign) == nil {
			return nil, args.Error("Expected '='.", args.Current())
		}
		value, err := args.ParseExpression()
		if err != nil {
			return nil, err
		}
		stmt.pairs[key.Val] = value

		if args.Match(tokens.Comma) == nil {
			break
		}
	}

	if !args.End() {
		return nil, errors.New("")
	}

	return stmt, nil
}

func init() {
	All.Register("with", withParser)
}
