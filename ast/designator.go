package ast

// Designator is an AST node that represents the declaration of the designator
// for a device.
type Designator struct {
	WithRange
	Value Node
}

func (n *Designator) walkChildNodes(cb internalWalkFunc) {
	cb(n.Value)
}
