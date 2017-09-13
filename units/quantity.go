package units

import (
	"fmt"
	"math/big"
)

// Quantity is the combination of a value and a unit.
type Quantity struct {
	value *big.Float
	unit  *Unit
}

// MakeQuantity initializes a Quantity for a given value and unit.
func MakeQuantity(value *big.Float, unit *Unit) Quantity {
	if unit == nil {
		unit = dimless
	}
	return Quantity{
		value: value,
		unit:  unit,
	}
}

// Value returns the value of the receiving quantity.
func (q Quantity) Value() *big.Float {
	// Since big floats are mutable, we return a copy to prevent
	// the caller from altering our internal state.
	return (&big.Float{}).Copy(q.value)
}

// Value returns the unit of the receiving quantity.
func (q Quantity) Unit() *Unit {
	return q.unit
}

// String returns a compact, human-readable representation of the receiver.
//
// It is primarily intended for debugging and is thus not optimized.
func (q Quantity) String() string {
	valStr := q.value.String()
	unitStr := q.unit.String()

	if unitStr == "" {
		return valStr
	}

	return fmt.Sprintf("%s %s", valStr, unitStr)
}
