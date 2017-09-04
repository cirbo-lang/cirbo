package ast

// Pinout is an AST node that represents a mapping of device terminals to
// pads within a land.
type Pinout struct {
	WithRange
	Name        string
	Device      Node   // Expression that evaluates to a device; nil when implied by context
	Land        Node   // Expression that evaluates to a land
	Connections []Node // Connection nodes that define relatinships between terminals and pads
}

func (n *Pinout) walkChildNodes(cb internalWalkFunc) {
	if n.Device != nil {
		cb(n.Device)
	}
	if n.Land != nil {
		cb(n.Land)
	}
	for _, c := range n.Connections {
		cb(c)
	}
}
