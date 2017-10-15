package cbty

import (
	"github.com/cirbo-lang/cirbo/source"
)

type typeCallable interface {
	CallSignature() *CallSignature
	Call(callee Value, args CallArgs) (Value, source.Diags)
}

// CallSignature describes the formal parameters of a callable value.
type CallSignature struct {
	// Parameters describes the explicit parameters of a callable, keyed
	// by the parameter name.
	Parameters map[string]CallParameter

	// Positional lists the names of the subset of parameters that may be
	// passed positionally, in the order that they are expected. If any
	// of the positional parameters are required, they will always be at
	// the start of this list.
	Positional []string

	// AcceptsVariadicPositional indicates that additional positional arguments
	// (further to those listed explicitly in Positional) are accepted.
	AcceptsVariadicPositional bool

	// AcceptsVariadicNamed indicates that additional named arguments (further
	// to those listed explicitly in Parameters) are accepted.
	AcceptsVariadicNamed bool

	// Result is the type of values returned from any call.
	Result Type
}

// CallParameter describes a single parameter within a CallSignature.
type CallParameter struct {
	Type     Type
	Required bool
}

// CallArgs represents a set of arguments being passed to Value.Call.
type CallArgs struct {
	Explicit      map[string]Value
	PosVariadic   []Value
	NamedVariadic map[string]Value

	// TargetName, if non-empty, is the declared name of a symbol that the
	// call result will be written to. Always empty if the result will not
	// be used to define a symbol
	//
	// Most callables can ignore this, but this can optionally be used to
	// establish the primary name of an object that needs such a name. This
	// includes circuit and device instances since they must be addressable
	// within the assignments file of a project.
	TargetName string
}

// Same returns true if the receiver and the other given signature are
// equivalent.
func (s *CallSignature) Same(o *CallSignature) bool {
	if !s.Result.Same(o.Result) {
		return false
	}

	if s.AcceptsVariadicPositional != o.AcceptsVariadicPositional {
		return false
	}
	if s.AcceptsVariadicNamed != o.AcceptsVariadicNamed {
		return false
	}

	if len(s.Positional) != len(o.Positional) {
		return false
	}
	for i := range s.Positional {
		if s.Positional[i] != o.Positional[i] {
			return false
		}
	}

	if len(s.Parameters) != len(o.Parameters) {
		return false
	}
	for n := range s.Parameters {
		if s.Parameters[n].Required != o.Parameters[n].Required {
			return false
		}

		sty := s.Parameters[n].Type
		oty := o.Parameters[n].Type

		if !sty.Same(oty) {
			return false
		}
	}

	// If we got through all the above without returning then our signatures
	// are indeed identical.
	return true
}
