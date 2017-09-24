package cty

import "github.com/cirbo-lang/cirbo/source"

type functionImpl struct {
	isType
	sig *CallSignature
}

// Function returns a function type with the given call signature.
func Function(sig *CallSignature) Type {
	return Type{functionImpl{sig: sig}}
}

// FunctionImpl represents the implementation of a particular function.
type FunctionImpl struct {
	Signature *CallSignature
	Callback  func(args CallArgs) (Value, source.Diags)
}

// FunctionVal creates a function value with the given implementation.
func FunctionVal(impl FunctionImpl) Value {
	return Value{
		ty: Function(impl.Signature),
		v:  &impl,
	}
}

func (i functionImpl) Name() string {
	// The full description of a function type is too complex for a compact
	// string serialization, so we'll just simplify here and assume that if
	// we're producing an error relating to a function's signature we will
	// refer specifically to individual parameters, rather than talking about
	// the type as a whole.
	return "Function"
}

func (i functionImpl) Same(o Type) bool {
	oi, isFunc := o.impl.(functionImpl)
	if !isFunc {
		return false
	}

	return i.sig.Same(oi.sig)
}

func (i functionImpl) Equal(a, b Value) Value {
	// Each call to FunctionVal produces a distinct function value, and
	// functions are compared by identity.
	return BoolVal(a.v == b.v)
}

func (i functionImpl) CallSignature() *CallSignature {
	return i.sig
}

func (i functionImpl) Call(callee Value, args CallArgs) (Value, source.Diags) {
	impl := callee.v.(*FunctionImpl)
	return impl.Callback(args)
}
