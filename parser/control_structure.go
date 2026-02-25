package parser

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type ControlStructureParser func(parser *Parser, args *Parser) (nodes.ControlStructure, error)

func (p *Parser) ParseControlStructureBlock() (*nodes.ControlStructureBlock, error) {
	begin := p.Match(tokens.BlockBegin)
	if begin == nil {
		return nil, errors.Errorf(`Expected "%s" got "%s"`, p.Config.BlockStartString, p.Current())
	}

	name := p.Match(tokens.Name)
	if name == nil {
		return nil, p.Error("Expected a controlStructure name here", p.Current())
	}

	controlStructureParser, exists := p.controlStructures.Get(name.Val)
	if !exists {
		return nil, p.Error(fmt.Sprintf("ControlStructure '%s' not found (or beginning not provided)", name.Val), name)
	}

	var args []*tokens.Token
	for p.Current(tokens.BlockEnd) == nil && !p.Stream().End() {
		args = append(args, p.Next())
	}

	end := p.Match(tokens.BlockEnd)
	if end == nil {
		return nil, p.Error(fmt.Sprintf(`Expected end of block "%s"`, p.Config.BlockEndString), p.Current())
	}
	if data := p.Current(tokens.Data); data != nil {
		data.Trim = data.Trim || len(end.Val) > 0 && end.Val[0] == '-'
		data.RemoveFirstLineReturn = p.Config.TrimBlocks && len(end.Val) > 0 && end.Val[0] != '+'
	}

	stream := tokens.NewStream(args)
	argParser := NewParser(p.identifier, stream, p.Config, p.Loader, p.controlStructures)

	controlStructure, err := controlStructureParser(p, argParser)
	if err != nil {
		return nil, errors.Wrapf(err, `Unable to parse controlStructure "%s"`, name.Val)
	}
	return &nodes.ControlStructureBlock{
		Location:         begin,
		Name:             name.Val,
		ControlStructure: controlStructure,
	}, nil
}
