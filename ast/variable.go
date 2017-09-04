package ast

type Variable struct {
	WithRange
	Name string
}

func (n *Variable) walkChildNodes(cb internalWalkFunc) {
	// Variable is a leaf node, so there are no child nodes to walk
}
