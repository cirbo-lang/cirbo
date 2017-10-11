package eval

import (
	"testing"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

func TestExprValue(t *testing.T) {
	scope := GlobalScope().NewChild()
	sym := scope.Declare("foo")

	ctx := GlobalContext().NewChild()
	ctx.DefineLiteral(sym, cbty.True)

	{
		expr := SymbolExpr(sym, source.NilRange)
		got, diags := expr.Value(ctx)
		want := cbty.True
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := LiteralExpr(cbty.False, source.NilRange)
		got, diags := expr.Value(ctx)
		want := cbty.False
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := ConcatExpr(
			LiteralExpr(cbty.True, source.NilRange),
			LiteralExpr(cbty.False, source.NilRange),
			source.NilRange,
		)
		got, diags := expr.Value(ctx)
		want := cbty.PlaceholderVal
		assertDiagCount(t, diags, 1) // invalid operand types
		assertExprResult(t, expr, got, want)
	}
}
