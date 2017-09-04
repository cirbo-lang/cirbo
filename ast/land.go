package ast

type Land struct {
	WithRange
	Name string
}

func (n *Land) walkChildNodes(cb internalWalkFunc) {
	// TODO: Implement child nodes
	panic("walkChildNodes not implemented for Land")
}
