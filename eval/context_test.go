package eval

import (
	"testing"

	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/source"
)

func TestContextDefine(t *testing.T) {
	scope := GlobalScope().NewChild()
	sym1 := scope.Declare("sym1")
	sym2 := scope.Declare("sym2")

	ctx := GlobalContext().NewChild()

	{
		expr := LiteralExpr(cty.True, source.NilRange)
		want := cty.True

		got, diags := ctx.Define(sym1, expr)
		assertDiagCount(t, diags, 0)
		assertExprResult(t, expr, got, want)

		got = ctx.Value(sym1)
		assertExprResult(t, expr, got, want)
	}
	{
		expr := NotExpr(LiteralExpr(cty.Zero, source.NilRange), source.NilRange)
		want := cty.UnknownVal(cty.Bool)

		got, diags := ctx.Define(sym2, expr)
		assertDiagCount(t, diags, 1) // invalid operand type
		assertExprResult(t, expr, got, want)

		got = ctx.Value(sym2)
		assertExprResult(t, expr, got, want)
	}
}
