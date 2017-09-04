package ast

type Object struct {
	WithRange
	Elements []*ObjectElem
}

func (n *Object) walkChildNodes(cb internalWalkFunc) {
	for _, c := range n.Elements {
		cb(c)
	}
}

type ObjectElem struct {
	WithRange
	Name  string
	Value Node
}

func (n *ObjectElem) walkChildNodes(cb internalWalkFunc) {
	cb(n.Value)
}
