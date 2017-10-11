package eval

import (
	"fmt"
	"testing"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
	"github.com/cirbo-lang/cirbo/units"
)

func TestExprImpls(t *testing.T) {
	// All of the testing here actually happens at compile time. We dress it
	// up like run-time tests just because that makes this test visible in
	// the test results when we pass.
	tests := []exprImpl{
		(*binaryOpExpr)(nil),
		(*callExpr)(nil),
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
	ctx.DefineLiteral(sym, cbty.True)
	undef := scope.Declare("undefined")

	tests := []struct {
		Expr      Expr
		Want      cbty.Value
		DiagCount int
	}{
		{
			SymbolExpr(sym, source.NilRange),
			cbty.True,
			0,
		},
		{
			SymbolExpr(undef, source.NilRange),
			cbty.PlaceholderVal,
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
		Want      cbty.Value
		DiagCount int
	}{
		{
			LiteralExpr(cbty.True, source.NilRange),
			cbty.True,
			0,
		},
		{
			LiteralExpr(cbty.UnknownVal(cbty.Bool), source.NilRange),
			cbty.UnknownVal(cbty.Bool),
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
		Want      cbty.Value
		DiagCount int
	}{

		{
			AddExpr(litExp(cbty.One), litExp(cbty.One), source.NilRange),
			cbty.NumberValInt(2),
			0,
		},
		{
			AddExpr(
				litExp(cbty.QuantityVal(units.MakeQuantityInt(2, units.Meter))),
				litExp(cbty.QuantityVal(units.MakeQuantityInt(1, units.Meter))),
				source.NilRange,
			),
			cbty.QuantityVal(units.MakeQuantityInt(3, units.Meter)),
			0,
		},
		{
			AddExpr(
				litExp(cbty.One),
				litExp(cbty.QuantityVal(units.MakeQuantityInt(1, units.Meter))),
				source.NilRange,
			),
			cbty.PlaceholderVal,
			1, // can't add Number to Length
		},

		{
			SubtractExpr(litExp(cbty.NumberValInt(4)), litExp(cbty.One), source.NilRange),
			cbty.NumberValInt(3),
			0,
		},

		{
			MultiplyExpr(litExp(cbty.NumberValInt(4)), litExp(cbty.NumberValInt(2)), source.NilRange),
			cbty.NumberValInt(8),
			0,
		},
		{
			MultiplyExpr(
				litExp(cbty.QuantityVal(units.MakeQuantityInt(4, units.Meter))),
				litExp(cbty.QuantityVal(units.MakeQuantityInt(3, units.Meter))),
				source.NilRange,
			),
			cbty.QuantityVal(units.MakeQuantityInt(12, units.Meter.ToPower(2))),
			0,
		},

		{
			DivideExpr(litExp(cbty.NumberValInt(12)), litExp(cbty.NumberValInt(2)), source.NilRange),
			cbty.NumberValInt(6),
			0,
		},
		{
			DivideExpr(
				litExp(cbty.QuantityVal(units.MakeQuantityInt(24, units.Meter))),
				litExp(cbty.QuantityVal(units.MakeQuantityInt(12, units.Second))),
				source.NilRange,
			),
			cbty.QuantityVal(units.MakeQuantityInt(2, units.Meter.Multiply(units.Second.ToPower(-1)))),
			0,
		},

		{
			ConcatExpr(litExp(cbty.StringVal("ab")), litExp(cbty.StringVal("cde")), source.NilRange),
			cbty.StringVal("abcde"),
			0,
		},
		{
			ConcatExpr(litExp(cbty.True), litExp(cbty.False), source.NilRange),
			cbty.PlaceholderVal,
			1, // cannot concatenate Bool values
		},

		{
			EqualExpr(litExp(cbty.True), litExp(cbty.False), source.NilRange),
			cbty.False,
			0,
		},
		{
			EqualExpr(litExp(cbty.True), litExp(cbty.True), source.NilRange),
			cbty.True,
			0,
		},
		{
			NotEqualExpr(litExp(cbty.True), litExp(cbty.True), source.NilRange),
			cbty.False,
			0,
		},

		{
			AndExpr(litExp(cbty.True), litExp(cbty.False), source.NilRange),
			cbty.False,
			0,
		},
		{
			AndExpr(litExp(cbty.True), litExp(cbty.True), source.NilRange),
			cbty.True,
			0,
		},
		{
			AndExpr(litExp(cbty.False), litExp(cbty.False), source.NilRange),
			cbty.False,
			0,
		},
		{
			AndExpr(litExp(cbty.True), litExp(cbty.Zero), source.NilRange),
			cbty.UnknownVal(cbty.Bool),
			1, // invalid operand types
		},

		{
			OrExpr(litExp(cbty.True), litExp(cbty.False), source.NilRange),
			cbty.True,
			0,
		},
		{
			OrExpr(litExp(cbty.True), litExp(cbty.True), source.NilRange),
			cbty.True,
			0,
		},
		{
			OrExpr(litExp(cbty.True), litExp(cbty.True), source.NilRange),
			cbty.True,
			0,
		},
		{
			OrExpr(litExp(cbty.True), litExp(cbty.Zero), source.NilRange),
			cbty.UnknownVal(cbty.Bool),
			1, // invalid operand types
		},

		{
			NotExpr(litExp(cbty.True), source.NilRange),
			cbty.False,
			0,
		},
		{
			NotExpr(litExp(cbty.False), source.NilRange),
			cbty.True,
			0,
		},
		{
			NotExpr(litExp(cbty.Zero), source.NilRange),
			cbty.UnknownVal(cbty.Bool),
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

func litExp(v cbty.Value) Expr {
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

func assertExprResult(t *testing.T, e Expr, got cbty.Value, want cbty.Value) bool {
	t.Helper()
	if !got.Same(want) {
		t.Errorf("wrong expression result\nexpr: %#v\ngot:  %#v\nwant: %#v", e, got, want)
		return false
	}
	return true
}
