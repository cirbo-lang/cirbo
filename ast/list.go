package ast

type List struct {
	WithRange
	Elements []Node
}

func (n *List) walkChildNodes(cb internalWalkFunc) {
	for _, c := range n.Elements {
		cb(c)
	}
}
