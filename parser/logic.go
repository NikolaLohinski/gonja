package parser

import (
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

var compareOps = []tokens.Type{
	tokens.Equals,
	tokens.Ne,
	tokens.GreaterThan,
	tokens.GreaterThanOrEqual,
	tokens.LowerThan,
	tokens.LowerThanOrEqual,
}

func BinOp(token *tokens.Token) *nodes.BinOperator {
	return &nodes.BinOperator{Token: token}
}

func (p *Parser) ParseLogicalExpression() (nodes.Expression, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (nodes.Expression, error) {
	var expr nodes.Expression

	expr, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.Current(tokens.Or) != nil {
		op := BinOp(p.Pop())
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		expr = &nodes.BinaryExpression{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}

	return expr, nil
}

func (p *Parser) parseAnd() (nodes.Expression, error) {
	var expr nodes.Expression

	expr, err := p.parseNot()
	if err != nil {
		return nil, err
	}

	for p.Current(tokens.And) != nil {
		op := BinOp(p.Pop())

		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}

		expr = &nodes.BinaryExpression{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}

	return expr, nil
}

func (p *Parser) parseNot() (nodes.Expression, error) {
	op := p.Match(tokens.Not)
	expr, err := p.parseCompare()
	if err != nil {
		return nil, err
	}

	if op != nil {
		expr = &nodes.Negation{
			Operator: op,
			Term:     expr,
		}
	}

	return expr, nil
}

func (p *Parser) parseCompare() (nodes.Expression, error) {
	var expr nodes.Expression

	expr, err := p.ParseMath()
	if err != nil {
		return nil, err
	}

	for p.Current(compareOps...) != nil {
		op := p.Pop()

		right, err := p.ParseMath()
		if err != nil {
			return nil, err
		}

		if right != nil {
			expr = &nodes.BinaryExpression{
				Left:     expr,
				Operator: BinOp(op),
				Right:    right,
			}
		}
	}

	expr, err = p.ParseTest(expr)
	if err != nil {
		return nil, err
	}

	return expr, nil
}
