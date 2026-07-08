package controlStructures

import (
	"bytes"
	"fmt"
	"io"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

// TransControlStructure implements Jinja2's {% trans %} ... {% endtrans %}
// with optional {% pluralize %} for plural forms.
//
// Supported forms:
//
//	{% trans %}Hello{% endtrans %}
//	{% trans name=user.name %}Hello {{ name }}{% endtrans %}
//	{% trans count=items|length %}{{ count }} item{% pluralize %}{{ count }} items{% endtrans %}
type TransControlStructure struct {
	Location     *tokens.Token
	Variables    map[string]nodes.Expression
	CountVar     string
	SingularBody *nodes.Wrapper
	PluralBody   *nodes.Wrapper // may be nil
}

func (c *TransControlStructure) Position() *tokens.Token { return c.Location }
func (c *TransControlStructure) String() string {
	t := c.Position()
	return fmt.Sprintf("TransControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (stmt *TransControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	sub := r.Inherit()
	ctx := sub.Environment.Context

	// Evaluate variables and inject into context.
	var countVal *exec.Value
	for name, expr := range stmt.Variables {
		v := r.Eval(expr)
		if v.IsError() {
			return v
		}
		ctx.Set(name, v)
		if name == stmt.CountVar {
			countVal = v
		}
	}

	// Decide which body to render.
	useBody := stmt.SingularBody
	if stmt.PluralBody != nil && countVal != nil && countVal.IsInteger() && countVal.Integer() != 1 {
		useBody = stmt.PluralBody
	}

	// Render the body into a buffer by swapping the renderer's Output writer.
	buf := new(bytes.Buffer)
	originalOutput := sub.Output
	sub.Output = buf
	err := sub.ExecuteWrapper(useBody)
	sub.Output = originalOutput
	if err != nil {
		return err
	}
	rendered := buf.String()

	// Optional translation hook lookup.
	if tf, ok := r.Environment.Context.Get("__gonja_trans_hook__"); ok {
		if hook, ok := tf.(func(string) string); ok {
			rendered = hook(rendered)
		}
	}

	// Write the (possibly translated) result to the outer writer.
	if _, err := io.WriteString(r.Output, rendered); err != nil {
		return err
	}
	return nil
}

// transParser parses {% trans var1=expr var2=expr %} ... {% endtrans %},
// optionally with {% pluralize %} dividing singular/plural bodies.
func transParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	ctrl := &TransControlStructure{
		Location:  p.Current(),
		Variables: map[string]nodes.Expression{},
	}

	// Parse key=expr pairs as trans arguments.
	for !args.End() {
		name := args.Match(tokens.Name)
		if name == nil {
			return nil, args.Error("expected an identifier in trans arguments", args.Current())
		}
		if args.Match(tokens.Assign) == nil {
			return nil, args.Error("expected '=' after trans argument", args.Current())
		}
		expr, err := args.ParseExpression()
		if err != nil {
			return nil, err
		}
		ctrl.Variables[name.Val] = expr
		if ctrl.CountVar == "" && (name.Val == "count" || name.Val == "n") {
			ctrl.CountVar = name.Val
		}
	}

	// Singular body up to {% pluralize %} or {% endtrans %}.
	wrapper, endargs, err := p.WrapUntil("pluralize", "endtrans")
	if err != nil {
		return nil, err
	}
	ctrl.SingularBody = wrapper
	if !endargs.End() {
		return nil, endargs.Error("trans block does not accept end-tag arguments", nil)
	}

	if wrapper.EndTag == "pluralize" {
		plural, endargs2, err := p.WrapUntil("endtrans")
		if err != nil {
			return nil, err
		}
		ctrl.PluralBody = plural
		if !endargs2.End() {
			return nil, endargs2.Error("endtrans does not accept arguments", nil)
		}
	}
	return ctrl, nil
}
