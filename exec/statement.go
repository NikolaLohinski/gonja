package exec

import (
	"github.com/nikolalohinski/gonja/v2/nodes"
)

type Statement interface {
	nodes.Statement
	Execute(*Renderer, *nodes.StatementBlock) error
}
