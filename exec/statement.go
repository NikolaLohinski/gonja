package exec

import (
	"github.com/nikolalohinski/gonja/nodes"
)

type Statement interface {
	nodes.Statement
	Execute(*Renderer, *nodes.StatementBlock) error
}
