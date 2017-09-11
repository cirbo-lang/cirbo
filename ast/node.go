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

// NodeDecl is an extension of Node that adds a DeclRange method.
type NodeDecl interface {
	Node

	// DeclRange returns the range that covers the part of a definition
	// node that serves as a declaration.
	//
	// For complex constructs that include nested statement blocks, this
	// returns just the range around the tokens that introduce the block,
	// and not the block itself. This result is, therefore, a better
	// range to use when presenting diagnostics relating to declarations
	// so as to focus the user's attention on the declaration portion.
	DeclRange() source.Range
}
