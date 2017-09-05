package ast

type ParenExpr struct {
	WithRange
	Content Node
}

func (n *ParenExpr) walkChildNodes(cb internalWalkFunc) {
	cb(n.Content)
}
