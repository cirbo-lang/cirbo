package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Board struct {
	WithRange
	Name string
	Body *StatementBlock

	HeaderRange source.Range
}

func (n *Board) walkChildNodes(cb internalWalkFunc) {
	cb(n.Body)
}

func (n *Board) DeclRange() source.Range {
	return n.HeaderRange
}
