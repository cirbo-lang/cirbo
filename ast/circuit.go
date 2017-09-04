package ast

type Circuit struct {
	WithRange
	Name string
}

func (n *Circuit) walkChildNodes(cb internalWalkFunc) {
	// TODO: Implement child nodes
	panic("walkChildNodes not implemented for Circuit")
}
