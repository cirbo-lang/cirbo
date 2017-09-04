package ast

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Import struct {
	WithRange
	Package string
	Name    string

	PackageRange source.Range
}

func (n *Import) walkChildNodes(cb internalWalkFunc) {
	// Import is a leaf node, so there are no child nodes to walk
}
