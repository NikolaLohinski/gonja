package exec

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type Template struct {
	source      string
	config      *config.Config
	environment *Environment
	loader      loaders.Loader
	tokens      *tokens.Stream
	parser      *parser.Parser
	root        *nodes.Template
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
		source:      source.String(),
		config:      config,
		loader:      loader,
		tokens:      tokens.Lex(source.String()),
		environment: environment,
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
	var output strings.Builder

	renderer := NewRenderer(&Environment{
		Tests:      t.environment.Tests,
		Filters:    t.environment.Filters,
		Statements: t.environment.Statements,
		Context:    t.environment.Context.Inherit().Update(ctx),
	}, &output, t.config, t.loader, t)

	err := renderer.Execute()
	if err != nil {
		return "", errors.Wrap(err, "unable to execute template")
	}

	return output.String(), nil
}

func (t *Template) Macros() map[string]*nodes.Macro {
	return t.root.Macros
}
