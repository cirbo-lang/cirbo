package eval

import (
	"github.com/cirbo-lang/cirbo/source"
)

type Stmt struct {
	s stmtImpl
}

// NilStmt is an invalid expression that serves as the zero value of Stmt.
//
// NilStmt indicates the absense of an expression and is not itself a valid
// expression. Any methods called on it will panic.
var NilStmt Stmt

type stmtImpl interface {
	definedSymbol() *Symbol
	requiredSymbols(scope *Scope) SymbolSet
	sourceRange() source.Range
}

type assignStmt struct {
	sym  *Symbol
	expr Expr
	rng
}

func AssignStmt(sym *Symbol, expr Expr, rng source.Range) Stmt {
	return Stmt{&assignStmt{
		sym:  sym,
		expr: expr,
		rng:  srcRange(rng),
	}}
}

func (s *assignStmt) definedSymbol() *Symbol {
	return s.sym
}

func (s *assignStmt) requiredSymbols(scope *Scope) SymbolSet {
	return s.expr.RequiredSymbols(scope)
}

type importStmt struct {
	ppath string
	sym   *Symbol
	rng
	nonExprStmt
}

func ImportStmt(ppath string, sym *Symbol, rng source.Range) Stmt {
	return Stmt{&importStmt{
		ppath: ppath,
		sym:   sym,
		rng:   srcRange(rng),
	}}
}

func (s *importStmt) definedSymbol() *Symbol {
	return s.sym
}
