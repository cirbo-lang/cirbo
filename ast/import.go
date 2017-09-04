package ast

type Import struct {
	WithRange
	Package string
}

func (n *Import) walkChildNodes(cb internalWalkFunc) {
	// Import is a leaf node, so there are no child nodes to walk
}
