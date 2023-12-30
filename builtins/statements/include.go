package statements

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type IncludeStmt struct {
	location           *tokens.Token
	filenameExpression nodes.Expression
	template           *nodes.Template
	ignoreMissing      bool
	withContext        bool
	isEmpty            bool
}

func (stmt *IncludeStmt) Position() *tokens.Token {
	return stmt.location
}

func (stmt *IncludeStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("IncludeStmt(Filename=%s Line=%d Col=%d)", stmt.filenameExpression, t.Line, t.Col)
}

func (stmt *IncludeStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	if stmt.isEmpty {
		return nil
	}
	sub := r.Inherit()

	filenameValue := r.Eval(stmt.filenameExpression)
	if filenameValue.IsError() {
		return errors.Wrap(filenameValue, `Unable to evaluate filename`)
	}

	filename := filenameValue.String()
	loader, err := r.Loader.Inherit(filename)
	if err != nil {
		return errors.Errorf("failed to inherit loader: %s", err)
	}

	included, err := exec.NewTemplate(filename, r.Config, loader, r.Environment)
	if err != nil {
		if stmt.ignoreMissing {
			return nil
		} else {
			return fmt.Errorf("Unable to load template '%s': %s", filename, err)
		}
	}
	sub = exec.NewRenderer(r.Environment, r.Output, r.Config.Inherit(), loader, included)

	return sub.Execute()
}

func includeParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &IncludeStmt{
		location: p.Current(),
	}

	filenameExpression, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	stmt.filenameExpression = filenameExpression

	if args.MatchName("ignore") != nil {
		if args.MatchName("missing") != nil {
			stmt.ignoreMissing = true
		} else {
			args.Stream().Backup()
		}
	}

	if tok := args.MatchName("with", "without"); tok != nil {
		if args.MatchName("context") != nil {
			stmt.withContext = tok.Val == "with"
		} else {
			args.Stream().Backup()
		}
	}

	if !args.End() {
		return nil, args.Error("Malformed 'include'-tag args.", nil)
	}

	return stmt, nil
}

func init() {
	All.Register("include", includeParser)
}
