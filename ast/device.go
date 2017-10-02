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
	cb(n.Params)
	cb(n.Body)
}

func (n *Device) DeclRange() source.Range {
	return n.HeaderRange
}
