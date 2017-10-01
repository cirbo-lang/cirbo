package eval

import (
	"testing"

	"github.com/cirbo-lang/cirbo/source"

	"github.com/cirbo-lang/cirbo/cty"
)

func TestExprValue(t *testing.T) {
	scope := GlobalScope().NewChild()
	sym := scope.Declare("foo")

	ctx := GlobalContext().NewChild()
	ctx.DefineLiteral(sym, cty.True)

	{
		expr := SymbolExpr(sym, source.NilRange)
		got, diags := expr.Value(ctx)
		want := cty.True
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := LiteralExpr(cty.False, source.NilRange)
		got, diags := expr.Value(ctx)
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
		got, diags := expr.Value(ctx)
		want := cty.PlaceholderVal
		assertDiagCount(t, diags, 1) // invalid operand types
		assertExprResult(t, expr, got, want)
	}
}
