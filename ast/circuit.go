package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Circuit struct {
	WithRange
	Name   string
	Params *Arguments
	Body   *StatementBlock

	HeaderRange source.Range
}

func (n *Circuit) walkChildNodes(cb internalWalkFunc) {
	cb(n.Params)
	cb(n.Body)
}

func (n *Circuit) DeclRange() source.Range {
	return n.HeaderRange
}
