package ast

type Call struct {
	WithRange
	Callee Node
	Args   *Arguments
}

func (n *Call) walkChildNodes(cb internalWalkFunc) {
	cb(n.Callee)
	cb(n.Args)
}

type Arguments struct {
	WithRange
	Positional []Node
	Named      []*NamedArgument
}

func (n *Arguments) walkChildNodes(cb internalWalkFunc) {
	for _, c := range n.Positional {
		cb(c)
	}
	for _, c := range n.Named {
		cb(c)
	}
}

type NamedArgument struct {
	WithRange
	Name  string
	Value Node
}

func (n *NamedArgument) walkChildNodes(cb internalWalkFunc) {
	cb(n.Value)
}
