package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Circuit struct {
	WithRange
	Name string
	Body []Node

	HeaderRange source.Range
}

func (n *Circuit) walkChildNodes(cb internalWalkFunc) {
	// TODO: Implement child nodes
	panic("walkChildNodes not implemented for Circuit")
}

func (n *Circuit) DeclRange() source.Range {
	return n.HeaderRange
}
