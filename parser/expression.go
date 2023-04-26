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

			filter, err := p.ParseFilter()
			if err != nil {
				return nil, err
			}

			filtered.Filters = append(filtered.Filters, filter)
		}
		expr = filtered
	}

	log.WithFields(log.Fields{
		"expr": expr,
	}).Trace("ParseFilterExpression return")
	return expr, nil
}

// ParseExpression parses an expression with optional filters
// Nested expression should call this method
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
	}

	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.Error("Expected an expression.", p.Current())
	}
	node.Expression = expr

	if p.MatchName("if") != nil {
		condition, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if condition == nil {
			return nil, p.Error("Expected a condition", p.Current())
		}
		node.Condition = condition

		if p.MatchName("else") != nil {
			alternative, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.Error("Expected an alternative", p.Current())
			}
			node.Alternative = alternative
		}
	}

	tok = p.Match(tokens.VariableEnd)
	if tok == nil {
		return nil, p.Error("'}}' expected here", p.Current())
	}
	node.End = tok
	if data := p.Current(tokens.Data); data != nil {
		data.Trim = data.Trim || len(node.End.Val) > 0 && node.End.Val[0] == '-'
	}

	log.WithFields(log.Fields{
		"node": node,
	}).Trace("parseExpressionNode return")
	return node, nil
}
