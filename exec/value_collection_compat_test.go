package exec_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2/exec"
)

type testAttributeGetter struct{}

func (testAttributeGetter) GetAttribute(name string) (*exec.Value, bool) {
	if name == "name" {
		return exec.AsValue("alice"), true
	}
	return exec.AsValue(nil), false
}

type testInt64Value struct{}

func (testInt64Value) Int64() int64 {
	return 42
}

type testFloat64Value struct{}

func (testFloat64Value) Float64() float64 {
	return 3.5
}

func TestValueCompatibilityHelpers(t *testing.T) {
	t.Run("attribute getter is used", func(t *testing.T) {
		value, ok := exec.AsValue(testAttributeGetter{}).GetAttribute("name")
		if !ok {
			t.Fatal("expected custom attribute getter to resolve the attribute")
		}
		if got := value.String(); got != "alice" {
			t.Fatalf("unexpected attribute value: %q", got)
		}
	})

	t.Run("integer parses strings and custom int64 values", func(t *testing.T) {
		if got := exec.AsValue("12.9").Integer(); got != 12 {
			t.Fatalf("expected string conversion to int, got %d", got)
		}
		if got := exec.AsValue(testInt64Value{}).Integer(); got != 42 {
			t.Fatalf("expected custom Int64 conversion, got %d", got)
		}
	})

	t.Run("float parses strings and custom float64 values", func(t *testing.T) {
		if got := exec.AsValue("12.5").Float(); got != 12.5 {
			t.Fatalf("expected string conversion to float, got %v", got)
		}
		if got := exec.AsValue(testFloat64Value{}).Float(); got != 3.5 {
			t.Fatalf("expected custom Float64 conversion, got %v", got)
		}
	})
}
