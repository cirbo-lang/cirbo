package eval

// nonDefStmt can be embedded into a statement type that does not define
// anything, to get a do-nothing implementation of definedSymbol.
type nonDefStmt struct {
}

func (nonDefStmt) definedSymbol() *Symbol {
	return nil
}

// nonExprStmt can be embedded into a statement type that does not have
// any expressions, and thus get a do-nothing implementation of requiredSymbols.
type nonExprStmt struct {
}

func (nonExprStmt) requiredSymbols(*Scope) SymbolSet {
	return nil
}
