package tokens_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/config"
	"github.com/nikolalohinski/gonja/tokens"
	"github.com/stretchr/testify/assert"
)

type tok struct {
	typ tokens.Type
	val string
}

func (t tok) String() string {
	return `"` + t.val + `"`
}

var (
	EOF            = tok{tokens.EOF, ""}
	varBegin       = tok{tokens.VariableBegin, "{{"}
	varEnd         = tok{tokens.VariableEnd, "}}"}
	blockBegin     = tok{tokens.BlockBegin, "{%"}
	blockBeginTrim = tok{tokens.BlockBegin, "{%-"}
	blockEnd       = tok{tokens.BlockEnd, "%}"}
	blockEndTrim   = tok{tokens.BlockEnd, "-%}"}
	lParen         = tok{tokens.Lparen, "("}
	rParen         = tok{tokens.Rparen, ")"}
	lBrace         = tok{tokens.Lbrace, "{"}
	rBrace         = tok{tokens.Rbrace, "}"}
	lBracket       = tok{tokens.Lbracket, "["}
	rBracket       = tok{tokens.Rbracket, "]"}
	space          = tok{tokens.Whitespace, " "}
)

func data(text string) tok {
	return tok{tokens.Data, text}
}

func name(text string) tok {
	return tok{tokens.Name, text}
}

func str(text string) tok {
	return tok{tokens.String, text}
}

func error(text string) tok {
	return tok{tokens.Error, text}
}

var lexerCases = []struct {
	name     string
	input    string
	expected []tok
}{
	{"empty", "", []tok{EOF}},
	{"data", "Hello World", []tok{
		data("Hello World"),
		EOF,
	}},
	{"comment", "{# a comment #}", []tok{
		tok{tokens.CommentBegin, "{#"},
		data(" a comment "),
		tok{tokens.CommentEnd, "#}"},
		EOF,
	}},
	{"mixed comment", "Hello, {# comment #}World", []tok{
		data("Hello, "),
		tok{tokens.CommentBegin, "{#"},
		data(" comment "),
		tok{tokens.CommentEnd, "#}"},
		data("World"),
		EOF,
	}},
	{"simple variable", "{{ foo }}", []tok{
		varBegin,
		space,
		name("foo"),
		space,
		varEnd,
		EOF,
	}},
	{"basic math expression", "{{ (a - b) + c }}", []tok{
		varBegin, space,
		lParen, name("a"), space, tok{tokens.Sub, "-"}, space, name("b"), rParen,
		space, tok{tokens.Add, "+"}, space, name("c"),
		space, varEnd,
		EOF,
	}},
	{"blocks", "Hello.  {% if true %}World{% else %}Nobody{% endif %}", []tok{
		data("Hello.  "),
		blockBegin, space, name("if"), space, name("true"), space, blockEnd,
		data("World"),
		blockBegin, space, name("else"), space, blockEnd,
		data("Nobody"),
		blockBegin, space, name("endif"), space, blockEnd,
		EOF,
	}},
	{"blocks with trim control", "Hello.  {%- if true -%}World{%- else -%}Nobody{%- endif -%}", []tok{
		data("Hello.  "),
		blockBeginTrim, space, name("if"), space, name("true"), space, blockEndTrim,
		data("World"),
		blockBeginTrim, space, name("else"), space, blockEndTrim,
		data("Nobody"),
		blockBeginTrim, space, name("endif"), space, blockEndTrim,
		EOF,
	}},
	{"Ignore tags in comment", "<html>{# ignore {% tags %} in comments ##}</html>", []tok{
		data("<html>"),
		tok{tokens.CommentBegin, "{#"},
		data(" ignore {% tags %} in comments #"),
		tok{tokens.CommentEnd, "#}"},
		data("</html>"),
		EOF,
	}},
	{"Mixed content", "{# comment #}{% if foo -%} bar {%- elif baz %} bing{%endif    %}", []tok{
		tok{tokens.CommentBegin, "{#"},
		data(" comment "),
		tok{tokens.CommentEnd, "#}"},
		blockBegin, space, name("if"), space, name("foo"), space, blockEndTrim,
		data(" bar "),
		blockBeginTrim, space, name("elif"), space, name("baz"), space, blockEnd,
		data(" bing"),
		blockBegin, name("endif"), tok{tokens.Whitespace, "    "}, blockEnd,
		EOF,
	}},
	{"mixed tokens with doubles", "{{ +--+ /+//,|*/**=>>=<=< == }}", []tok{
		varBegin,
		space,
		tok{tokens.Add, "+"}, tok{tokens.Sub, "-"}, tok{tokens.Sub, "-"}, tok{tokens.Add, "+"},
		space,
		tok{tokens.Div, "/"}, tok{tokens.Add, "+"}, tok{tokens.Floordiv, "//"},
		tok{tokens.Comma, ","},
		tok{tokens.Pipe, "|"},
		tok{tokens.Mul, "*"},
		tok{tokens.Div, "/"},
		tok{tokens.Pow, "**"},
		tok{tokens.Assign, "="},
		tok{tokens.Gt, ">"},
		tok{tokens.Gteq, ">="},
		tok{tokens.Lteq, "<="},
		tok{tokens.Lt, "<"},
		space,
		tok{tokens.Eq, "=="},
		space,
		varEnd,
		EOF,
	}},
	{"delimiters", "{{ ([{}]()) }}", []tok{
		varBegin, space,
		lParen, lBracket, lBrace, rBrace, rBracket, lParen, rParen, rParen,
		space, varEnd,
		EOF,
	}},
	{"Unbalanced delimiters", "{{ ([{]) }}", []tok{
		varBegin, space,
		lParen, lBracket, lBrace,
		error(`Unbalanced delimiters, expected "}", got "]"`),
	}},
	{"Unexpeced delimiter", "{{ ()) }}", []tok{
		varBegin, space,
		lParen, rParen,
		error(`Unexpected delimiter ")"`),
	}},
	{"Unbalance over end block", "{{ ({a:b, {a:b}}) }}", []tok{
		varBegin, space,
		lParen,
		lBrace, name("a"), tok{tokens.Colon, ":"}, name("b"), tok{tokens.Comma, ","},
		space,
		lBrace, name("a"), tok{tokens.Colon, ":"}, name("b"), rBrace, rBrace,
		rParen,
		space, varEnd,
		EOF,
	}},
	{"string with double quote", `{{ "Hello, " + "World" }}`, []tok{
		varBegin, space,
		str("Hello, "),
		space, tok{tokens.Add, "+"}, space,
		str("World"),
		space, varEnd,
		EOF,
	}},
	{"string with simple quote", `{{ 'Hello, ' + 'World' }}`, []tok{
		varBegin, space,
		str("Hello, "),
		space, tok{tokens.Add, "+"}, space,
		str("World"),
		space, varEnd,
		EOF,
	}},
	{"single quotes inside double quotes string", `{{ "'quoted' test" }}`, []tok{
		varBegin, space, str("'quoted' test"), space, varEnd, EOF,
	}},
	{"escaped string", `{{ "Hello, \"World\"" }}`, []tok{
		varBegin, space,
		str(`Hello, "World"`),
		space, varEnd,
		EOF,
	}},
	{"escaped string mixed", `{{ "Hello,\n \'World\'" }}`, []tok{
		varBegin, space,
		str(`Hello,\n 'World'`),
		space, varEnd,
		EOF,
	}},
	{"if statement", `{% if 5.5 == 5.500000 %}5.5 is 5.500000{% endif %}`, []tok{
		blockBegin, space, name("if"), space,
		tok{tokens.Float, "5.5"}, space, tok{tokens.Eq, "=="}, space, tok{tokens.Float, "5.500000"},
		space, blockEnd,
		data("5.5 is 5.500000"),
		blockBegin, space, name("endif"), space, blockEnd,
		EOF,
	}},
	{"logical 'and' expression", `{{ a is defined and a == "x" }}`, []tok{
		varBegin,
		space, name("a"), space, name("is"), space, name("defined"),
		space, {tokens.And, "and"},
		space, name("a"), space, {tokens.Eq, "=="}, space, str("x"), space,
		varEnd,
		EOF,
	}},
	{"logical 'or' expression", `{{ a is defined or a == "x" }}`, []tok{
		varBegin,
		space, name("a"), space, name("is"), space, name("defined"),
		space, {tokens.Or, "or"},
		space, name("a"), space, {tokens.Eq, "=="}, space, str("x"), space,
		varEnd,
		EOF,
	}},
	{"logical 'not' expression", `{{ not a }}`, []tok{
		varBegin,
		space, {tokens.Not, "not"},
		space, name("a"), space,
		varEnd,
		EOF,
	}},
}

func tokenSlice(c chan *tokens.Token) []*tokens.Token {
	toks := []*tokens.Token{}
	for token := range c {
		toks = append(toks, token)
	}
	return toks
}

func TestLexer(t *testing.T) {
	for _, lc := range lexerCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			lexer := tokens.NewLexer(test.input)
			go lexer.Run()
			toks := tokenSlice(lexer.Tokens)

			assert := assert.New(t)
			assert.Equal(len(test.expected), len(toks))
			actual := []tok{}
			for _, token := range toks {
				actual = append(actual, tok{token.Type, token.Val})
			}
			assert.Equal(test.expected, actual)
		})
	}
}

func streamResult(s *tokens.Stream) []tok {
	out := []tok{}
	for !s.End() {
		token := s.Current()
		out = append(out, tok{token.Type, token.Val})
		s.Next()
	}
	return out
}

func asStreamResult(toks []tok) ([]tok, bool) {
	out := []tok{}
	isError := false
	for _, token := range toks {
		if token.typ == tokens.Error {
			isError = true
			break
		}
		if token.typ != tokens.Whitespace && token.typ != tokens.EOF {
			out = append(out, token)
		}
	}
	return out, isError
}

func TestLex(t *testing.T) {
	for _, lc := range lexerCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			stream := tokens.Lex(test.input)
			expected, _ := asStreamResult(test.expected)

			actual := streamResult(stream)

			assert := assert.New(t)
			assert.Equal(len(expected), len(actual))
			assert.Equal(expected, actual)
		})
	}
}
func TestRegexpClashingDelimiters(t *testing.T) {
	t.Run("Expression delimiters clashing the regexp parsing", func(t *testing.T) {

		config.DefaultConfig.VariableStartString = "[["
		config.DefaultConfig.VariableEndString = "]]"
		config.DefaultConfig.BlockStartString = "[%"
		config.DefaultConfig.BlockEndString = "%]"
		config.DefaultConfig.CommentStartString = "[#"
		config.DefaultConfig.CommentEndString = "#]"
		defer func() {
			config.DefaultConfig = config.NewConfig()
		}()

		lexer := tokens.NewLexer("[[ variable ]][% block %][# comment #]")
		go lexer.Run()
		toks := tokenSlice(lexer.Tokens)

		stream := tokens.NewStream(toks)
		expected, _ := asStreamResult([]tok{
			tok{tokens.VariableBegin, "[["}, space, name("variable"), space, tok{tokens.VariableEnd, "]]"},
			tok{tokens.BlockBegin, "[%"}, space, name("block"), space, tok{tokens.BlockEnd, "%]"},
			tok{tokens.CommentBegin, "[#"}, data(" comment "), tok{tokens.CommentEnd, "#]"},
			EOF,
		})

		actual := streamResult(stream)

		assert := assert.New(t)
		assert.Equal(len(expected), len(actual))
		assert.Equal(expected, actual)
	})
}
func TestStreamSlice(t *testing.T) {
	for _, lc := range lexerCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			lexer := tokens.NewLexer(test.input)
			go lexer.Run()
			toks := tokenSlice(lexer.Tokens)

			stream := tokens.NewStream(toks)
			expected, _ := asStreamResult(test.expected)

			actual := streamResult(stream)

			assert := assert.New(t)
			assert.Equal(len(expected), len(actual))
			assert.Equal(expected, actual)
		})
	}
}

const positionsCase = `Hello
{#
    Multiline comment
#}
World
`

func TestLexerPosition(t *testing.T) {
	assert := assert.New(t)

	lexer := tokens.NewLexer(positionsCase)
	go lexer.Run()
	toks := tokenSlice(lexer.Tokens)
	assert.Equal([]*tokens.Token{
		{tokens.Data, "Hello\n", 0, 1, 1, false},
		{tokens.CommentBegin, "{#", 6, 2, 1, false},
		{tokens.Data, "\n    Multiline comment\n", 8, 2, 3, false},
		{tokens.CommentEnd, "#}", 31, 4, 1, false},
		{tokens.Data, "\nWorld\n", 33, 4, 3, false},
		{tokens.EOF, "", 40, 6, 1, false},
	}, toks)
}
