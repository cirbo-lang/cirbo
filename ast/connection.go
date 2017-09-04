package ast

// Connection is an AST node that represents a sequence of connections
// between nets, terminals, devices, and buses.
type Connection struct {
	WithRange
	Seq []Node
}

func (n *Connection) walkChildNodes(cb internalWalkFunc) {
	for _, c := range n.Seq {
		cb(c)
	}
}

// NoConnection is an AST node that represents a declaration that a particular
// terminal is intentionally unconnected.
type NoConnection struct {
	WithRange
	Terminal Node
}

func (n *NoConnection) walkChildNodes(cb internalWalkFunc) {
	cb(n.Terminal)
}
