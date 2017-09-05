package ast

// Invalid is a funny sort of AST node that represents a portion of the tree
// that can't be completed because of a parse error.
//
// It is never used in a valid AST, but may be included if the tree was
// returned with error diagnostics. The intent is to still make the AST
// have a sound structure so that it can still be used for certain limited
// sorts of static analysis in the presence of errors.
type Invalid struct {
	WithRange
}

func (n *Invalid) walkChildNodes(cb internalWalkFunc) {
	// Invalid is a leaf node
}
