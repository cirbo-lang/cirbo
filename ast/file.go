package ast

import (
	"github.com/cirbo-lang/cirbo/projpath"
)

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

func (n *File) Filename() projpath.FilePath {
	return n.WithRange.Filename
}
