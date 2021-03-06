package cbty

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/units"
)

type Type struct {
	impl typeImpl
}

// NilType is an invalid type that serves as the zero value of type Type.
//
// NilType is not a real type and so is used only to signal the absense of
// a type when returning from functions.
var NilType Type

// Name returns a name for the receiving type that is suitable for display
// to cirbo end-users.
func (t Type) Name() string {
	return t.impl.Name()
}

// Same returns true if and only if the given type is the same as the
// receiver.
func (t Type) Same(o Type) bool {
	type Samer interface {
		Same(o Type) bool
	}

	if s, canSame := t.impl.(Samer); canSame {
		return s.Same(o)
	}

	// Default implementation works for simple typeImpls; will panic if
	// the impl does not support ==, so such impls must implement the
	// Same method from above.
	return t.impl == o.impl
}

// GoString returns a representation of the receiving type as Go syntax,
// suitable for display in tests and other internal debug messages.
//
// This result must never be displayed to cirbo end-users.
func (t Type) GoString() string {
	if s, isStringer := t.impl.(fmt.GoStringer); isStringer {
		return s.GoString()
	}
	return fmt.Sprintf("cty.%s", t.Name())
}

// HasArithmetic returns true if and only if the recieving type supports
// the arithmetic operators.
func (t Type) HasArithmetic() bool {
	_, has := t.impl.(typeWithArithmetic)
	return has
}

// CanSum returns true if the given type can support the Add and Subtract
// operations with values of the other given type.
//
// Always returns false if the receiver doesn't support arithmetic at all.
func (t Type) CanSum(o Type) bool {
	if !t.HasArithmetic() {
		return false
	}
	return t.impl.(typeWithArithmetic).CanSum(o)
}

// CanProduct returns true if the given type can support the Multiply and Divide
// operations with values of the other given type.
//
// Always returns false if the receiver doesn't support arithmetic at all.
func (t Type) CanProduct(o Type) bool {
	if !t.HasArithmetic() {
		return false
	}
	return t.impl.(typeWithArithmetic).CanProduct(o)
}

// CanConcat returns true if the given type can concatenate values of the other
// given type.
//
// Always returns false if the receiver doesn't support concatenation at all.
func (t Type) CanConcat(o Type) bool {
	concatter, canConcat := t.impl.(typeWithConcat)
	if !canConcat {
		return false
	}
	return concatter.CanConcat(o)
}

// HasAttr returns true if the receiver has an attribute of the given name.
func (t Type) HasAttr(name string) bool {
	return t.AttrType(name) != NilType
}

// AttrType returns the type of the attribute of the given name, or NilType
// if the receiver has no such attribute.
func (t Type) AttrType(name string) Type {
	withAttrs, has := t.impl.(typeWithAttributes)
	if !has {
		return NilType
	}

	uv := withAttrs.GetAttr(UnknownVal(t), name)
	if uv == NilValue {
		return NilType
	}

	return uv.Type()
}

// CallSignature returns the expected signature for calls to values of the
// recieving type, or nil if the type cannot be called at all.
func (t Type) CallSignature() *CallSignature {
	callImpl, isCallable := t.impl.(typeCallable)
	if !isCallable {
		return nil
	}

	return callImpl.CallSignature()
}

// IsModel returns true if the receiver is a model type.
func (t Type) IsModel() bool {
	_, isModel := t.impl.(*modelImpl)
	return isModel
}

// ModelImpl returns the underlying implementation of the receiver if it
// is a model type, or nil otherwise.
//
// This should only be used by the package that created the type, to avoid
// exposing implementation details.
func (t Type) ModelImpl() interface{} {
	m, isModel := t.impl.(*modelImpl)
	if !isModel {
		return nil
	}
	return m.pubImpl
}

// IsNumber returns true if the receiver is a number type. Number includes
// both dimensionless numbers and dimensioned quantities.
func (t Type) IsNumber() bool {
	_, isNumber := t.impl.(numberImpl)
	return isNumber
}

// NumberDimensionality returns the dimensionality of the receiving number
// type, or panics if the receiver is not a number type.
func (t Type) NumberDimensionality() units.Dimensionality {
	impl, isNumber := t.impl.(numberImpl)
	if !isNumber {
		panic("IsNumber on non-number value")
	}

	return impl.dim
}

type typeImpl interface {
	typeSigil() isType
	Name() string
	Equal(a, b Value) Value
}

type isType struct {
}

func (it isType) typeSigil() isType {
	return it
}

// GetAttr is a default implementation of GetAttr that always returns NilValue,
// indicating that the associated value has no attributes.
//
// Override this with another implementation of GetAttr in order to actually
// provide attributes. For example, embedding staticAttributes allows
// attributes to be provided as a map.
func (it isType) GetAttr(name string) Value {
	return NilValue
}
