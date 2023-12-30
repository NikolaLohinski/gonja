package statements

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type RawStmt struct {
	data *nodes.Data
}

func (stmt *RawStmt) Position() *tokens.Token { return stmt.data.Position() }
func (stmt *RawStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("RawStmt(Line=%d Col=%d)", t.Line, t.Col)
}

func (stmt *RawStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	_, err := r.Output.WriteString(stmt.data.Data.Val)
	return err
}

func rawParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &RawStmt{}

	wrapper, _, err := p.WrapUntil("endraw")
	if err != nil {
		return nil, err
	}
	node := wrapper.Nodes[0]
	data, ok := node.(*nodes.Data)
	if ok {
		stmt.data = data
	} else {
		return nil, p.Error("raw statement can only contains a single data node", node.Position())
	}

	if !args.End() {
		return nil, args.Error("raw statement doesn't accept parameters.", args.Current())
	}

	return stmt, nil
}

func init() {
	All.Register("raw", rawParser)
}
