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

// MakeDimensionless initializes a dimensionless Quantity with the given value
func MakeDimensionless(value *big.Float) Quantity {
	return Quantity{
		value: value,
		unit:  dimless,
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

// WithStandardUnits converts the quantity so it uses the standard units for
// each base dimension. The standard units are:
//
//    mass                 kg      (kilograms)
//    length               meter   (meters)
//    angle                deg     (degrees)
//    time                 s       (seconds)
//    electric current     A       (amps)
//    luminous intensity   cd      (candelas)
func (q Quantity) WithStandardUnits() Quantity {
	u := q.unit.ToStandardUnits()
	return q.Convert(u)
}

// Equal returns true if and only if the receiver is equal to the given
// quantity.
//
// Two quantities can be equal only if their units have the same
// dimensionality. If the two quantities have different units of the same
// dimensionality, they will both be converted to standard units before
// comparison.
//
// Unit conversions are done at high precision but the precision is not
// unlimited, so mismatches may occur when comparing quantities whose values
// have many significant figures.
func (q Quantity) Equal(o Quantity) bool {
	if q.unit.dim != o.unit.dim {
		return false
	}

	return q.Compare(o) == 0
}

// Same returns true if and only if the receiver has the same value _and_
// the same unit as the given quantity.
//
// This method is primarily provided for testing, to verify that some
// result has the expected unit and value.
func (q Quantity) Same(o Quantity) bool {
	if q.unit.dim != o.unit.dim {
		return false
	}
	if !q.unit.SameBaseUnits(o.unit) {
		return false
	}
	return q.Equal(o)
}

// Compare determines whether the receiver is less than, greater than or
// equal to the given quantity.
//
// Two quantities can be compared only if their units have the same
// dimensionality. This function will panic if given incommensurable
// quantities.
//
// If the two quantities have different units of the same dimensionality,
// they will both be converted to standard units before comparison.
//
// The result is as follows:
//
//     -1 if q < o
//      0 if q == 0
//      1 if q > o
func (q Quantity) Compare(o Quantity) int {
	if q.unit.dim != o.unit.dim {
		panic(fmt.Errorf("attempt to compare incommensurable quantities"))
	}

	if !q.unit.SameBaseUnits(o.unit) {
		q = q.WithStandardUnits()
		o = o.WithStandardUnits()
	}

	return q.value.Cmp(o.value)
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

// Divide computes the quotient of the receiver by the given quantity,
// dividing both the value and the units to produce a new quantity.
//
// The same normalization of unit applies as for the Multiply method.
func (q Quantity) Divide(o Quantity) Quantity {
	if !q.unit.SameBaseUnits(o.unit) {
		q = q.WithStandardUnits()
		o = o.WithStandardUnits()
	}

	nu := q.unit.Multiply(o.unit.Reciprocal())
	nv := (&big.Float{}).Quo(q.value, o.value)

	return Quantity{
		unit:  nu,
		value: nv,
	}
}

// Add computes the sum of the receiver and the given quantity, which must
// have commensurable units.
//
// If the units are not commensurable, this method will panic.
//
// If the units are identical (same base units) then the result will have the
// same units. Otherwise, the result will be in the standard units.
func (q Quantity) Add(o Quantity) Quantity {
	if !q.CommensurableWith(o) {
		panic("Attempt to Add non-commensurable quantities")
	}

	if !q.unit.SameBaseUnits(o.unit) {
		q = q.WithStandardUnits()
		o = o.WithStandardUnits()
	}

	nv := (&big.Float{}).Add(q.value, o.value)

	return Quantity{
		unit:  q.unit,
		value: nv,
	}
}

// Subtract computes the difference between the given receiver and the given
// quantity, which must both have commensurable units.
//
// If the units are not commensurable, this method will panic.
//
// If the units are identical (same base units) then the result will have the
// same units. Otherwise, the result will be in the standard units.
func (q Quantity) Subtract(o Quantity) Quantity {
	if !q.CommensurableWith(o) {
		panic("Attempt to Subtract non-commensurable quantities")
	}

	if !q.unit.SameBaseUnits(o.unit) {
		q = q.WithStandardUnits()
		o = o.WithStandardUnits()
	}

	nv := (&big.Float{}).Sub(q.value, o.value)

	return Quantity{
		unit:  q.unit,
		value: nv,
	}
}

// FormatValue returns a string representation of the value expressed in the
// given unit.
//
// This is intended to provide a convenient way to obtain a string
// representation of a quantity to include when serializing data into an
// application-specific file format that assumes a particular unit rather
// than supporting arbitrary units.
//
// The "format" and "prec" arguments take the same meaning as for
// big.Float.Text. The unit conversion is performed first at high
// precision, before rounding the result during formatting.
//
// The given unit must be commensurable with the quantity's own unit, or
// this method will panic.
func (q Quantity) FormatValue(format byte, prec int, unit *Unit) string {
	if q.unit.base != unit.base || q.unit.scale != unit.scale {
		q = q.Convert(unit)
	}

	return q.value.Text(format, prec)
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

func (q Quantity) GoString() string {
	if q.unit == unitByName[""] {
		return fmt.Sprintf("units.MakeDimensionless((&big.Float{}).Parse(%q))", q.value.String())
	}

	return fmt.Sprintf("units.MakeQuantity((&big.Float{}).Parse(%q), %#v)", q.value.String(), q.unit)
}
