package controlStructures

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

// Sentinel errors used by for-loop to detect break/continue.
type LoopBreakError struct{}

func (e *LoopBreakError) Error() string { return "break" }

type LoopContinueError struct{}

func (e *LoopContinueError) Error() string { return "continue" }

type BreakControlStructure struct {
	Location *tokens.Token
}

func (c *BreakControlStructure) Position() *tokens.Token { return c.Location }
func (c *BreakControlStructure) String() string {
	t := c.Position()
	return fmt.Sprintf("BreakControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}
func (stmt *BreakControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	return &LoopBreakError{}
}

type ContinueControlStructure struct {
	Location *tokens.Token
}

func (c *ContinueControlStructure) Position() *tokens.Token { return c.Location }
func (c *ContinueControlStructure) String() string {
	t := c.Position()
	return fmt.Sprintf("ContinueControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}
func (stmt *ContinueControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	return &LoopContinueError{}
}

func breakParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	ctrl := &BreakControlStructure{Location: p.Current()}
	if !args.End() {
		return nil, args.Error("break does not accept arguments", args.Current())
	}
	return ctrl, nil
}

func continueParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	ctrl := &ContinueControlStructure{Location: p.Current()}
	if !args.End() {
		return nil, args.Error("continue does not accept arguments", args.Current())
	}
	return ctrl, nil
}
