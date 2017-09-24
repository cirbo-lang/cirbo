package ast

import (
	"github.com/cirbo-lang/cirbo/cbo"
)

type Terminal struct {
	WithRange
	Name       string
	Type       cbo.TerminalType
	Dir        cbo.TerminalDir
	Role       cbo.TerminalRole
	OutputType cbo.TerminalOutputType
}

func (n *Terminal) walkChildNodes(cb internalWalkFunc) {
	// Terminal is a leaf node
}
