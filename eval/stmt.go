package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cbty"
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

func (s Stmt) RequiredSymbols(scope *Scope) SymbolSet {
	return s.s.requiredSymbols(scope)
}

type stmtImpl interface {
	definedSymbol() *Symbol
	requiredSymbols(scope *Scope) SymbolSet
	execute(exec *StmtBlockExecute, result *StmtBlockResult) source.Diags
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

func (s *assignStmt) execute(exec *StmtBlockExecute, result *StmtBlockResult) source.Diags {
	_, diags := exec.Context.Define(s.sym, s.expr)
	return diags
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

func (s *importStmt) execute(exec *StmtBlockExecute, result *StmtBlockResult) source.Diags {
	var diags source.Diags
	val, defined := exec.Packages[s.ppath]
	if !defined {
		val = cbty.PlaceholderVal
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Module not loaded",
			Detail:  fmt.Sprintf("The module %q has not yet been loaded. This is a bug in Cirbo that should be reported!", s.ppath),
			Ranges:  s.sourceRange().List(),
		})
	}
	exec.Context.DefineLiteral(s.sym, val)
	return diags
}

type exportStmt struct {
	value Expr
	rng
	nonDefStmt
}

func ExportStmt(value Expr, rng source.Range) Stmt {
	return Stmt{&exportStmt{
		value: value,
		rng:   srcRange(rng),
	}}
}

func (s *exportStmt) execute(exec *StmtBlockExecute, result *StmtBlockResult) source.Diags {
	val, diags := s.value.Value(exec.Context)
	if result.ExportValue != cbty.NilValue {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Duplicate export statements",
			Detail:  "Only one export statement is permitted per module.",
			Ranges:  s.sourceRange().List(),
		})
		return diags
	}
	result.ExportValue = val
	return diags
}

func (s *exportStmt) requiredSymbols(scope *Scope) SymbolSet {
	return s.value.RequiredSymbols(scope)
}

type attrStmt struct {
	sym *Symbol

	// Either defValue or valueType must be set, and not both
	defValue, valueType Expr

	rng
}

func AttrStmt(sym *Symbol, valueType Expr, rng source.Range) Stmt {
	return Stmt{&attrStmt{
		sym:       sym,
		valueType: valueType,
		rng:       srcRange(rng),
	}}
}

func AttrStmtDefault(sym *Symbol, defValue Expr, rng source.Range) Stmt {
	return Stmt{&attrStmt{
		sym:      sym,
		defValue: defValue,
		rng:      srcRange(rng),
	}}
}

func (s *attrStmt) definedSymbol() *Symbol {
	return s.sym
}

func (s *attrStmt) requiredSymbols(scope *Scope) SymbolSet {
	switch {
	case s.valueType != NilExpr:
		return s.valueType.RequiredSymbols(scope)
	case s.defValue != NilExpr:
		return s.defValue.RequiredSymbols(scope)
	default:
		return nil // should never happen
	}
}

func (s *attrStmt) execute(exec *StmtBlockExecute, result *StmtBlockResult) source.Diags {
	var diags source.Diags
	var val cbty.Value

	// Attribute values and types are always resolved in the parent context,
	// since they aren't allowed to refer to symbols in the current scope.
	ctx := exec.Context.parent
	if ctx == nil {
		// should never happen
		panic("attempt to execute attr statement in the global context")
	}

	// We allow Attrs to be nil so callers can evaluate statement blocks in
	// isolation for type checking and other similar purposes. In this case,
	// some of the result values will be unknown.
	if exec.Attrs != nil {
		val = exec.Attrs[s.sym]
	}

	var ty cbty.Type
	if s.valueType != NilExpr {
		tyVal, exprDiags := s.valueType.Value(exec.Context)
		diags = append(diags, exprDiags...)
		if tyVal.Type() != cbty.TypeType {
			// This is also checked by stmtBlock.Attributes, so in normal
			// codepaths we'll never get here but we still need to produce
			// a reasonable result in case we're in an analysis codepath that
			// wants to extract as much as it can from a broken program.
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Invalid attribute type",
				Detail:  fmt.Sprintf("Expected a type, but given a value of type %s. To assign a default value, use the '=' (equals) symbol.", tyVal.Type().Name()),
				Ranges:  s.valueType.sourceRange().List(),
			})
			exec.Context.DefineLiteral(s.sym, cbty.PlaceholderVal)
			return diags
		}
		ty = tyVal.UnwrapType()
		if val == cbty.NilValue {
			val = cbty.UnknownVal(ty)
		}
	} else if s.defValue != NilExpr {
		defVal, exprDiags := s.defValue.Value(exec.Context)
		diags = append(diags, exprDiags...)
		ty = defVal.Type()
		if val == cbty.NilValue {
			val = defVal
		}
	}

	if !val.Type().Same(ty) {
		// This should actually never happen because the caller should've
		// already type-checked the call arguments before we get in here,
		// but we'll check anyway. Checking in here is bad because we report
		// the error from the perspective of the declaration rather
		// then from the perspective of the call, and that's not helpful
		// to the user.

		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Invalid attribute value",
			Detail:  fmt.Sprintf("Attribute %q expects a value of type %s, not %s.", s.sym.DeclaredName(), ty.Name(), val.Type().Name()),
			Ranges:  s.sourceRange().List(),
		})

		// Place an unknown value of the correct type to suppress any
		// downstream errors that would otherwise result from this problem.
		val = cbty.UnknownVal(ty)
	}

	exec.Context.DefineLiteral(s.sym, val)

	return diags
}
