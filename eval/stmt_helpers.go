package eval

// nonDefStmt can be embedded into a statement type that does not define
// anything, to get a do-nothing implementation of definedSymbol.
type nonDefStmt struct {
}

func (nonDefStmt) definedSymbol() *Symbol {
	return nil
}
