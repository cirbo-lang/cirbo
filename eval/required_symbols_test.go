package eval

import (
	"testing"

	"github.com/cirbo-lang/cirbo/source"
)

func TestRequiredSymbols(t *testing.T) {
	scope1 := GlobalScope().NewChild()
	scope2 := scope1.NewChild()

	symScope1Only := scope1.Declare("scope1Only")
	symScope2Only := scope2.Declare("scope2Only")
	symBothScopes1 := scope1.Declare("bothScopes")
	symBothScopes2 := scope2.Declare("bothScopes")
	scope2.Declare("notUsed")

	expr := AddExpr(
		AddExpr(
			SymbolExpr(symScope1Only, source.NilRange),
			SymbolExpr(symBothScopes1, source.NilRange),
			source.NilRange,
		),
		AddExpr(
			SymbolExpr(symScope2Only, source.NilRange),
			SymbolExpr(symBothScopes2, source.NilRange),
			source.NilRange,
		),
		source.NilRange,
	)

	got := expr.RequiredSymbols(scope2)
	want := NewSymbolSet(symScope2Only, symBothScopes2)

	if !got.Equal(want) {
		t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got.Names(), want.Names())
	}
}

func TestSymbolReferences(t *testing.T) {
	scope := GlobalScope().NewChild()

	symbol1 := scope.Declare("symbol1")
	symbol2 := scope.Declare("symbol2")

	symExpr1a := SymbolExpr(symbol1, source.NilRange)
	symExpr1b := SymbolExpr(symbol1, source.NilRange)
	symExpr2a := SymbolExpr(symbol2, source.NilRange)
	symExpr2b := SymbolExpr(symbol2, source.NilRange)

	expr := AddExpr(
		AddExpr(
			symExpr1a,
			symExpr2a,
			source.NilRange,
		),
		AddExpr(
			symExpr1b,
			symExpr2b,
			source.NilRange,
		),
		source.NilRange,
	)

	got := expr.SymbolReferences(symbol1)
	want := NewExprSet(symExpr1a, symExpr1b)

	if !got.Equal(want) {
		t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, want)
	}
}
