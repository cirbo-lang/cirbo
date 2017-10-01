package eval

import (
	"fmt"
	"testing"

	"github.com/cirbo-lang/cirbo/units"

	"github.com/cirbo-lang/cirbo/source"

	"github.com/cirbo-lang/cirbo/cty"
)

func TestExprImpls(t *testing.T) {
	// All of the testing here actually happens at compile time. We dress it
	// up like run-time tests just because that makes this test visible in
	// the test results when we pass.
	tests := []Expr{
		(*binaryOpExpr)(nil),
		(*literalExpr)(nil),
		(*symbolExpr)(nil),
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test)[6:], func(t *testing.T) {
			// Nothing to do! If we managed to compile then we passed the test.
		})
	}
}

func TestSymbolExpr(t *testing.T) {
	// We'll construct a standalone scope and context to stand in for the
	// global scope and context for the sake of this test.
	scope := (*Scope)(nil).NewChild()
	ctx := (*Context)(nil).NewChild()

	sym := scope.Declare("foo")
	ctx.Define(sym, cty.True)
	undef := scope.Declare("undefined")

	tests := []struct {
		Expr      Expr
		Want      cty.Value
		DiagCount int
	}{
		{
			SymbolExpr(sym, source.NilRange),
			cty.True,
			0,
		},
		{
			SymbolExpr(undef, source.NilRange),
			cty.PlaceholderVal,
			1, // symbol is not defined
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Expr), func(t *testing.T) {
			got, diags := test.Expr.value(ctx, nil)
			assertDiagCount(t, diags, test.DiagCount)
			assertExprResult(t, test.Expr, got, test.Want)
		})
	}
}

func TestLiteralExpr(t *testing.T) {
	tests := []struct {
		Expr      Expr
		Want      cty.Value
		DiagCount int
	}{
		{
			LiteralExpr(cty.True, source.NilRange),
			cty.True,
			0,
		},
		{
			LiteralExpr(cty.UnknownVal(cty.Bool), source.NilRange),
			cty.UnknownVal(cty.Bool),
			0,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Expr), func(t *testing.T) {
			got, diags := test.Expr.value(GlobalContext(), nil)
			assertDiagCount(t, diags, test.DiagCount)
			assertExprResult(t, test.Expr, got, test.Want)
		})
	}
}

func TestOperatorExpr(t *testing.T) {
	tests := []struct {
		Expr      Expr
		Want      cty.Value
		DiagCount int
	}{

		{
			AddExpr(litExp(cty.One), litExp(cty.One), source.NilRange),
			cty.NumberValInt(2),
			0,
		},
		{
			AddExpr(
				litExp(cty.QuantityVal(units.MakeQuantityInt(2, units.Meter))),
				litExp(cty.QuantityVal(units.MakeQuantityInt(1, units.Meter))),
				source.NilRange,
			),
			cty.QuantityVal(units.MakeQuantityInt(3, units.Meter)),
			0,
		},
		{
			AddExpr(
				litExp(cty.One),
				litExp(cty.QuantityVal(units.MakeQuantityInt(1, units.Meter))),
				source.NilRange,
			),
			cty.PlaceholderVal,
			1, // can't add Number to Length
		},

		{
			SubtractExpr(litExp(cty.NumberValInt(4)), litExp(cty.One), source.NilRange),
			cty.NumberValInt(3),
			0,
		},

		{
			MultiplyExpr(litExp(cty.NumberValInt(4)), litExp(cty.NumberValInt(2)), source.NilRange),
			cty.NumberValInt(8),
			0,
		},
		{
			MultiplyExpr(
				litExp(cty.QuantityVal(units.MakeQuantityInt(4, units.Meter))),
				litExp(cty.QuantityVal(units.MakeQuantityInt(3, units.Meter))),
				source.NilRange,
			),
			cty.QuantityVal(units.MakeQuantityInt(12, units.Meter.ToPower(2))),
			0,
		},

		{
			DivideExpr(litExp(cty.NumberValInt(12)), litExp(cty.NumberValInt(2)), source.NilRange),
			cty.NumberValInt(6),
			0,
		},
		{
			DivideExpr(
				litExp(cty.QuantityVal(units.MakeQuantityInt(24, units.Meter))),
				litExp(cty.QuantityVal(units.MakeQuantityInt(12, units.Second))),
				source.NilRange,
			),
			cty.QuantityVal(units.MakeQuantityInt(2, units.Meter.Multiply(units.Second.ToPower(-1)))),
			0,
		},

		{
			ConcatExpr(litExp(cty.StringVal("ab")), litExp(cty.StringVal("cde")), source.NilRange),
			cty.StringVal("abcde"),
			0,
		},
		{
			ConcatExpr(litExp(cty.True), litExp(cty.False), source.NilRange),
			cty.PlaceholderVal,
			1, // cannot concatenate Bool values
		},

		{
			EqualExpr(litExp(cty.True), litExp(cty.False), source.NilRange),
			cty.False,
			0,
		},
		{
			EqualExpr(litExp(cty.True), litExp(cty.True), source.NilRange),
			cty.True,
			0,
		},
		{
			NotEqualExpr(litExp(cty.True), litExp(cty.True), source.NilRange),
			cty.False,
			0,
		},

		{
			AndExpr(litExp(cty.True), litExp(cty.False), source.NilRange),
			cty.False,
			0,
		},
		{
			AndExpr(litExp(cty.True), litExp(cty.True), source.NilRange),
			cty.True,
			0,
		},
		{
			AndExpr(litExp(cty.False), litExp(cty.False), source.NilRange),
			cty.False,
			0,
		},
		{
			AndExpr(litExp(cty.True), litExp(cty.Zero), source.NilRange),
			cty.UnknownVal(cty.Bool),
			1, // invalid operand types
		},

		{
			OrExpr(litExp(cty.True), litExp(cty.False), source.NilRange),
			cty.True,
			0,
		},
		{
			OrExpr(litExp(cty.True), litExp(cty.True), source.NilRange),
			cty.True,
			0,
		},
		{
			OrExpr(litExp(cty.True), litExp(cty.True), source.NilRange),
			cty.True,
			0,
		},
		{
			OrExpr(litExp(cty.True), litExp(cty.Zero), source.NilRange),
			cty.UnknownVal(cty.Bool),
			1, // invalid operand types
		},

		{
			NotExpr(litExp(cty.True), source.NilRange),
			cty.False,
			0,
		},
		{
			NotExpr(litExp(cty.False), source.NilRange),
			cty.True,
			0,
		},
		{
			NotExpr(litExp(cty.Zero), source.NilRange),
			cty.UnknownVal(cty.Bool),
			1, // invalid operand type
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Expr), func(t *testing.T) {
			got, diags := test.Expr.value(GlobalContext(), nil)
			assertDiagCount(t, diags, test.DiagCount)
			assertExprResult(t, test.Expr, got, test.Want)
		})
	}
}

func litExp(v cty.Value) Expr {
	return LiteralExpr(v, source.NilRange)
}

func assertDiagCount(t *testing.T, diags source.Diags, want int) bool {
	t.Helper()
	if len(diags) != want {
		t.Errorf("unexpected diagnostics count %d; want %d", len(diags), want)
		for _, diag := range diags {
			t.Logf(fmt.Sprintf("- %s", diag))
		}
		return false
	}
	return true
}

func assertExprResult(t *testing.T, e Expr, got cty.Value, want cty.Value) bool {
	t.Helper()
	if !got.Same(want) {
		t.Errorf("wrong expression result\nexpr: %#v\ngot:  %#v\nwant: %#v", e, got, want)
		return false
	}
	return true
}
