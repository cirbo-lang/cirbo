package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Device struct {
	WithRange

	Name   string
	Params *Arguments
	Body   *StatementBlock

	HeaderRange source.Range
}

func (n *Device) walkChildNodes(cb internalWalkFunc) {
	// TODO: Implement child nodes
	panic("walkChildNodes not implemented for Device")
}
