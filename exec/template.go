package exec

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/config"
	"github.com/nikolalohinski/gonja/loaders"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	"github.com/nikolalohinski/gonja/tokens"
)

type Template struct {
	source      string
	config      *config.Config
	environment *Environment
	tokens      *tokens.Stream
	parser      *parser.Parser
	root        *nodes.Template
	macros      MacroSet
}

func NewTemplate(identifier string, config *config.Config, loader loaders.Loader, environment *Environment) (*Template, error) {
	input, err := loader.Read(identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to reader template '%s': %s", identifier, err)
	}

	source := new(strings.Builder)
	if _, err := io.Copy(source, input); err != nil {
		return nil, fmt.Errorf("failed to copy '%s' to string buffer: %s", source, err)
	}

	t := &Template{
		source: source.String(),
		config: config,
		tokens: tokens.Lex(source.String()),
	}

	t.parser = parser.NewParser(identifier, t.tokens, config, loader, environment.Statements)

	root, err := t.parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse template '%s': %s", source, err)
	}
	t.root = root

	return t, nil
}

// Executes the template and returns the rendered template as a string
func (t *Template) Execute(ctx *Context) (string, error) {
	var b strings.Builder

	renderingContext := t.environment.Context.Inherit()
	renderingContext.Update(ctx)

	var builder strings.Builder
	renderer := NewRenderer(t.environment, &builder, t.config, t)

	err := renderer.Execute()
	if err != nil {
		return "", errors.Wrap(err, "unable to execute template")
	}

	b.WriteString(renderer.String())

	return b.String(), nil
}
