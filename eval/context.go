package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cty"
)

// Context represents the current values for a set of symbols used during
// evaluation.
//
// A Context is roughly analogous to a stack frame, allowing local values
// to be defined for a particular call/instantiation without interfering with
// values for the same symbol in other possibly-concurrent calls.
type Context struct {
	parent *Context
	values map[*Symbol]cty.Value
}

// NewChild creates a new, empty context that is a child of the receiever.
func (ctx *Context) NewChild() *Context {
	return &Context{
		parent: ctx,
		values: map[*Symbol]cty.Value{},
	}
}

// Define records a value for the given symbol in the receiver.
//
// This method will panic if the given symbol is already defined in the
// receiver. This does not apply if it is defined only in parent contexts.
// Use Defined to determine if a given symbol is already defined.
//
// This method will also panic if the receiver is the immutable global context.
func (ctx *Context) Define(sym *Symbol, val cty.Value) {
	if ctx == globalContext {
		panic(fmt.Errorf("attempt to define %#v as %#v in the immutable global scope", sym, val))
	}
	if _, defined := ctx.values[sym]; defined {
		panic(fmt.Errorf("attempt to re-define %#v as %#v in context %#v", sym, val, ctx))
	}

	ctx.values[sym] = val
}

// Defined returns true if the given symbol is defined directly in the receiver.
// It returns false if the symbol is defined only in ancestor contexts or if
// no context in the inheritance chain defines it.
func (ctx *Context) Defined(sym *Symbol) bool {
	_, defined := ctx.values[sym]
	return defined
}

// Value returns the value of the given symbol in the receiver or the nearest
// defining ancestor context, or NilValue if the given symbol is not yet
// defined in any context in the inheritance chain.
func (ctx *Context) Value(sym *Symbol) cty.Value {
	current := ctx
	for current != nil {
		val, has := current.values[sym]
		if has {
			return val
		}
		current = current.parent
	}

	return cty.NilValue
}

// AllValues returns a map describing the values of all of all of the symbols
// defined in the given scope, using their definition names.
//
// If any of the symbols are not yet defined, they will map to NilValue.
func (ctx *Context) AllValues(s *Scope) map[string]cty.Value {
	ret := make(map[string]cty.Value, len(s.symbols))
	for name, symbol := range s.symbols {
		ret[name] = ctx.Value(symbol)
	}
	return ret
}