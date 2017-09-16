package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Land struct {
	WithRange
	Name   string
	Params *Arguments
	Body   *StatementBlock

	HeaderRange source.Range
}

func (n *Land) walkChildNodes(cb internalWalkFunc) {
	// TODO: Implement child nodes
	panic("walkChildNodes not implemented for Land")
}
