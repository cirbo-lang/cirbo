package ast

import (
	"math/big"
)

type NumberLit struct {
	WithRange
	Value *big.Float
}

func (n *NumberLit) walkChildNodes(cb internalWalkFunc) {
	// NumberLit is a leaf node, so there are no child nodes to walk
}

type StringLit struct {
	WithRange
	Value string
}

func (n *StringLit) walkChildNodes(cb internalWalkFunc) {
	// StringLit is a leaf node, so there are no child nodes to walk
}

type BooleanLit struct {
	WithRange
	Value bool
}

func (n *BooleanLit) walkChildNodes(cb internalWalkFunc) {
	// BooleanLit is a leaf node, so there are no child nodes to walk
}

type QuantityLit struct {
	WithRange
	Value *big.Float
	Unit  string
}

func (n *QuantityLit) walkChildNodes(cb internalWalkFunc) {
	// QuantityLit is a leaf node, so there are no child nodes to walk
}
