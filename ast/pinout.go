package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

// Pinout is an AST node that represents a mapping of device terminals to
// pads within a land.
type Pinout struct {
	WithRange
	Name   string
	Device Node // Expression that evaluates to a device; nil when implied by context
	Land   Node // Expression that evaluates to a land
	Body   *StatementBlock

	HeaderRange source.Range
}

func (n *Pinout) walkChildNodes(cb internalWalkFunc) {
	if n.Device != nil {
		cb(n.Device)
	}
	if n.Land != nil {
		cb(n.Land)
	}
	cb(n.Body)
}

func (n *Pinout) DeclRange() source.Range {
	return n.HeaderRange
}
