package ast

// Attr is an AST node that represents the declaration of an attribute of
// some parameterized object.
type Attr struct {
	WithRange
	Name  string
	Type  Node
	Value Node
}

func (n *Attr) walkChildNodes(cb internalWalkFunc) {
	if n.Type != nil {
		cb(n.Type)
	}
	if n.Value != nil {
		cb(n.Value)
	}
}
