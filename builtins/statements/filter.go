package statements

import (
	// "bytes"

	// "github.com/nikolalohinski/gonja/exec"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	"github.com/nikolalohinski/gonja/tokens"
)

type FilterStmt struct {
	position    *tokens.Token
	bodyWrapper *nodes.Wrapper
	filterChain []*nodes.FilterCall
}

func (stmt *FilterStmt) Position() *tokens.Token { return stmt.position }
func (stmt *FilterStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("FilterStmt(Line=%d Col=%d)", t.Line, t.Col)
}

func (node *FilterStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	var out strings.Builder
	sub := r.Inherit()
	sub.Output = &out
	// temp := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB size

	err := sub.ExecuteWrapper(node.bodyWrapper)
	if err != nil {
		return err
	}

	value := exec.AsValue(out.String())

	for _, call := range node.filterChain {
		value = r.Evaluator().ExecuteFilter(call, value)
		if value.IsError() {
			return errors.Wrapf(value, `Unable to apply filter %s (Line: %d Col: %d, near %s`,
				call.Name, call.Token.Line, call.Token.Col, call.Token.Val)
		}
	}

	_, err = r.Output.WriteString(value.String())

	return err
}

func filterParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &FilterStmt{
		position: p.Current(),
	}

	wrapper, _, err := p.WrapUntil("endfilter")
	if err != nil {
		return nil, err
	}
	stmt.bodyWrapper = wrapper

	for !args.End() {
		filterCall, err := args.ParseFilter()
		if err != nil {
			return nil, err
		}

		stmt.filterChain = append(stmt.filterChain, filterCall)

		if args.Match(tokens.Pipe) == nil {
			break
		}
	}

	if !args.End() {
		return nil, p.Error("Malformed filter-tag args.", nil)
	}

	return stmt, nil
}

func init() {
	All.Register("filter", filterParser)
}
