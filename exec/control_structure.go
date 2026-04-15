package exec

import (
	"github.com/ardanlabs/gonja/nodes"
)

type ControlStructure interface {
	nodes.ControlStructure
	Execute(*Renderer, *nodes.ControlStructureBlock) error
}
