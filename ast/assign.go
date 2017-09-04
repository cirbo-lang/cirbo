package ast

type Assign struct {
	WithRange
	Name  string
	Value Node
}

func (n *Assign) walkChildNodes(cb internalWalkFunc) {
	cb(n.Value)
}
