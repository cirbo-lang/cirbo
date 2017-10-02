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
	cb(n.Params)
	cb(n.Body)
}

func (n *Land) DeclRange() source.Range {
	return n.HeaderRange
}
