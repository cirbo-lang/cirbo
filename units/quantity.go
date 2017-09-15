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

// CommensurableWith returns true if the receiver and the other given quantity
// have commensurable units.
//
// It is functionally equivalent to q.Unit().CommensurableWith(other.Unit()) .
func (q Quantity) CommensurableWith(other Quantity) bool {
	return q.unit.CommensurableWith(other.unit)
}

// ConvertableTo returns true if the receiver's unit is commensurable with the
// given unit, which indicates that the receiver could be converted to the
// given unit.
//
// It is functionally equivalent to q.Unit().CommensurableWith(new) .
func (q Quantity) ConvertableTo(new *Unit) bool {
	return q.unit.CommensurableWith(new)
}

// Convert returns a new Quantity that is equivalent to the receiver but
// is expressed in the given unit.
//
// Will panic if the receiver's unit is not commensurable with the
// given unit. Use ConvertableTo before calling if unsure.
func (q Quantity) Convert(new *Unit) Quantity {
	nf := q.Value() // creates a copy, so we can mutate
	old := q.Unit()

	if new == old {
		// NOTE: the above is a pointer comparison, so it only works for
		// named units where we guarantee singleton values and thus
		// comparable pointers. For constructed units, we'll fall through
		// to the machinery below.
		return q
	}

	// Eliminate the scale before we begin, so we only have to worry
	// about the base units.
	switch {
	case old.scale > 0:
		fs := (&big.Float{}).SetInt64(old.scale)
		nf.Mul(nf, fs)
	case old.scale < 0:
		fs := (&big.Float{}).SetInt64(-old.scale)
		nf.Quo(nf, fs)
	}

	// Now for each base dimension we'll convert to the primary unit and
	// then to the target, unless units already match.
	if old.base.Mass != new.base.Mass {
		nf.Quo(nf, powerScale(&old.base.Mass.Scale, old.dim.Mass))
		nf.Mul(nf, powerScale(&new.base.Mass.Scale, new.dim.Mass))
	}
	if old.base.Length != new.base.Length {
		nf.Quo(nf, powerScale(&old.base.Length.Scale, old.dim.Length))
		nf.Mul(nf, powerScale(&new.base.Length.Scale, new.dim.Length))
	}
	if old.base.Angle != new.base.Angle {
		nf.Quo(nf, &old.base.Angle.Scale)
		nf.Mul(nf, &new.base.Angle.Scale)
	}
	if old.base.Time != new.base.Time {
		nf.Quo(nf, &old.base.Time.Scale)
		nf.Mul(nf, &new.base.Time.Scale)
	}
	if old.base.ElectricCurrent != new.base.ElectricCurrent {
		nf.Quo(nf, &old.base.ElectricCurrent.Scale)
		nf.Mul(nf, &new.base.ElectricCurrent.Scale)
	}
	if old.base.LuminousIntensity != new.base.LuminousIntensity {
		nf.Quo(nf, &old.base.LuminousIntensity.Scale)
		nf.Mul(nf, &new.base.LuminousIntensity.Scale)
	}

	// Finally, apply any scale required by the new unit.
	switch {
	case new.scale > 0:
		fs := (&big.Float{}).SetInt64(new.scale)
		nf.Quo(nf, fs)
	case new.scale < 0:
		fs := (&big.Float{}).SetInt64(-new.scale)
		nf.Mul(nf, fs)
	}

	return MakeQuantity(nf, new)
}

func (q Quantity) WithStandardUnits() Quantity {
	u := q.unit.ToStandardUnits()
	return q.Convert(u)
}

// Multiply computes the product of the receiver and the given quantity,
// multiplying both the value and the units to produce a new quantity.
//
// If the two values use the same units for the base dimensions, or have
// non-overlapping base dimensions, these base units will be left untouched.
// If they differ, both quantities will be converted to standard units so that
// the result has a consistent set of units.
//
// For example, inches divided by seconds yields a result in inches per second,
// but inches multiplied by yards yields a result in square meters since the
// differing length units must be normalized to avoid producing the nonsense
// unit "inch-yards".
func (q Quantity) Multiply(o Quantity) Quantity {
	if !q.unit.SameBaseUnits(o.unit) {
		q = q.WithStandardUnits()
		o = o.WithStandardUnits()
	}

	nu := q.unit.Multiply(o.unit)
	nv := (&big.Float{}).Mul(q.value, o.value)

	return Quantity{
		unit:  nu,
		value: nv,
	}
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
