package exec_test

import (
	"errors"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/tokens"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// constant evaluates to a fixed value.
type constant struct{ value any }

func (constant) Position() *tokens.Token { return &tokens.Token{} }
func (constant) String() string          { return "constant" }
func (c constant) Evaluate(*exec.Evaluator) *exec.Value {
	return exec.AsValue(c.value)
}

// twice evaluates its operand through the evaluator and doubles it.
type twice struct{ operand nodes.Expression }

func (twice) Position() *tokens.Token { return &tokens.Token{} }
func (twice) String() string          { return "twice" }
func (t twice) Evaluate(evaluator *exec.Evaluator) *exec.Value {
	value := evaluator.Eval(t.operand)
	if value.IsError() {
		return value
	}
	return exec.AsValue(value.Integer() * 2)
}

// opaque is an expression node that cannot evaluate itself.
type opaque struct{}

func (opaque) Position() *tokens.Token { return &tokens.Token{} }
func (opaque) String() string          { return "opaque" }

var _ = Context("evaluable expressions", func() {
	var evaluator = new(*exec.Evaluator)
	BeforeEach(func() {
		*evaluator = &exec.Evaluator{Config: gonja.DefaultConfig, Environment: gonja.DefaultEnvironment}
	})
	It("should evaluate a custom expression node", func() {
		Expect((*evaluator).Eval(constant{value: 42}).Integer()).To(Equal(42))
	})
	It("should let a custom node evaluate its operands", func() {
		Expect((*evaluator).Eval(twice{operand: &nodes.Integer{Val: 21}}).Integer()).To(Equal(42))
	})
	It("should evaluate a custom node inside a parsed template", func() {
		template, err := gonja.FromString("{{ placeholder | upper }}")
		Expect(err).ToNot(HaveOccurred())
		output, ok := template.Root().Nodes[0].(*nodes.Output)
		Expect(ok).To(BeTrue())
		filtered, ok := output.Expression.(*nodes.FilteredExpression)
		Expect(ok).To(BeTrue())
		filtered.Expression = constant{value: "custom"}
		Expect(template.ExecuteToString(exec.EmptyContext())).To(Equal("CUSTOM"))
	})
	It("should propagate errors from a custom node", func() {
		cause := errors.New("custom failure")
		value := (*evaluator).Eval(constant{value: cause})
		Expect(value.IsError()).To(BeTrue())
		Expect(value.Error()).To(Equal("custom failure"))
	})
	It("should still reject expression types that cannot evaluate themselves", func() {
		value := (*evaluator).Eval(opaque{})
		Expect(value.IsError()).To(BeTrue())
		Expect(value.Error()).To(ContainSubstring(`Unknown expression type "exec_test.opaque"`))
	})
})
