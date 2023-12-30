package statements

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type IfStmt struct {
	location   *tokens.Token
	conditions []nodes.Expression
	wrappers   []*nodes.Wrapper
}

func (stmt *IfStmt) Position() *tokens.Token { return stmt.location }
func (stmt *IfStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("IfStmt(Line=%d Col=%d)", t.Line, t.Col)
}

func (node *IfStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	for i, condition := range node.conditions {
		result := r.Eval(condition)
		if result.IsError() {
			return result
		}

		if result.IsTrue() {
			return r.ExecuteWrapper(node.wrappers[i])
		}
		// Last condition?
		if len(node.conditions) == i+1 && len(node.wrappers) > i+1 {
			return r.ExecuteWrapper(node.wrappers[i+1])
		}
	}
	return nil
}

func ifParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	log.WithFields(log.Fields{
		"arg":     args.Current(),
		"current": p.Current(),
	}).Trace("ParseIf")
	ifNode := &IfStmt{
		location: args.Current(),
	}

	// Parse first and main IF condition
	condition, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	ifNode.conditions = append(ifNode.conditions, condition)

	if !args.End() {
		return nil, args.Error("If-condition is malformed.", nil)
	}

	// Check the rest
	for {
		wrapper, tagArgs, err := p.WrapUntil("elif", "else", "endif")
		if err != nil {
			return nil, err
		}
		ifNode.wrappers = append(ifNode.wrappers, wrapper)

		if wrapper.EndTag == "elif" {
			// elif can take a condition
			condition, err = tagArgs.ParseExpression()
			if err != nil {
				return nil, err
			}
			ifNode.conditions = append(ifNode.conditions, condition)

			if !tagArgs.End() {
				return nil, tagArgs.Error("Elif-condition is malformed.", nil)
			}
		} else {
			if !tagArgs.End() {
				// else/endif can't take any conditions
				return nil, tagArgs.Error("Arguments not allowed here.", nil)
			}
		}

		if wrapper.EndTag == "endif" {
			break
		}
	}

	return ifNode, nil
}

func init() {
	All.Register("if", ifParser)
}
