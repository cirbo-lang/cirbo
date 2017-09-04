package ast

type Slice struct {
	WithRange
	Source Node
	Start  Node
	End    Node
}

func (n *Slice) walkChildNodes(cb internalWalkFunc) {
	cb(n.Source)
	cb(n.Start)
	cb(n.End)
}
