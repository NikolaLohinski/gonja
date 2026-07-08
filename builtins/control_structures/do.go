package controlStructures

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type DoControlStructure struct {
	Location   *tokens.Token
	Expression nodes.Expression
}

func (c *DoControlStructure) Position() *tokens.Token { return c.Location }
func (c *DoControlStructure) String() string {
	t := c.Position()
	return fmt.Sprintf("DoControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (stmt *DoControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	// In gonja v2, r.Eval returns *exec.Value (which implements error via IsError)
	val := r.Eval(stmt.Expression)
	if val.IsError() {
		return val
	}
	return nil
}

func doParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	ctrl := &DoControlStructure{
		Location: p.Current(),
	}
	expr, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	ctrl.Expression = expr
	if !args.End() {
		return nil, args.Error("Malformed do-tag args.", args.Current())
	}
	return ctrl, nil
}
