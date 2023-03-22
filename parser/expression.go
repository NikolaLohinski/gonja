package parser

import (
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/tokens"
	log "github.com/sirupsen/logrus"
)

// ParseFilterExpression parses an optionnal filter chain for a node
func (p *Parser) ParseFilterExpression(expr nodes.Expression) (nodes.Expression, error) {
	log.WithFields(log.Fields{
		"current": p.Current(),
	}).Trace("ParseFilterExpression")

	if p.Current(tokens.Pipe) != nil {

		filtered := &nodes.FilteredExpression{
			Expression: expr,
		}
		for p.Match(tokens.Pipe) != nil {
			// Parse one single filter
			filter, err := p.ParseFilter()
			if err != nil {
				return nil, err
			}

			// Check sandbox filter restriction
			// if _, isBanned := p.template.set.bannedFilters[filter.name]; isBanned {
			// 	return nil, p.Error(fmt.Sprintf("Usage of filter '%s' is not allowed (sandbox restriction active).", filter.name), nil)
			// }

			filtered.Filters = append(filtered.Filters, filter)
		}
		expr = filtered
	}

	log.WithFields(log.Fields{
		"expr": expr,
	}).Trace("ParseFilterExpression return")
	return expr, nil
}

// ParseExpression parses an expression with optionnal filters
// Nested expression shoulds call this method
func (p *Parser) ParseExpression() (nodes.Expression, error) {
	log.WithFields(log.Fields{
		"current": p.Current(),
	}).Trace("ParseExpression")
	var expr nodes.Expression

	expr, err := p.ParseLogicalExpression()
	if err != nil {
		return nil, err
	}

	expr, err = p.ParseFilterExpression(expr)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"expr": expr,
	}).Trace("ParseExpression return")
	return expr, nil
}

func (p *Parser) ParseExpressionNode() (nodes.Node, error) {
	log.WithFields(log.Fields{
		"current": p.Current(),
	}).Trace("ParseExpressionNode")

	tok := p.Match(tokens.VariableBegin)
	if tok == nil {
		return nil, p.Error("'{{' expected here", p.Current())
	}

	node := &nodes.Output{
		Start: tok,
		Trim: &nodes.Trim{
			Left: tok.Val[len(tok.Val)-1] == '-',
		},
	}

	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.Error("Expected an expression.", p.Current())
	}
	node.Expression = expr

	tok = p.Match(tokens.VariableEnd)
	if tok == nil {
		return nil, p.Error("'}}' expected here", p.Current())
	}
	node.End = tok
	node.Trim.Right = tok.Val[0] == '-'

	log.WithFields(log.Fields{
		"node": node,
	}).Trace("parseExpressionNode return")
	return node, nil
}
