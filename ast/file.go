package ast

type File struct {
	WithRange
	TopLevel []Node
}

func (n *File) walkChildNodes(cb internalWalkFunc) {
	for _, c := range n.TopLevel {
		cb(c)
	}
}

func (n *File) Filename() string {
	return n.WithRange.Filename
}
