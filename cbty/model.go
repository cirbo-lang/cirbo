package cbty

import (
	"fmt"
	"reflect"

	"github.com/cirbo-lang/cirbo/source"
)

// modelImpl is the typeImpl for exposing model objects (e.g. from the cbo
// package) into the language. Model types are based on Go types, and support
// only the equality test and attribute access operations.
type modelImpl struct {
	isType
	pubImpl ModelImpl
}

// ModelImpl is an interface that can be implemented in order to create a
// model type, which is a type that wraps a Go object (usually a model from
// the cbo package) so it can be passed around in the language and have
// attributes accessed on it.
type ModelImpl interface {
	Name() string
	SuitableValue(raw interface{}) bool

	GetAttr(raw interface{}, name string) Value

	CallSignature() *CallSignature
	Call(callee interface{}, args CallArgs) (Value, source.Diags)
}

// Model creates a new model type with the given model implementation.
//
// Each call to Model produces a distinct type. That is, Same will return true
// with two Type values that were the result of the same call, but false for
// any two Type values that were the result of different calls, even if the
// underlying implementation is identical.
func Model(impl ModelImpl) Type {
	// modelImpl is used via a pointer to ensure that each call to Model
	// produces a distinct type.
	return Type{&modelImpl{
		pubImpl: impl,
	}}
}

// ModelVal creates a new value of a model type.
//
// The given raw value must be a pointer and must be of a type appropriate
// for the given model type, or else this function will panic. What constitutes
// an "appropriate" value depends on the model type.
//
// This function will panic also if the given type is not a model type.
func ModelVal(ty Type, raw interface{}) Value {
	if reflect.ValueOf(raw).Kind() != reflect.Ptr {
		panic("ModelVal raw value must be pointer")
	}
	if i, isModel := ty.impl.(*modelImpl); isModel {
		if !i.pubImpl.SuitableValue(raw) {
			panic(fmt.Sprintf("ModelVal called with unsuitable %s value %#v", i.Name(), raw))
		}

		return Value{
			ty: ty,
			v:  raw,
		}
	}

	panic("ModelVal called with non-model type")
}

func (i *modelImpl) Name() string {
	return i.pubImpl.Name()
}

func (i *modelImpl) Equal(a, b Value) Value {
	return BoolVal(a.v == b.v)
}

func (i *modelImpl) GetAttr(val Value, name string) Value {
	var raw interface{}
	if val.IsKnown() {
		raw = val.v
	}
	return i.pubImpl.GetAttr(raw, name)
}

func (i *modelImpl) CallSignature() *CallSignature {
	return i.pubImpl.CallSignature()
}

func (i *modelImpl) Call(callee Value, args CallArgs) (Value, source.Diags) {
	return i.pubImpl.Call(callee.v, args)
}

func (i *modelImpl) GoString() string {
	return fmt.Sprintf("cbty.Model(%#v)", i.pubImpl)
}
