package ast

import (
	"math/big"

	"github.com/cirbo-lang/cirbo/units"
)

type NumberLit struct {
	WithRange
	Value *big.Float
	Unit  string // empty for dimensionless values
}

func (n *NumberLit) walkChildNodes(cb internalWalkFunc) {
	// NumberLit is a leaf node, so there are no child nodes to walk
}

// Quantity returns a units.Quantity corresponding to the value in the
// receiver.
//
// The Unit string must be a valid unit name, meaning that if passed to
// IsQuantityUnitKeyword the result would be true. Otherwise, this method
// will panic.
//
// If the unit string is empty, the result is a dimensionless quantity
// which can then serve as a plain number for calculations.
func (n *NumberLit) Quantity() units.Quantity {
	if n.Unit == "" {
		return units.MakeDimensionless(n.Value)
	}

	unit := units.ByName(n.Unit)
	if unit == nil {
		panic("attempt to call Quantity on NumberLit with invalid Unit")
	}

	return units.MakeQuantity(n.Value, unit)
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

func IsQuantityUnitKeyword(kw string) bool {
	return units.ByName(kw) != nil
}
