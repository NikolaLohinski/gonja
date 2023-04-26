package exec

/* Incomplete:
   -----------

   verbatim (only the "name" argument is missing for verbatim)

   Reconsideration:
   ----------------

   debug (reason: not sure what to output yet)
   regroup / Grouping on other properties (reason: maybe too python-specific; not sure how useful this would be in Go)

   Following built-in tags wont be added:
   --------------------------------------

   csrf_token (reason: web-framework specific)
   load (reason: python-specific)
   url (reason: web-framework specific)
*/

import (
	// "fmt"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	// "github.com/nikolalohinski/gonja/tokens"
)

// type NodeStatement interface {
// 	astNode
// }

// This is the function signature of the tag's parser you will have
// to implement in order to create a new tag.
//
// 'doc' is providing access to the whole document while 'arguments'
// is providing access to the user's arguments to the tag:
//
//     {% your_tag_name some "arguments" 123 %}
//
// start_token will be the *Token with the tag's name in it (here: your_tag_name).
//
// Please see the Parser documentation on how to use the parser.
// See RegisterTag()'s documentation for more information about
// writing a tag as well.

// type StatementExecutor func(*nodes.Node, *ExecutionContext) *Value

type Statement interface {
	nodes.Statement
	Execute(*Renderer, *nodes.StatementBlock) error
}

type StatementSet map[string]parser.StatementParser

// Exists returns true if the given test is already registered
func (ss StatementSet) Exists(name string) bool {
	_, existing := ss[name]
	return existing
}

// Registers a new tag. You usually want to call this
// function in the tag's init() function:
// http://golang.org/doc/effective_go.html#init
func (ss *StatementSet) Register(name string, parser parser.StatementParser) error {
	if ss.Exists(name) {
		return errors.Errorf("Statement '%s' is already registered", name)
	}
	(*ss)[name] = parser
	// &statement{
	// 	name:   name,
	// 	parser: parserFn,
	// }
	return nil
}

// Replaces an already registered tag with a new implementation. Use this
// function with caution since it allows you to change existing tag behaviour.
func (ss *StatementSet) Replace(name string, parser parser.StatementParser) error {
	if !ss.Exists(name) {
		return errors.Errorf("Statement '%s' does not exist (therefore cannot be overridden)", name)
	}
	(*ss)[name] = parser
	// statements[name] = &statement{
	// 	name:   name,
	// 	parser: parserFn,
	// }
	return nil
}

func (ss *StatementSet) Update(other StatementSet) StatementSet {
	for name, parser := range other {
		(*ss)[name] = parser
	}
	return *ss
}

// func (ss StatementSet) Parsers() map[string]parser.StatementParser {
// 	parsers := map[string]parser.StatementParser{}
// 	for key, specs := range ss {
// 		parsers[key] = specs.Parse
// 	}
// 	return parsers
// }

// // Tag = "{%" IDENT ARGS "%}"
// func (p *Parser) ParseStatement() (ast.Statement, *Error) {
// 	p.Consume() // consume "{%"
// 	tokenName := p.MatchType(TokenIdentifier)

// 	// Check for identifier
// 	if tokenName == nil {
// 		return nil, p.Error("Statement name must be an identifier.", nil)
// 	}

// 	// Check for the existing statement
// 	stmt, exists := statements[tokenName.Val]
// 	if !exists {
// 		// Does not exists
// 		return nil, p.Error(fmt.Sprintf("Statement '%s' not found (or beginning not provided)", tokenName.Val), tokenName)
// 	}

// 	// Check sandbox tag restriction
// 	if _, isBanned := p.bannedStmts[tokenName.Val]; isBanned {
// 		return nil, p.Error(fmt.Sprintf("Usage of statement '%s' is not allowed (sandbox restriction active).", tokenName.Val), tokenName)
// 	}

// 	var argsToken []*Token
// 	for p.Peek(TokenSymbol, "%}") == nil && p.Remaining() > 0 {
// 		// Add token to args
// 		argsToken = append(argsToken, p.Current())
// 		p.Consume() // next token
// 	}

// 	// EOF?
// 	if p.Remaining() == 0 {
// 		return nil, p.Error("Unexpectedly reached EOF, no statement end found.", p.lastToken)
// 	}

// 	p.Match(TokenSymbol, "%}")

// 	argParser := newParser(p.name, argsToken, p.template)
// 	if len(argsToken) == 0 {
// 		// This is done to have nice EOF error messages
// 		argParser.lastToken = tokenName
// 	}

// 	p.template.level++
// 	defer func() { p.template.level-- }()
// 	return stmt.parser(p, tokenName, argParser)
// }
