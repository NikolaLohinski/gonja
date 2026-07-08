package controlStructures

import (
	"bytes"
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
	"github.com/pkg/errors"
)

type SetControlStructure struct {
	location    *tokens.Token
	target      nodes.Expression
	expression  nodes.Expression
	condition   nodes.Expression
	alternative nodes.Expression
	body        *nodes.Wrapper
}

func (scs *SetControlStructure) Position() *tokens.Token {
	return scs.location
}
func (scs *SetControlStructure) String() string {
	t := scs.Position()
	return fmt.Sprintf("SetControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (scs *SetControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	var value *exec.Value

	// --- Block form: {% set x %}...{% endset %} ---
	if scs.body != nil {
		buf := new(bytes.Buffer)
		sub := r.Inherit()
		originalOutput := sub.Output
		sub.Output = buf
		err := sub.ExecuteWrapper(scs.body)
		sub.Output = originalOutput
		if err != nil {
			return err
		}
		// Jinja2's block-set captures the rendered string as a Markup
		// (safe) value. Match that semantic.
		value = exec.AsSafeValue(buf.String())
	} else {
		// --- Expression form: {% set x = expr [if cond else alt] %} ---
		if scs.condition != nil && scs.alternative != nil {
			condition := r.Eval(scs.condition)
			if condition.IsError() {
				return condition
			}
			if !condition.IsNil() && condition.IsTrue() {
				value = r.Eval(scs.expression)
			} else {
				value = r.Eval(scs.alternative)
			}
		} else {
			value = r.Eval(scs.expression)
		}
	}

	if value == nil {
		return errors.Errorf(`Invalid value in 'set' tag: %s`, scs.expression)
	}
	if value.IsError() {
		return value
	}

	switch n := scs.target.(type) {
	case *nodes.Name:
		r.Environment.Context.Set(n.Name.Val, value.Interface())
	case *nodes.GetAttribute:
		target := r.Eval(n.Node)
		if target.IsError() {
			return errors.Wrapf(target, `Unable to evaluate target %s`, n)
		}
		if err := target.Set(exec.AsValue(n.Attribute), value.Interface()); err != nil {
			return errors.Wrapf(err, `Unable to set value on "%s"`, n.Attribute)
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

func setParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	cs := &SetControlStructure{
		location: p.Current(),
	}

	// Parse target (variable name / attr / item).
	ident, err := args.ParseVariableOrLiteral()
	if err != nil {
		return nil, errors.Wrap(err, `unable to parse identifier`)
	}
	switch n := ident.(type) {
	case *nodes.Name, *nodes.Call, *nodes.GetItem, *nodes.GetAttribute:
		cs.target = n
	default:
		return nil, errors.Errorf(`unexpected set target %s`, n)
	}

	if args.Match(tokens.Assign) == nil {
		if !args.End() {
			return nil, args.Error("Expected '=' or end of tag for block-set.", args.Current())
		}
		wrapper, endargs, err := p.WrapUntil("endset")
		if err != nil {
			return nil, err
		}
		if !endargs.End() {
			return nil, endargs.Error("endset takes no arguments", nil)
		}
		cs.body = wrapper
		return cs, nil
	}

	// --- Expression form: {% set x = expr [if cond else alt] %} ---
	expr, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	cs.expression = expr

	condition, alternative, err := args.ParseCondition()
	if err != nil {
		return nil, err
	}
	if condition != nil {
		cs.condition = condition
		if alternative == nil {
			return nil, args.Error("Malformed 'set' if else condition", args.Current())
		}
		cs.alternative = alternative
	}

	if !args.End() {
		return nil, args.Error("Malformed 'set' tag args.", args.Current())
	}

	return cs, nil
}
