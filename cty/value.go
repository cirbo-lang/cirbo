package cty

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/units"
)

type Value struct {
	ty Type
	v  interface{}
}

var NilValue Value

// Type returns the type of the receiever.
func (v Value) Type() Type {
	return v.ty
}

// Same returns true if and only if the given value and the reciever are
// identical.
//
// Identity is different than equality in that two unknown values are
// identical if their types are identical, whereas an equality test would
// return an unknown boolean value.
//
// This method is primarily for test assertions. All code implementing the
// language itself should use Equal and handle unknown values.
func (v Value) Same(o Value) bool {
	if !v.SameType(o) {
		return false
	}

	if v.v == nil || o.v == nil {
		return v.v == o.v
	}

	type Samer interface {
		ValueSame(a, b Value) bool
	}

	if s, canSame := v.ty.impl.(Samer); canSame {
		return s.ValueSame(v, o)
	}

	eq := v.Equal(o)

	if eq.IsUnknown() {
		return false
	}
	return eq.True()
}

// Equal returns True if and only if the value and the receiver represent
// the same value.
//
// No two values of different types are ever equal. If either value is
// unknown then the result itself is an unknown boolean.
func (v Value) Equal(o Value) Value {
	if !v.SameType(o) {
		return False
	}

	if v.IsUnknown() || o.IsUnknown() {
		return UnknownVal(Bool)
	}

	return v.ty.impl.Equal(v, o)
}

// True returns true if the receiver is True, false if the receiver is False,
// and panics otherwise.
//
// This method converts a known cty.Bool value into a native Go bool value.
func (v Value) True() bool {
	switch {
	case v == True:
		return true
	case v == False:
		return false
	case !v.IsKnown():
		panic("True called on unknown value")
	default:
		panic("True called on non-boolean value")
	}
}

// IsKnown returns true if the receiver is a known value.
//
// If false is returned, only the type is known.
func (v Value) IsKnown() bool {
	return v.v != nil
}

// IsUnknown is the opposite of IsKnown, for convenience.
func (v Value) IsUnknown() bool {
	return v.v == nil
}

// SameType returns true if and only if the given value has the same type
// as the reciever.
func (v Value) SameType(o Value) bool {
	return v.ty.Same(o.ty)
}

// Add returns the sum of the receiver and the given other value.
//
// This function will panic if the value type does not support arithmetic or
// cannot add a value of the other type.
func (v Value) Add(o Value) Value {
	if !v.Type().CanSum(o.Type()) {
		panic(fmt.Errorf("attempt to add %#v and %#v", v.Type(), o.Type()))
	}

	return v.ty.impl.(typeWithArithmetic).Add(v, o)
}

// Subtract returns the difference between the receiver and the given other value.
//
// This function will panic if the value type does not support arithmetic or
// cannot subtract a value of the other type.
func (v Value) Subtract(o Value) Value {
	if !v.Type().CanSum(o.Type()) {
		panic(fmt.Errorf("attempt to subtract %#v from %#v", v.Type(), o.Type()))
	}

	return v.ty.impl.(typeWithArithmetic).Subtract(v, o)
}

// Multiply returns the product of the receiver and the given other value.
//
// This function will panic if the value type does not support arithmetic or
// cannot multiply a value of the other type.
func (v Value) Multiply(o Value) Value {
	if !v.Type().CanProduct(o.Type()) {
		panic(fmt.Errorf("attempt to multiply %#v and %#v", v.Type(), o.Type()))
	}

	return v.ty.impl.(typeWithArithmetic).Multiply(v, o)
}

// Divide returns the quotient of the receiver by the given other value.
//
// This function will panic if the value type does not support arithmetic or
// cannot divide by a value of the other type.
func (v Value) Divide(o Value) Value {
	if !v.Type().CanProduct(o.Type()) {
		panic(fmt.Errorf("attempt to divide %#v by %#v", v.Type(), o.Type()))
	}

	return v.ty.impl.(typeWithArithmetic).Divide(v, o)
}

// Concat concatenates the other given value onto the end of the reciever
// and returns the result.
//
// This function will panic if the receiver type does not support concatenation
// with the other value's type.
func (v Value) Concat(o Value) Value {
	if !v.Type().CanConcat(o.Type()) {
		panic(fmt.Errorf("attempt to concatenate %#v onto %#v", o.Type(), v.Type()))
	}

	return v.ty.impl.(typeWithConcat).Concat(v, o)
}

// Not returns the inverse of the reciever, which must be of type Bool.
//
// If the receiver is not a boolean value, this method will panic.
func (v Value) Not() Value {
	if !v.Type().Same(Bool) {
		panic(fmt.Errorf("attempt to NOT %#v", v.Type()))
	}

	return v.ty.impl.(boolImpl).Not(v)
}

// And returns True if the receiver and the other given value are both True,
// an unknown Bool if either is unknown, or False otherwise.
//
// If the either value is not of type Bool, this method will panic.
func (v Value) And(o Value) Value {
	if !v.Type().Same(Bool) || !o.Type().Same(Bool) {
		panic(fmt.Errorf("attempt to AND %#v and %#v", v.Type(), o.Type()))
	}

	return v.ty.impl.(boolImpl).And(v, o)
}

// Or returns True if either the receiver or the other given value is True,
// or if both are true. It returns an unknown Bool if either is unknown.
// Otherwise, it returns False.
//
// If the either value is not of type Bool, this method will panic.
func (v Value) Or(o Value) Value {
	if !v.Type().Same(Bool) || !o.Type().Same(Bool) {
		panic(fmt.Errorf("attempt to OR %#v and %#v", v.Type(), o.Type()))
	}

	return v.ty.impl.(boolImpl).Or(v, o)
}

func (v Value) GoString() string {
	_, isQuantity := v.ty.impl.(numberImpl)
	switch {
	case v.IsUnknown():
		return fmt.Sprintf("cty.UnknownVal(%#v)", v.Type())
	case isQuantity:
		quantity := v.v.(units.Quantity)
		return fmt.Sprintf("cty.Quantity(%#v)", quantity)
	case v == True:
		return "cty.True"
	case v == False:
		return "cty.False"
	case v.Type().Same(String):
		return fmt.Sprintf("cty.StringVal(%q)", v.v)
	default:
		return fmt.Sprintf("cty.Value{ty: %#v, v: %#v}", v.ty, v.v)
	}
}
