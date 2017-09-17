package ast

type Export struct {
	WithRange
	Value Node
	Name  string
}

func (n *Export) walkChildNodes(cb internalWalkFunc) {
	cb(n.Value)
}
