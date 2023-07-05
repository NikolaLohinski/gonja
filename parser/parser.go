package parser

import (
	"fmt"
	"strings"

	"github.com/nikolalohinski/gonja/config"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/tokens"
)

// The parser provides you a comprehensive and easy tool to
// work with the template document and arguments provided by
// the user for your custom tag.
//
// The parser works on a token list which will be provided by gonja.
// A token is a unit you can work with. Tokens are either of type identifier,
// string, number, keyword, HTML or symbol.
//
// (See Token's documentation for more about tokens)
type Parser struct {
	Name   string
	Stream *tokens.Stream
	Config *config.Config

	Template       *nodes.Template
	Statements     map[string]StatementParser
	Level          int8
	TemplateParser TemplateParser
}

// Creates a new parser to parse tokens.
// Used inside gonja to parse documents and to provide an easy-to-use
// parser for tag authors
func NewParser(name string, cfg *config.Config, stream *tokens.Stream) *Parser {
	return &Parser{
		Name:   name,
		Stream: stream,
		Config: cfg,
	}
}

func Parse(input string) (*nodes.Template, error) {
	stream := tokens.Lex(input)
	p := NewParser("parser", config.DefaultConfig, stream)
	return p.Parse()
}

func (p *Parser) Parse() (*nodes.Template, error) {
	return p.ParseTemplate()
}

// Consume one token. It will be gone forever.
func (p *Parser) Consume() {
	p.Stream.Next()
}

// Next returns and consume the current token
func (p *Parser) Next() *tokens.Token {
	// t := p.Stream.Next()
	// p.Consume()
	// return t
	return p.Stream.Next()
}

func (p *Parser) End() bool {
	return p.Stream.End()
}

// Match returns the CURRENT token if the given type matches.
// Consumes this token on success.
func (p *Parser) Match(types ...tokens.Type) *tokens.Token {
	tok := p.Stream.Current()
	for _, t := range types {
		if tok.Type == t {
			p.Stream.Next()
			return tok
		}
	}
	return nil
}

func (p *Parser) MatchName(names ...string) *tokens.Token {
	t := p.Current(tokens.Name)
	if t != nil {
		for _, name := range names {
			if t.Val == name {
				return p.Pop()
			}
		}
	}
	return nil
}

// Pop returns the current token and advance to the next
func (p *Parser) Pop() *tokens.Token {
	t := p.Stream.Current()
	p.Stream.Next()
	return t
}

// Current returns the current token without consuming
// it and only if it matches one of the given types
func (p *Parser) Current(types ...tokens.Type) *tokens.Token {
	tok := p.Stream.Current()
	if types == nil {
		return tok
	}
	for _, t := range types {
		if tok.Type == t {
			return tok
		}
	}
	return nil
}

func (p *Parser) Peek(types ...tokens.Type) *tokens.Token {
	tok := p.Stream.Peek()
	if types == nil {
		return tok
	}
	for _, t := range types {
		if tok.Type == t {
			return tok
		}
	}
	return nil
}

func (p *Parser) CurrentName(names ...string) *tokens.Token {
	t := p.Current(tokens.Name)
	if t != nil {
		for _, name := range names {
			if t.Val == name {
				return t
			}
		}
	}
	return nil
}

// WrapUntil wraps all nodes between starting tag and "{% endtag %}" and provides
// one simple interface to execute the wrapped nodes.
// It returns a parser to process provided arguments to the tag.
func (p *Parser) WrapUntil(names ...string) (*nodes.Wrapper, *Parser, error) {
	wrapper := &nodes.Wrapper{
		Location: p.Current(),
		Trim:     &nodes.Trim{},
	}

	var args []*tokens.Token

	for !p.Stream.End() {
		// New tag, check whether we have to stop wrapping here
		if begin := p.Match(tokens.BlockBegin); begin != nil {
			endTag := p.CurrentName(names...)

			if endTag != nil {
				p.Consume()
				for {
					if end := p.Match(tokens.BlockEnd); end != nil {
						wrapper.EndTag = endTag.Val
						if data := p.Current(tokens.Data); data != nil {
							data.Trim = data.Trim || len(end.Val) > 0 && end.Val[0] == '-'
						}
						stream := tokens.NewStream(args)
						return wrapper, NewParser(p.Name, p.Config, stream), nil
					}
					if p.End() || p.Current(tokens.EOF) != nil {
						return nil, nil, p.Error("Unexpected EOF.", p.Current())
					}
					args = append(args, p.Next())
				}
			}
			p.Stream.Backup()
		}

		// Otherwise process next element to be wrapped
		node, err := p.parseDocElement()
		if err != nil {
			return nil, nil, err
		}
		wrapper.Nodes = append(wrapper.Nodes, node)
	}

	return nil, nil, p.Error(fmt.Sprintf("Unexpected EOF, expected tag %s.", strings.Join(names, " or ")),
		p.Current())
}
