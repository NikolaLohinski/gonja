package statements

import (
	"fmt"

	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	"github.com/nikolalohinski/gonja/tokens"
)

type ExtendsStmt struct {
	location    *tokens.Token
	filename    string
	withContext bool
}

func (stmt *ExtendsStmt) Position() *tokens.Token {
	return stmt.location
}

func (stmt *ExtendsStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("ExtendsStmt(Filename=%s Line=%d Col=%d)", stmt.filename, t.Line, t.Col)
}

func (node *ExtendsStmt) Execute(r *exec.Renderer) error {
	return nil
}

func extendsParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &ExtendsStmt{
		location: p.Current(),
	}

	if p.Template.Parent != nil {
		return nil, args.Error("this template has already one parent", args.Current())
	}

	// var filename nodes.Node
	if filename := args.Match(tokens.String); filename != nil {
		stmt.filename = filename.Val

		extended, err := p.Extend(stmt.filename)
		if err != nil {
			return nil, fmt.Errorf("unable to load template '%s': %s", filename, err)
		}

		p.Template.Parent = extended
	} else {
		return nil, args.Error("tag 'extends' requires a template filename as string", args.Current())
	}

	if tok := args.MatchName("with", "without"); tok != nil {
		if args.MatchName("context") != nil {
			stmt.withContext = tok.Val == "with"
		} else {
			args.Stream().Backup()
		}
	}

	if !args.End() {
		return nil, args.Error("tag 'extends' only takes 1 argument", nil)
	}

	return stmt, nil
}

func init() {
	All.Register("extends", extendsParser)
}
