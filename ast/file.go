package ast

type File struct {
	WithRange
	Source   []byte
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
