package exec

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type namedString string

func TestUndefinedValueSemantics(t *testing.T) {
	lenient := UndefinedValue(false, "missing")
	strict := UndefinedValue(true, "required")

	if !lenient.IsUndefined() || lenient.IsStrictUndefined() || lenient.IsNil() {
		t.Fatalf("unexpected lenient undefined state: undefined=%v strict=%v nil=%v", lenient.IsUndefined(), lenient.IsStrictUndefined(), lenient.IsNil())
	}
	if !strict.IsUndefined() || !strict.IsStrictUndefined() {
		t.Fatalf("unexpected strict undefined state: undefined=%v strict=%v", strict.IsUndefined(), strict.IsStrictUndefined())
	}
	if AsValue(nil).IsUndefined() || AsValue("value").IsStrictUndefined() {
		t.Fatal("ordinary values were classified as undefined")
	}
	if lenient.String() != "" || lenient.IsTrue() || !lenient.Negate().IsTrue() {
		t.Fatal("lenient undefined does not have Jinja's empty/falsy behavior")
	}
	if got := strict.undefinedError().Error(); got != "required is undefined" {
		t.Fatalf("strict undefined error = %q", got)
	}
	if got := UndefinedValue(true, "").undefinedError().Error(); got != "value is undefined" {
		t.Fatalf("unnamed undefined error = %q", got)
	}

	// Undefined must survive the interface boundary used by contexts and calls.
	roundTrip := ToValue(lenient.Interface())
	if !roundTrip.IsUndefined() || roundTrip.IsNil() {
		t.Fatalf("undefined identity was lost: %#v", roundTrip.Interface())
	}
	if !lenient.EqualValueTo(strict) || lenient.EqualValueTo(AsValue(nil)) {
		t.Fatal("undefined equality should distinguish undefined from None")
	}
}

func TestUndefinedValueRepresentations(t *testing.T) {
	undefined := UndefinedValue(false, "missing")
	if got, want := (ValuesList{AsValue("x"), undefined}).String(), "['x', Undefined]"; got != want {
		t.Fatalf("ValuesList.String() = %q, want %q", got, want)
	}
	pair := &Pair{Key: AsValue("key"), Value: undefined}
	if got, want := pair.String(), "'key': Undefined"; got != want {
		t.Fatalf("Pair.String() = %q, want %q", got, want)
	}
	if got, want := AsValue([]any{undefined}).String(), "[Undefined]"; got != want {
		t.Fatalf("slice String() = %q, want %q", got, want)
	}
	if got, want := AsValue(map[string]any{"key": undefined}).String(), "{'key': Undefined}"; got != want {
		t.Fatalf("map String() = %q, want %q", got, want)
	}
}

func TestContainsCheckedUndefinedSemantics(t *testing.T) {
	strict := UndefinedValue(true, "needle")
	lenient := UndefinedValue(false, "needle")

	tests := []struct {
		name      string
		haystack  *Value
		needle    *Value
		want      bool
		wantError string
	}{
		{"strict haystack", UndefinedValue(true, "items"), AsValue(1), false, "items is undefined"},
		{"lenient haystack", UndefinedValue(false, "items"), AsValue(1), false, ""},
		{"None haystack", AsValue(nil), AsValue(1), false, "None is not iterable"},
		{"strict needle in empty sequence", AsValue([]int{}), strict, false, ""},
		{"strict needle in sequence", AsValue([]int{1}), strict, false, "needle is undefined"},
		{"strict value in sequence", AsValue([]any{UndefinedValue(true, "element")}), AsValue(1), false, "element is undefined"},
		{"lenient needle in string", AsValue("abc"), lenient, false, "string containment requires a string"},
		{"None needle in string", AsValue("abc"), AsValue(nil), false, "string containment requires a string"},
		{"substring", AsValue("abc"), AsValue("b"), true, ""},
		{"slice member", AsValue([]int{1, 2}), AsValue(2), true, ""},
		{"map key", AsValue(map[string]int{"a": 1}), AsValue("a"), true, ""},
		{"unsupported map key", AsValue(map[string]int{"a": 1}), AsValue(1.5), false, ""},
		{"unsupported haystack", AsValue(42), AsValue(4), false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.haystack.ContainsChecked(tt.needle)
			if got != tt.want {
				t.Fatalf("ContainsChecked() = %v, want %v", got, tt.want)
			}
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("ContainsChecked() error = %v", err)
				}
			} else if err == nil || err.Error() != tt.wantError {
				t.Fatalf("ContainsChecked() error = %v, want %q", err, tt.wantError)
			}
		})
	}

	// The compatibility helper intentionally discards the checked error.
	if UndefinedValue(true, "items").Contains(AsValue(1)) {
		t.Fatal("strict undefined unexpectedly contained a value")
	}
}

func TestOrderedDictValueAPIs(t *testing.T) {
	dict := &Dict{Pairs: []*Pair{
		{Key: AsValue("z"), Value: AsValue(1)},
		{Key: AsValue("a"), Value: AsValue(2)},
	}}
	value := AsValue(dict)
	items := value.Items()
	if len(items) != 2 || items[0].Key.String() != "z" || items[1].Key.String() != "a" {
		t.Fatalf("Items() did not preserve insertion order: %v", items)
	}
	// Items returns a separate slice, while retaining the dictionary's pairs.
	items[0] = nil
	if dict.Pairs[0] == nil || dict.Pairs[0].Key.String() != "z" {
		t.Fatal("Items() exposed the dictionary's pair slice")
	}
	if got := value.Keys().String(); got != "['z', 'a']" {
		t.Fatalf("Keys() = %q", got)
	}
	if got := dict.Get(AsValue(namedString("a"))); got.Integer() != 2 {
		t.Fatalf("named string key did not compare by text: %v", got)
	}
	if !AsValue(namedString("same")).EqualValueTo(AsValue("same")) {
		t.Fatal("named and built-in strings with the same text were not equal")
	}

	plainMapItems := AsValue(map[string]int{"only": 3}).Items()
	if len(plainMapItems) != 1 || plainMapItems[0].Key.String() != "only" || plainMapItems[0].Value.Integer() != 3 {
		t.Fatalf("map Items() = %v", plainMapItems)
	}
	if got := AsValue(3).Items(); len(got) != 0 {
		t.Fatalf("scalar Items() = %v", got)
	}
}

func TestEvaluatorSynthesizesUndefinedValues(t *testing.T) {
	tok := &tokens.Token{Val: "test"}
	for _, strict := range []bool{false, true} {
		evaluator := &Evaluator{Config: &config.Config{StrictUndefined: strict}}
		value := evaluator.Eval(&nodes.Undefined{Location: tok})
		if !value.IsUndefined() || value.IsStrictUndefined() != strict {
			t.Fatalf("strict=%v: Eval(Undefined) returned %#v", strict, value.Interface())
		}
	}

	evaluator := &Evaluator{Config: config.New()}
	conditional := &nodes.ConditionalExpression{
		Expression: &nodes.String{Location: tok, Val: "unused"},
		Condition:  &nodes.Bool{Location: tok, Val: false},
	}
	value := evaluator.Eval(conditional)
	if !value.IsUndefined() || value.IsStrictUndefined() {
		t.Fatalf("omitted conditional alternative returned %#v", value.Interface())
	}

	strictEvaluator := &Evaluator{Config: &config.Config{StrictUndefined: true}}
	negated := strictEvaluator.Eval(&nodes.Negation{
		Operator: tok,
		Term:     &nodes.Undefined{Location: tok},
	})
	if !negated.IsError() || negated.Error() != "value is undefined" {
		t.Fatalf("negating strict undefined returned %q", negated.Error())
	}
}

func TestEvaluatorUndefinedBinaryOperations(t *testing.T) {
	tok := func(kind tokens.Type, value string) *tokens.Token {
		return &tokens.Token{Type: kind, Val: value}
	}
	undefined := &nodes.Undefined{Location: tok(tokens.Name, "undefined")}
	integer := &nodes.Integer{Location: tok(tokens.Integer, "1"), Val: 1}
	empty := &nodes.List{Location: tok(tokens.LeftBracket, "["), Val: nil}
	nonempty := &nodes.List{Location: tok(tokens.LeftBracket, "["), Val: []nodes.Expression{integer}}
	binary := func(left nodes.Expression, operator tokens.Type, symbol string, right nodes.Expression) *nodes.BinaryExpression {
		return &nodes.BinaryExpression{
			Left:     left,
			Operator: &nodes.BinOperator{Token: tok(operator, symbol)},
			Right:    right,
		}
	}

	lenient := &Evaluator{Config: config.New()}
	if got := lenient.Eval(binary(undefined, tokens.Equals, "==", undefined)); !got.IsBool() || !got.Bool() {
		t.Fatalf("undefined == undefined returned %v", got)
	}
	if got := lenient.Eval(binary(undefined, tokens.Addition, "+", integer)); !got.IsError() {
		t.Fatalf("undefined addition returned %v", got)
	}
	if got := lenient.Eval(binary(&nodes.None{Location: tok(tokens.Name, "None")}, tokens.Addition, "+", integer)); !got.IsError() {
		t.Fatalf("None addition returned %v", got)
	}

	strictConfig := config.New()
	strictConfig.StrictUndefined = true
	strict := &Evaluator{Config: strictConfig}
	if got := strict.Eval(binary(undefined, tokens.Addition, "+", integer)); !got.IsError() || got.Error() != "value is undefined" {
		t.Fatalf("strict undefined addition returned %q", got.Error())
	}
	if got := strict.Eval(binary(undefined, tokens.In, "in", empty)); !got.IsBool() || got.Bool() {
		t.Fatalf("strict undefined in empty list returned %v", got)
	}
	if got := strict.Eval(binary(undefined, tokens.In, "in", nonempty)); !got.IsError() || got.Error() != "value is undefined" {
		t.Fatalf("strict undefined in non-empty list returned %q", got.Error())
	}

	if got := lenient.Eval(&nodes.UnaryExpression{
		Operator: tok(tokens.Addition, "+"),
		Term:     &nodes.String{Location: tok(tokens.String, "text"), Val: "text"},
	}); !got.IsError() {
		t.Fatalf("unary plus on text returned %v", got)
	}
	if got := lenient.Eval(&nodes.UnaryExpression{
		Operator: tok(tokens.Subtraction, "-"),
		Negative: true,
		Term:     undefined,
	}); !got.IsError() || got.Error() != "value is undefined" {
		t.Fatalf("unary minus on undefined returned %q", got.Error())
	}
}

func TestUndefinedIterationIsEmpty(t *testing.T) {
	for _, value := range []*Value{AsValue(nil), UndefinedValue(false, "items")} {
		called := false
		empty := false
		value.Iterate(func(_, _ int, _, _ *Value) bool {
			called = true
			return true
		}, func() { empty = true })
		if called || !empty {
			t.Fatalf("iteration callbacks: called=%v empty=%v", called, empty)
		}
	}
}
