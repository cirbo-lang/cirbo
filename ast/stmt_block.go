package ast

type StatementBlock struct {
	WithRange

	Statements []Node
}

func (n *StatementBlock) walkChildNodes(cb internalWalkFunc) {
	for _, cn := range n.Statements {
		cb(cn)
	}
}
