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

func Number(dim units.Dimensionality) Type {
	return Type{numberImpl{dim: dim}}
}

func NumberVal(q units.Quantity) Value {
	return Value{
		v:  q,
		ty: Number(q.Unit().Dimensionality()),
	}
}

func (i numberImpl) Name() string {
	return fmt.Sprintf("Number(%s)", i.dim.String())
}

func (i numberImpl) GoString() string {
	return fmt.Sprintf("cty.Number(%#v)", i.dim)
}

func (i numberImpl) Add(a, b Value) Value {
	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return NumberVal(av.Add(bv))
}

func (i numberImpl) Subtract(a, b Value) Value {
	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return NumberVal(av.Subtract(bv))
}

func (i numberImpl) Multiply(a, b Value) Value {
	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return NumberVal(av.Multiply(bv))
}

func (i numberImpl) Divide(a, b Value) Value {
	av := a.v.(units.Quantity)
	bv := b.v.(units.Quantity)
	return NumberVal(av.Divide(bv))
}
