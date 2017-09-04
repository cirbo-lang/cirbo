package ast

type Export struct {
	WithRange
	Expr Node
}

func (n *Export) walkChildNodes(cb internalWalkFunc) {
	cb(n.Expr)
}
