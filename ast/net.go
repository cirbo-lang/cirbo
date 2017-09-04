package ast

type NetCons struct {
	WithRange
	Elements []Node
}

func (n *NetCons) walkChildNodes(cb internalWalkFunc) {
	for _, c := range n.Elements {
		cb(c)
	}
}
