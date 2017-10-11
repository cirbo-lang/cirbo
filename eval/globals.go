package eval

import (
	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/cbty/globals"
)

// GlobalScope returns the global scope, which is the top-most scope that
// is visible in all source files.
func GlobalScope() *Scope {
	return globalScope
}

// GlobalContext returns the global context, which is a singleton context
// containing the shared values for the global scope.
func GlobalContext() *Context {
	return globalContext
}

var globalScope *Scope
var globalContext *Context

func init() {
	globalScope = &Scope{
		symbols: map[string]*Symbol{},
		final:   true,
	}
	globalContext = &Context{
		values: map[*Symbol]cbty.Value{},
		final:  true,
	}

	vals := globals.Table()

	for name, val := range vals {
		sym := &Symbol{
			scope: globalScope,
			name:  name,
		}
		globalScope.symbols[name] = sym
		globalContext.values[sym] = val
	}
}
