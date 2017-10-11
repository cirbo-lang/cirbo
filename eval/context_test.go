package eval

import (
	"testing"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

func TestContextDefine(t *testing.T) {
	scope := GlobalScope().NewChild()
	sym1 := scope.Declare("sym1")
	sym2 := scope.Declare("sym2")

	ctx := GlobalContext().NewChild()

	{
		expr := LiteralExpr(cbty.True, source.NilRange)
		want := cbty.True

		got, diags := ctx.Define(sym1, expr)
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)

		got = ctx.Value(sym1)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := NotExpr(LiteralExpr(cbty.Zero, source.NilRange), source.NilRange)
		want := cbty.UnknownVal(cbty.Bool)

		got, diags := ctx.Define(sym2, expr)
		assertDiagCount(t, diags, 1) // invalid operand type
		assertExprResult(t, expr, got, want)

		got = ctx.Value(sym2)
		assertExprResult(t, expr, got, want)
	}
}
