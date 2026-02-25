package parser

import (
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

func (p *Parser) ParseTest(expr nodes.Expression) (nodes.Expression, error) {
	expr, err := p.ParseFilterExpression(expr)
	if err != nil {
		return nil, err
	}
	if p.Current(
		tokens.GreaterThan,
		tokens.GreaterThanOrEqual,
		tokens.LowerThan,
		tokens.LowerThanOrEqual,
		tokens.Not,
		tokens.In,
		tokens.Is,
	) != nil {
		_ = p.Match(tokens.Is) // ignore the is keyword entirely if present

		not := p.Match(tokens.Not)
		ident := p.Next()

		test := &nodes.TestCall{
			Token:  ident,
			Name:   ident.Val,
			Args:   []nodes.Expression{},
			Kwargs: map[string]nodes.Expression{},
		}
		// avoid trying to parse "else" as test arguments
		if p.CurrentName("else") == nil {
			arg, err := p.ParseVariableOrLiteral()
			if err == nil && arg != nil {
				test.Args = append(test.Args, arg)
			}
		}
		expr = &nodes.TestExpression{
			Expression: expr,
			Test:       test,
		}

		if not != nil {
			expr = &nodes.Negation{
				Term:     expr,
				Operator: not,
			}
		}
	}

	return expr, nil
}
