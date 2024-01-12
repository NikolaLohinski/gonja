package controlStructures

import (
	"fmt"
	"io"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type RawControlStructure struct {
	data *nodes.Data
}

func (controlStructure *RawControlStructure) Position() *tokens.Token {
	return controlStructure.data.Position()
}
func (controlStructure *RawControlStructure) String() string {
	t := controlStructure.Position()
	return fmt.Sprintf("RawControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (controlStructure *RawControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	_, err := io.WriteString(r.Output, controlStructure.data.Data.Val)
	return err
}

func rawParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	controlStructure := &RawControlStructure{}

	wrapper, _, err := p.WrapUntil("endraw")
	if err != nil {
		return nil, err
	}
	node := wrapper.Nodes[0]
	data, ok := node.(*nodes.Data)
	if ok {
		controlStructure.data = data
	} else {
		return nil, p.Error("raw controlStructure can only contains a single data node", node.Position())
	}

	if !args.End() {
		return nil, args.Error("raw controlStructure doesn't accept parameters.", args.Current())
	}

	return controlStructure, nil
}
