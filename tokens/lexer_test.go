package tokens_test

import (
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/tokens"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Context("lexer", func() {
	var (
		lexer = new(*tokens.Lexer)

		returnedTokens = new([]*tokens.Token)
	)

	JustBeforeEach(func() {
		go (*lexer).Run()

		*returnedTokens = make([]*tokens.Token, 0)
		for message := range (*lexer).Tokens {
			*returnedTokens = append(*returnedTokens, message)
		}
	})
	Context("default", func() {
		for _, testCase := range []struct {
			description  string
			input        string
			tokenMatcher []Fields
		}{
			{
				"is empty",
				"",
				[]Fields{
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"unbalanced expression delimiter",
				"{{ 42",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Integer), "Val": Equal("42")},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"just contains regular data",
				"Hello World",
				[]Fields{
					{"Type": Equal(tokens.Data), "Val": Equal("Hello World")},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a comment",
				"{# a comment #}",
				[]Fields{
					{"Type": Equal(tokens.CommentBegin)},
					{"Type": Equal(tokens.Data), "Val": Equal(" a comment ")},
					{"Type": Equal(tokens.CommentEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"has data and a comment",
				"Hello, {# comment #}World",
				[]Fields{
					{"Type": Equal(tokens.Data), "Val": Equal("Hello, ")},
					{"Type": Equal(tokens.CommentBegin)},
					{"Type": Equal(tokens.Data), "Val": Equal(" comment ")},
					{"Type": Equal(tokens.CommentEnd)},
					{"Type": Equal(tokens.Data), "Val": Equal("World")},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a simple variable",
				"{{ foo }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("foo")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a simple math expression",
				"{{ (a - b) + c }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.LeftParenthesis)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Subtraction)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("b")},
					{"Type": Equal(tokens.RightParenthesis)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Addition)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("c")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"contains numeric literals with underscores, exponents and prefixes",
				"{{ 12_34_56 1_2.3_4e5_6 0b1001_1111 0o1_23 0x12_3abc }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Integer), "Val": Equal("12_34_56")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Float), "Val": Equal("1_2.3_4e5_6")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Integer), "Val": Equal("0b1001_1111")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Integer), "Val": Equal("0o1_23")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Integer), "Val": Equal("0x12_3abc")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"contains a scientific float with a signed exponent",
				"{{ 25.6e-10 }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Float), "Val": Equal("25.6e-10")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"keeps float method access separate from the numeric token",
				"{{ 42.00.hex() }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Float), "Val": Equal("42.00")},
					{"Type": Equal(tokens.Dot)},
					{"Type": Equal(tokens.Name), "Val": Equal("hex")},
					{"Type": Equal(tokens.LeftParenthesis)},
					{"Type": Equal(tokens.RightParenthesis)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"contains data and and 'if ... else' block",
				"Hello.  {% if true %}World{% else %}Nobody{% endif %}",
				[]Fields{
					{"Type": Equal(tokens.Data), "Val": Equal("Hello.  ")},
					{"Type": Equal(tokens.BlockBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("if")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("true")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd)},
					{"Type": Equal(tokens.Data), "Val": Equal("World")},
					{"Type": Equal(tokens.BlockBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("else")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd)},
					{"Type": Equal(tokens.Data), "Val": Equal("Nobody")},
					{"Type": Equal(tokens.BlockBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("endif")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a block with trim control",
				"come      {%- if true -%}  -  {%- endif -%}closer",
				[]Fields{
					{"Type": Equal(tokens.Data), "Val": Equal("come      ")},
					{"Type": Equal(tokens.BlockBegin), "Val": MatchRegexp("^.*-$")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("if")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("true")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd), "Val": MatchRegexp("^-.*$")},
					{"Type": Equal(tokens.Data), "Val": Equal("  -  ")},
					{"Type": Equal(tokens.BlockBegin), "Val": MatchRegexp("^.*-$")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("endif")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd), "Val": MatchRegexp("^-.*$")},
					{"Type": Equal(tokens.Data), "Val": Equal("closer")},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a comment with a nested block",
				"<html>{# ignore {% tags %} in comments ##}</html>",
				[]Fields{
					{"Type": Equal(tokens.Data), "Val": Equal("<html>")},
					{"Type": Equal(tokens.CommentBegin)},
					{"Type": Equal(tokens.Data), "Val": Equal(" ignore {% tags %} in comments #")},
					{"Type": Equal(tokens.CommentEnd)},
					{"Type": Equal(tokens.Data), "Val": Equal("</html>")},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a little bit of everything",
				"{# comment #}{% if foo -%} bar {%- elif baz %} bing{%endif    %}",
				[]Fields{
					{"Type": Equal(tokens.CommentBegin)},
					{"Type": Equal(tokens.Data), "Val": Equal(" comment ")},
					{"Type": Equal(tokens.CommentEnd)},
					{"Type": Equal(tokens.BlockBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("if")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("foo")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd), "Val": MatchRegexp("^-.*$")},
					{"Type": Equal(tokens.Data), "Val": Equal(" bar ")},
					{"Type": Equal(tokens.BlockBegin), "Val": MatchRegexp("^.*-$")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("elif")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("baz")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd)},
					{"Type": Equal(tokens.Data), "Val": Equal(" bing")},
					{"Type": Equal(tokens.BlockBegin)},
					{"Type": Equal(tokens.Name), "Val": Equal("endif")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.BlockEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"contains all possible operators",
				"{{ +--+ /+//,|*/**=>>=<=< == % }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Addition)},
					{"Type": Equal(tokens.Subtraction)},
					{"Type": Equal(tokens.Subtraction)},
					{"Type": Equal(tokens.Addition)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Division)},
					{"Type": Equal(tokens.Addition)},
					{"Type": Equal(tokens.FloorDivision)},
					{"Type": Equal(tokens.Comma)},
					{"Type": Equal(tokens.Pipe)},
					{"Type": Equal(tokens.Multiply)},
					{"Type": Equal(tokens.Division)},
					{"Type": Equal(tokens.Power)},
					{"Type": Equal(tokens.Assign)},
					{"Type": Equal(tokens.GreaterThan)},
					{"Type": Equal(tokens.GreaterThanOrEqual)},
					{"Type": Equal(tokens.LowerThanOrEqual)},
					{"Type": Equal(tokens.LowerThan)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Equals)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Modulo)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"contains all possible delimiters",
				"{{ ([{}]()) }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.LeftParenthesis)},
					{"Type": Equal(tokens.LeftBracket)},
					{"Type": Equal(tokens.LeftBrace)},
					{"Type": Equal(tokens.RightBrace)},
					{"Type": Equal(tokens.RightBracket)},
					{"Type": Equal(tokens.LeftParenthesis)},
					{"Type": Equal(tokens.RightParenthesis)},
					{"Type": Equal(tokens.RightParenthesis)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"contains unbalanced delimiters",
				"{{ ([{]) }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.LeftParenthesis)},
					{"Type": Equal(tokens.LeftBracket)},
					{"Type": Equal(tokens.LeftBrace)},
					{"Type": Equal(tokens.Error), "Val": Equal(`Unbalanced delimiters, expected "}", got "]"`)},
				},
			},
			{
				"contains an unexpected delimiter",
				"{{ ()) }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.LeftParenthesis)},
					{"Type": Equal(tokens.RightParenthesis)},
					{"Type": Equal(tokens.Error), "Val": Equal(`Unexpected delimiter ")"`)},
				},
			},
			{
				"contains what looks likes to be an unbalanced variable block but is actually valid",
				"{{ ({a:b, {a:b}}) }}",
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.LeftParenthesis)},
					{"Type": Equal(tokens.LeftBrace)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Colon)},
					{"Type": Equal(tokens.Name), "Val": Equal("b")},
					{"Type": Equal(tokens.Comma)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.LeftBrace)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Colon)},
					{"Type": Equal(tokens.Name), "Val": Equal("b")},
					{"Type": Equal(tokens.RightBrace)},
					{"Type": Equal(tokens.RightBrace)},
					{"Type": Equal(tokens.RightParenthesis)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a variable block with single and double quoted strings",
				`{{ "Hello, " + 'World' }}`,
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.String), "Val": Equal("Hello, ")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Addition)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.String), "Val": Equal("World")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a variable block with a double quoted string containing single quotes",
				`{{ "foo 'bar'" }}`,
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.String), "Val": Equal("foo 'bar'")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a variable block with a single quoted string containing double quotes",
				`{{ 'foo "bar"' }}`,
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.String), "Val": Equal(`foo "bar"`)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a variable block with an escaped double quoted string",
				`{{ "foo \"bar\"" }}`,
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.String), "Val": Equal(`foo "bar"`)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a logical 'and' expression",
				`{{ a is defined and a == "x" }}`,
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Is)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("defined")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.And)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Equals)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.String), "Val": Equal("x")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a logical 'or' expression",
				`{{ a is defined or a == "x" }}`,
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Is)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("defined")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Or)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Equals)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.String), "Val": Equal("x")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"is a logical 'not' expression",
				`{{ not a }}`,
				[]Fields{
					{"Type": Equal(tokens.VariableBegin)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Not)},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.Name), "Val": Equal("a")},
					{"Type": Equal(tokens.Whitespace)},
					{"Type": Equal(tokens.VariableEnd)},
					{"Type": Equal(tokens.EOF)},
				},
			},
			{
				"when when the input is multiline and exact positions matter",
				heredoc.Doc(`
					Hello
					{#
					    Multiline comment
					#}
					World
				`),
				[]Fields{
					{"Type": Equal(tokens.Data), "Val": Equal("Hello\n"), "Pos": Equal(0), "Line": Equal(1), "Col": Equal(1)},
					{"Type": Equal(tokens.CommentBegin), "Val": Equal("{#"), "Pos": Equal(6), "Line": Equal(2), "Col": Equal(1)},
					{"Type": Equal(tokens.Data), "Val": Equal("\n    Multiline comment\n"), "Pos": Equal(8), "Line": Equal(2), "Col": Equal(3)},
					{"Type": Equal(tokens.CommentEnd), "Val": Equal("#}"), "Pos": Equal(31), "Line": Equal(4), "Col": Equal(1)},
					{"Type": Equal(tokens.Data), "Val": Equal("\nWorld"), "Pos": Equal(33), "Line": Equal(4), "Col": Equal(3)},
					{"Type": Equal(tokens.EOF), "Val": Equal(""), "Pos": Equal(39), "Line": Equal(5), "Col": Equal(6)},
				},
			},
		} {
			t := testCase
			Context(fmt.Sprintf("when the input %s", t.description), func() {
				BeforeEach(func() {
					*lexer = tokens.NewLexer(t.input, config.New())
				})
				elements := make(Elements)
				for index, fields := range t.tokenMatcher {
					elements[strconv.Itoa(index)] = PointTo(MatchFields(IgnoreExtras, fields))
				}
				It("should return the expected tokens", func() {
					Expect(*returnedTokens).To(MatchAllElementsWithIndex(
						func(index int, _ any) string { return strconv.Itoa(index) },
						elements,
					))
				})
			})
		}
	})
	Context("when overriding the default delimiters", func() {
		BeforeEach(func() {
			config := config.New()

			config.BlockStartString = "<@"
			config.BlockEndString = "@>"
			config.VariableStartString = "{$"
			config.VariableEndString = "$}"
			config.CommentStartString = "(###"
			config.CommentEndString = "###)"

			*lexer = tokens.NewLexer(`<@ block @>{$ variable $}(### comment ###)`, config)
		})
		It("should return the expected tokens", func() {
			Expect(*returnedTokens).To(MatchAllElementsWithIndex(
				func(index int, _ any) string { return strconv.Itoa(index) },
				Elements{
					"0": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.BlockBegin),
					})),
					"1": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.Whitespace),
					})),
					"2": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.Name),
						"Val":  Equal("block"),
					})),
					"3": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.Whitespace),
					})),
					"4": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.BlockEnd),
					})),
					"5": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.VariableBegin),
					})),
					"6": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.Whitespace),
					})),
					"7": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.Name),
						"Val":  Equal("variable"),
					})),
					"8": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.Whitespace),
					})),
					"9": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.VariableEnd),
					})),
					"10": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.CommentBegin),
					})),
					"11": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.Data),
						"Val":  Equal(" comment "),
					})),
					"12": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.CommentEnd),
					})),
					"13": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.EOF),
					})),
				},
			))
		})
	})
	Context("when trimming comments aggressively", func() {
		BeforeEach(func() {
			*lexer = tokens.NewLexer("{#-#}", config.New())
		})

		It("should not panic and should return the expected tokens", func() {
			Expect(*returnedTokens).To(MatchAllElementsWithIndex(
				func(index int, _ any) string { return strconv.Itoa(index) },
				Elements{
					"0": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.CommentBegin),
						"Val":  Equal("{#-"),
					})),
					"1": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.CommentEnd),
						"Val":  Equal("#}"),
					})),
					"2": PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(tokens.EOF),
					})),
				},
			))
		})
	})
})

var _ = Context("raw block end patterns", func() {
	It("should share compiled patterns between lexers with the same block delimiter", func() {
		first := tokens.NewLexer("", config.New())
		second := tokens.NewLexer("", config.New())
		Expect(first.RawControlStructures["raw"]).To(BeIdenticalTo(second.RawControlStructures["raw"]))
		Expect(first.RawControlStructures["comment"]).To(BeIdenticalTo(second.RawControlStructures["comment"]))
	})
	It("should keep each lexer's raw control structures independent", func() {
		first := tokens.NewLexer("", config.New())
		second := tokens.NewLexer("", config.New())
		delete(first.RawControlStructures, "raw")
		Expect(second.RawControlStructures).To(HaveKey("raw"))
	})
	It("should compile separate patterns for each block delimiter", func() {
		custom := config.New()
		custom.BlockStartString = "<@"
		Expect(tokens.NewLexer("", custom).RawControlStructures["raw"].String()).To(Equal(`<@[-+]?\s*endraw`))
		Expect(tokens.NewLexer("", config.New()).RawControlStructures["raw"].String()).To(Equal(`\{%[-+]?\s*endraw`))
	})
	It("should keep raw block contents as data with overridden delimiters", func() {
		custom := config.New()
		custom.BlockStartString = "<@"
		custom.BlockEndString = "@>"
		custom.VariableStartString = "{$"
		custom.VariableEndString = "$}"
		stream := tokens.Lex(`<@ raw @>{$ kept $}<@ endraw @>`, custom)
		var data []string
		for !stream.End() {
			if token := stream.Next(); token.Type == tokens.Data {
				data = append(data, token.Val)
			}
		}
		Expect(data).To(Equal([]string{"{$ kept $}"}))
	})
})
