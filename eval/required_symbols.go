package eval

// RequiredSymbols returns a set of symbols from the given scope that are
// required by the given expression.
//
// This can be used to understand the dependency relationships between different
// symbol definitions, and thus to process definitions in a suitable order
// to ensure that all prerequisites are satisfied.
func RequiredSymbols(expr Expr, scope *Scope) SymbolSet {
	ret := make(SymbolSet)

	var cb walkCb
	cb = func(oe Expr) {
		if se, ok := oe.(*symbolExpr); ok {
			if se.sym.scope == scope {
				ret.Add(se.sym)
			}
		}
		oe.eachChild(cb)
	}
	expr.eachChild(cb)

	return ret
}

// SymbolReferences returns a set of sub-expressions within the given
// expression that refer to the given symbol.
//
// The primary reason to use this function is to obtain the source
// locations of invalid references.
func SymbolReferences(expr Expr, sym *Symbol) ExprSet {
	ret := make(ExprSet)

	var cb walkCb
	cb = func(oe Expr) {
		if se, ok := oe.(*symbolExpr); ok {
			if se.sym == sym {
				ret.Add(oe)
			}
		}
		oe.eachChild(cb)
	}
	expr.eachChild(cb)

	return ret
}
