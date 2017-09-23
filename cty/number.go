package cty

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/units"
)

// numberImpl is the typeImpl for numeric values, including those with units.
type numberImpl struct {
	isType
	dim units.Dimensionality
}

func Quantity(dim units.Dimensionality) Type {
	return Type{numberImpl{dim: dim}}
}

func QuantityVal(q units.Quantity) Value {
	return Value{
		v:  q,
		ty: Quantity(q.Unit().Dimensionality()),
	}
}

func (i numberImpl) Name() string {
	name := numberTypeNames[i.dim]
	if name != "" {
		return name
	}

	// If our dimensionality doesn't have a friendly name then we'll just
	// call it "Quantity" and qualify it with the dimensionality string.
	return fmt.Sprintf("Quantity(%s)", i.dim.String())
}

func (i numberImpl) GoString() string {
	name := numberTypeNames[i.dim]
	if name != "" {
		return "cty." + name
	}

	return fmt.Sprintf("cty.Quantity(%#v)", i.dim)
}

func (i numberImpl) Equal(a, b Value) Value {
	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return BoolVal(av.Equal(bv))
}

func (i numberImpl) CanSum(other Type) bool {
	otherNum, isNumber := other.impl.(numberImpl)
	if !isNumber {
		return false
	}
	return i.dim == otherNum.dim
}

func (i numberImpl) CanProduct(other Type) bool {
	_, isNumber := other.impl.(numberImpl)
	return isNumber
}

func (i numberImpl) Add(a, b Value) Value {
	if a.IsUnknown() || b.IsUnknown() {
		return UnknownVal(a.Type())
	}

	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return QuantityVal(av.Add(bv))
}

func (i numberImpl) Subtract(a, b Value) Value {
	if a.IsUnknown() || b.IsUnknown() {
		return UnknownVal(a.Type())
	}

	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return QuantityVal(av.Subtract(bv))
}

func (i numberImpl) Multiply(a, b Value) Value {
	if a.IsUnknown() || b.IsUnknown() {
		retTy := Quantity(a.ty.impl.(numberImpl).dim.Multiply(b.ty.impl.(numberImpl).dim))
		return UnknownVal(retTy)
	}

	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return QuantityVal(av.Multiply(bv))
}

func (i numberImpl) Divide(a, b Value) Value {
	if a.IsUnknown() || b.IsUnknown() {
		retTy := Quantity(a.ty.impl.(numberImpl).dim.Multiply(b.ty.impl.(numberImpl).dim.Reciprocal()))
		return UnknownVal(retTy)
	}

	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return QuantityVal(av.Divide(bv))
}
