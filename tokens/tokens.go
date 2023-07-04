package tokens

import "fmt"

// TokenType identifies the type of a token
type Type int

// Known tokens
const (
	Error Type = iota
	Add
	Assign
	Colon
	Comma
	Div
	Dot
	Eq
	Floordiv
	Gt
	Gteq
	Lbrace
	Lbracket
	Lparen
	Lt
	Lteq
	Not
	Is
	In
	And
	Or
	Mod
	Mul
	Ne
	Pipe
	Pow
	Rbrace
	Rbracket
	Rparen
	Semicolon
	Sub
	Tilde
	Whitespace
	Float
	Integer
	Name
	String
	Operator
	BlockBegin
	BlockEnd
	VariableBegin
	VariableEnd
	RawBegin
	RawEnd
	CommentBegin
	CommentEnd
	Comment
	LinestatementBegin
	LinestatementEnd
	LinecommentBegin
	LinecommentEnd
	Linecomment
	Data
	Initial
	EOF
)

// TokenNames maps token types to their human readable name
var Names = map[Type]string{
	Error:              "Error",
	Add:                "Add",
	Assign:             "Assign",
	Colon:              "Colon",
	Comma:              "Comma",
	Div:                "Div",
	Dot:                "Dot",
	Eq:                 "Eq",
	Not:                "Not",
	Is:                 "Is",
	In:                 "In",
	Floordiv:           "Floordiv",
	Gt:                 "Gt",
	Gteq:               "Gteq",
	Lbrace:             "Lbrace",
	Lbracket:           "Lbracket",
	Lparen:             "Lparen",
	Lt:                 "Lt",
	Lteq:               "Lteq",
	And:                "And",
	Or:                 "Or",
	Mod:                "Mod",
	Mul:                "Mul",
	Ne:                 "Ne",
	Pipe:               "Pipe",
	Pow:                "Pow",
	Rbrace:             "Rbrace",
	Rbracket:           "Rbracket",
	Rparen:             "Rparen",
	Semicolon:          "Semicolon",
	Sub:                "Sub",
	Tilde:              "Tilde",
	Whitespace:         "Whitespace",
	Float:              "Float",
	Integer:            "Integer",
	Name:               "Name",
	String:             "String",
	Operator:           "Operator",
	BlockBegin:         "BlockBegin",
	BlockEnd:           "BlockEnd",
	VariableBegin:      "VariableBegin",
	VariableEnd:        "VariableEnd",
	RawBegin:           "RawBegin",
	RawEnd:             "RawEnd",
	CommentBegin:       "CommentBegin",
	CommentEnd:         "CommentEnd",
	Comment:            "Comment",
	LinestatementBegin: "LinestatementBegin",
	LinestatementEnd:   "LinestatementEnd",
	LinecommentBegin:   "LinecommentBegin",
	LinecommentEnd:     "LinecommentEnd",
	Linecomment:        "Linecomment",
	Data:               "Data",
	Initial:            "Initial",
	EOF:                "EOF",
}

// Token represents a unit of lexing
type Token struct {
	Type Type
	Val  string
	Pos  int
	Line int
	Col  int
	Trim bool
}

func (t Token) String() string {
	val := t.Val
	if len(val) > 1000 {
		val = fmt.Sprintf("%s...%s", val[:10], val[len(val)-5:])
	}

	return fmt.Sprintf("<Token[%s] Val='%s' Pos=%d Line=%d Col=%d>",
		Names[t.Type], val, t.Pos, t.Line, t.Col)
}
