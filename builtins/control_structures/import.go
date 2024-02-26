package controlStructures

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type ImportControlStructure struct {
	location           *tokens.Token
	filenameExpression nodes.Expression
	as                 string
	withContext        bool
}

func (controlStructure *ImportControlStructure) Position() *tokens.Token {
	return controlStructure.location
}

func (controlStructure *ImportControlStructure) String() string {
	t := controlStructure.Position()
	return fmt.Sprintf("ImportControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (controlStructure *ImportControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {

	filenameValue := r.Eval(controlStructure.filenameExpression)
	if filenameValue.IsError() {
		return errors.Wrap(filenameValue, `Unable to evaluate filename`)
	}

	filename, err := r.Loader.Resolve(filenameValue.String())
	if err != nil {
		return errors.Errorf("failed to resolve filename: %s", err)
	}

	loader, err := r.Loader.Inherit(filename)
	if err != nil {
		return fmt.Errorf("failed to inherit loader from '%s': %s", filename, r.Loader)
	}

	template, err := exec.NewTemplate(filename, r.Config, loader, r.Environment)
	if err != nil {
		return fmt.Errorf("unable to load template '%s': %s", filename, err)
	}

	macros := map[string]exec.Macro{}
	for name, macro := range template.Macros() {
		fn, err := exec.MacroNodeToFunc(macro, r)
		if err != nil {
			return errors.Wrapf(err, `Unable to import macro '%s'`, name)
		}
		macros[name] = fn
	}
	r.Environment.Context.Set(controlStructure.as, macros)

	return nil
}

type FromImportControlStructure struct {
	location           *tokens.Token
	FilenameExpression nodes.Expression
	WithContext        bool
	Template           *nodes.Template
	As                 map[string]string
	Macros             map[string]*nodes.Macro // alias/name -> macro instance
}

func (controlStructure *FromImportControlStructure) Position() *tokens.Token {
	return controlStructure.location
}

func (controlStructure *FromImportControlStructure) String() string {
	t := controlStructure.Position()
	return fmt.Sprintf("FromImportControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (controlStructure *FromImportControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {

	filenameValue := r.Eval(controlStructure.FilenameExpression)
	if filenameValue.IsError() {
		return errors.Wrap(filenameValue, `Unable to evaluate filename`)
	}

	filename, err := r.Loader.Resolve(filenameValue.String())
	if err != nil {
		return errors.Errorf("failed to resolve filename: %s", err)
	}

	loader, err := r.Loader.Inherit(filename)
	if err != nil {
		return fmt.Errorf("failed to inherit loader from '%s': %s", filename, r.Loader)
	}

	template, err := exec.NewTemplate(filename, r.Config, loader, r.Environment)
	if err != nil {
		return fmt.Errorf("unable to load template '%s': %s", filename, err)
	}

	imported := template.Macros()
	for alias, name := range controlStructure.As {
		node := imported[name]
		fn, err := exec.MacroNodeToFunc(node, r)
		if err != nil {
			return errors.Wrapf(err, `Unable to import macro '%s'`, name)
		}
		r.Environment.Context.Set(alias, fn)
	}
	return nil
}

func importParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	controlStructure := &ImportControlStructure{
		location: p.Current(),
		// Macros:   map[string]*nodes.Macro{},
	}

	if args.End() {
		return nil, args.Error("You must at least specify one macro to import.", nil)
	}

	expression, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	controlStructure.filenameExpression = expression
	if args.MatchName("as") == nil {
		return nil, args.Error(`Expected "as" keyword`, args.Current())
	}

	alias := args.Match(tokens.Name)
	if alias == nil {
		return nil, args.Error("Expected macro alias name (identifier)", args.Current())
	}
	controlStructure.as = alias.Val

	if tok := args.MatchName("with", "without"); tok != nil {
		if args.MatchName("context") != nil {
			controlStructure.withContext = tok.Val == "with"
		} else {
			args.Stream().Backup()
		}
	}
	return controlStructure, nil
}

func fromParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	controlStructure := &FromImportControlStructure{
		location: p.Current(),
		As:       map[string]string{},
	}

	if args.End() {
		return nil, args.Error("You must at least specify one macro to import.", nil)
	}

	filename, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	controlStructure.FilenameExpression = filename

	if args.MatchName("import") == nil {
		return nil, args.Error("Expected import keyword", args.Current())
	}

	for !args.End() {
		name := args.Match(tokens.Name)
		if name == nil {
			return nil, args.Error("Expected macro name (identifier).", args.Current())
		}

		// asName := macroNameToken.Val
		if args.MatchName("as") != nil {
			alias := args.Match(tokens.Name)
			if alias == nil {
				return nil, args.Error("Expected macro alias name (identifier).", nil)
			}
			// asName = aliasToken.Val
			controlStructure.As[alias.Val] = name.Val
		} else {
			controlStructure.As[name.Val] = name.Val
		}

		if tok := args.MatchName("with", "without"); tok != nil {
			if args.MatchName("context") != nil {
				controlStructure.WithContext = tok.Val == "with"
				break
			} else {
				args.Stream().Backup()
			}
		}

		if args.End() {
			break
		}

		if args.Match(tokens.Comma) == nil {
			return nil, args.Error("Expected ','.", nil)
		}
	}

	return controlStructure, nil
}
