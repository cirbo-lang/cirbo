package eval

import (
	"testing"

	"github.com/cirbo-lang/cirbo/source"

	"github.com/cirbo-lang/cirbo/cty"
)

func TestValue(t *testing.T) {
	scope := GlobalScope().NewChild()
	sym := scope.Declare("foo")

	ctx := GlobalContext().NewChild()
	ctx.Define(sym, cty.True)

	{
		expr := SymbolExpr(sym, source.NilRange)
		got, diags := Value(expr, ctx)
		want := cty.True
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := LiteralExpr(cty.False, source.NilRange)
		got, diags := Value(expr, ctx)
		want := cty.False
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := ConcatExpr(
			LiteralExpr(cty.True, source.NilRange),
			LiteralExpr(cty.False, source.NilRange),
			source.NilRange,
		)
		got, diags := Value(expr, ctx)
		want := cty.PlaceholderVal
		assertDiagCount(t, diags, 1) // invalid operand types
		assertExprResult(t, expr, got, want)
	}
}

func TestDefineSymbol(t *testing.T) {
	scope := GlobalScope().NewChild()
	sym1 := scope.Declare("sym1")
	sym2 := scope.Declare("sym2")

	ctx := GlobalContext().NewChild()

	{
		expr := LiteralExpr(cty.True, source.NilRange)
		want := cty.True

		got, diags := DefineSymbol(sym1, expr, ctx)
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)

		got = ctx.Value(sym1)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := NotExpr(LiteralExpr(cty.Zero, source.NilRange), source.NilRange)
		want := cty.UnknownVal(cty.Bool)

		got, diags := DefineSymbol(sym2, expr, ctx)
		assertDiagCount(t, diags, 1) // invalid operand type
		assertExprResult(t, expr, got, want)

		got = ctx.Value(sym2)
		assertExprResult(t, expr, got, want)
	}
}
