package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Node interface {
	SourceRange() source.Range

	walkChildNodes(cb internalWalkFunc)
}

type internalWalkFunc func(Node)

// WithRange is a mixin included in all nodes that provides the SourceRange
// method.
type WithRange struct {
	source.Range
}

func (n *WithRange) SourceRange() source.Range {
	return n.Range
}
