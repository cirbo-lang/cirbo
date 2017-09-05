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

func IsQuantityUnitKeyword(kw string) bool {
	// TODO: this will get more sophisticated once we have actual
	// unit handling in the evaluator.
	switch kw {

	case "m", "mm", "cm", "nm", "um", "in", "mil", "ft":
		return true

	case "g", "kg", "lb", "oz":
		return true

	case "A", "mA", "uA", "kA":
		return true

	case "V", "mV", "uV", "kV":
		return true

	case "W", "mW", "uW", "kW", "MW", "GW":
		return true

	case "ohm", "mohm", "uohm", "kohm", "Mohm":
		return true

	case "farad", "ufarad", "mfarad", "Mfarad":
		return true

	case "H", "uH", "mH", "kH", "MH":
		return true

	case "s", "ms", "us":
		return true

	case "Hz", "kHz", "MHz":
		return true

	case "K", "degC", "degF":
		return true

	default:
		return false

	}
}
