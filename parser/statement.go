package parser

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type StatementParser func(parser *Parser, args *Parser) (nodes.Statement, error)

func (p *Parser) ParseStatementBlock() (*nodes.StatementBlock, error) {
	log.WithFields(log.Fields{
		"current": p.Current(),
	}).Trace("ParseStatementBlock")

	begin := p.Match(tokens.BlockBegin)
	if begin == nil {
		return nil, errors.Errorf(`Expected "%s" got "%s"`, p.Config.BlockStartString, p.Current())
	}

	name := p.Match(tokens.Name)
	if name == nil {
		return nil, p.Error("Expected a statement name here", p.Current())
	}

	stmtParser, exists := p.statements[name.Val]
	if !exists {
		return nil, p.Error(fmt.Sprintf("Statement '%s' not found (or beginning not provided)", name.Val), name)
	}

	log.Trace("args")
	var args []*tokens.Token
	for p.Current(tokens.BlockEnd) == nil && !p.Stream().End() {
		log.Trace("for args")
		args = append(args, p.Next())
	}
	log.Trace("loop ended")

	end := p.Match(tokens.BlockEnd)
	if end == nil {
		return nil, p.Error(fmt.Sprintf(`Expected end of block "%s"`, p.Config.BlockEndString), p.Current())
	}
	if data := p.Current(tokens.Data); data != nil {
		data.Trim = data.Trim || len(end.Val) > 0 && end.Val[0] == '-'
	}

	log.WithFields(log.Fields{
		"args": args,
	}).Trace("Matched end block")

	stream := tokens.NewStream(args)
	log.WithFields(log.Fields{
		"stream": stream,
	}).Trace("Got stream")
	argParser := NewParser(p.identifier, stream, p.Config, p.Loader, p.statements)
	log.Trace("argparser")

	stmt, err := stmtParser(p, argParser)
	if err != nil {
		return nil, errors.Wrapf(err, `Unable to parse statement "%s"`, name.Val)
	}
	log.Trace("got stmt and return")
	return &nodes.StatementBlock{
		Location: begin,
		Name:     name.Val,
		Stmt:     stmt,
	}, nil
}
