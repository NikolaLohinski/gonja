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

// CallControlStructure implements Jinja2's
//
//	{% call macro_name(args...) %} body {% endcall %}
//
// The body is exposed inside the macro as the `caller()` function.
type CallControlStructure struct {
	Location *tokens.Token
	Call     *nodes.Call    // the macro invocation
	Body     *nodes.Wrapper // body to expose as caller()
}

func (c *CallControlStructure) Position() *tokens.Token { return c.Location }
func (c *CallControlStructure) String() string {
	t := c.Position()
	return fmt.Sprintf("CallControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (c *CallControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	// Build the caller() closure: renders the {% call %}...{% endcall %} body.
	callerFn := func(va *exec.VarArgs) *exec.Value {
		buf := new(bytes.Buffer)
		sub := r.Inherit()
		originalOutput := sub.Output
		sub.Output = buf
		err := sub.ExecuteWrapper(c.Body)
		sub.Output = originalOutput
		if err != nil {
			return exec.AsValue(err)
		}
		return exec.AsSafeValue(buf.String())
	}

	// Stash the caller in the context so MacroNodeToFunc can find and bind it.
	r.Environment.Context.Set("__gonja_caller__", callerFn)

	// Invoke the macro by evaluating the call expression. The result is the
	// macro's rendered output, which we write to the renderer.
	result := r.Eval(c.Call)
	if result.IsError() {
		return errors.Wrapf(result, "Unable to execute {%% call %%}")
	}

	// Ensure stash is cleared even if the macro didn't consume it.
	r.Environment.Context.Set("__gonja_caller__", nil)

	if _, err := r.Output.Write([]byte(result.String())); err != nil {
		return err
	}
	return nil
}

func callParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	cs := &CallControlStructure{Location: p.Current()}

	// Parse the macro-call expression: name(arg1, arg2, key=val, ...).
	expr, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	call, ok := expr.(*nodes.Call)
	if !ok {
		return nil, args.Error("Expected a macro call expression after {% call %}", args.Current())
	}
	cs.Call = call

	if !args.End() {
		return nil, args.Error("Malformed call-tag arguments.", args.Current())
	}

	// Parse body until {% endcall %}.
	wrapper, endargs, err := p.WrapUntil("endcall")
	if err != nil {
		return nil, err
	}
	cs.Body = wrapper
	if !endargs.End() {
		return nil, endargs.Error("endcall takes no arguments", nil)
	}

	_ = tokens.RightParenthesis // keep tokens import used if call args change
	return cs, nil
}
