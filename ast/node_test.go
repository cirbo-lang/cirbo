package ast

import (
	"fmt"
	"testing"
)

func TestNodeImpls(t *testing.T) {
	// Everything here actually gets checked at compile time, so this
	// test function is just here for visibility in verbose test output.
	var tests []Node = []Node{
		&Arguments{},
		&ArithmeticBinary{},
		&ArithmeticUnary{},
		&Assign{},
		&Attr{},
		&BooleanLit{},
		&Call{},
		&Circuit{},
		&Connection{},
		&Designator{},
		&Device{},
		&Export{},
		&File{},
		&GetAttr{},
		&GetIndex{},
		&Import{},
		&Land{},
		&List{},
		&NamedArgument{},
		&NetCons{},
		&NumberLit{},
		&Object{},
		&ObjectElem{},
		&Pinout{},
		&Slice{},
		&StatementBlock{},
		&StringLit{},
		&Terminal{},
		&Variable{},
	}

	for _, n := range tests {
		t.Run(fmt.Sprintf("%T", n)[5:], func(t *testing.T) {
			// nothing to do; if we compiled then we're good!
		})
	}
}
