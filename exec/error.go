package exec

import (
	"errors"
	"strings"

	"github.com/nikolalohinski/gonja/v2/nodes"
)

type RuntimeError struct {
	Message string
	Raw     string
	Line    int
	Column  int
	Err     error
}

func (e *RuntimeError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *RuntimeError) Unwrap() error {
	return e.Err
}

func describe(err error) string {
	context := ""
	for {
		inner := errors.Unwrap(err)
		if inner == nil {
			break
		}
		full := err.Error()
		if own, ok := strings.CutSuffix(full, ": "+inner.Error()); ok {
			context = own
		}
		err = inner
	}
	if context == "" {
		return err.Error()
	}
	return context + ": " + err.Error()
}

func newRuntimeError(err error, node nodes.Node) error {
	if err == nil {
		return nil
	}

	runtimeErr := &RuntimeError{
		Message: describe(err),
		Raw:     node.String(),
		Err:     err,
	}

	if pos := node.Position(); pos != nil {
		runtimeErr.Line = pos.Line
		runtimeErr.Column = pos.Col
	}

	return runtimeErr
}
