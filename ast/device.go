package ast

type Device struct {
	WithRange
	Name string
}

func (n *Device) walkChildNodes(cb internalWalkFunc) {
	// TODO: Implement child nodes
	panic("walkChildNodes not implemented for Device")
}
