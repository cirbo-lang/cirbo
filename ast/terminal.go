package ast

import (
	"github.com/cirbo-lang/cirbo/cbo"
)

type Terminal struct {
	WithRange
	Name       string
	Type       cbo.ERCType
	Dir        cbo.ERCDir
	Role       cbo.TerminalRole
	OutputType cbo.ERCOutputType
}

func (n *Terminal) walkChildNodes(cb internalWalkFunc) {
	// Terminal is a leaf node
}
