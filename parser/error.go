package parser

import (
	"errors"
	"fmt"

	"github.com/nikolalohinski/gonja/v2/tokens"
)

type SyntaxError struct {
	Message string
	Raw      string
	Line     int
	Column   int
}

func (e *SyntaxError) Error() string {
	if e.Raw == "" && e.Line == 0 && e.Column == 0 {
		return e.Message
	}

	return fmt.Sprintf(`%s (Line: %d Col: %d, near "%s")`, e.Message, e.Line, e.Column, e.Raw)
}

func (p *Parser) Error(message string, token *tokens.Token) error {
	if token == nil {
		return errors.New(message)
	}

	return &SyntaxError{
		Message: message,
		Raw:     token.Val,
		Line:    token.Line,
		Column:  token.Col,
	}
}
