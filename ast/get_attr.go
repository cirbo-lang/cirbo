package ast

type GetAttr struct {
	WithRange
	Source Node
	Name   string
}

func (n *GetAttr) walkChildNodes(cb internalWalkFunc) {
	cb(n.Source)
}
