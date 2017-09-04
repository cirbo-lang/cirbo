package ast

type GetIndex struct {
	WithRange
	Source Node
	Index  Node
}

func (n *GetIndex) walkChildNodes(cb internalWalkFunc) {
	cb(n.Source)
	cb(n.Index)
}
