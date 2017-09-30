package eval

import (
	"fmt"
	"testing"

	"github.com/cirbo-lang/cirbo/source"

	"github.com/cirbo-lang/cirbo/cty"
)

func TestExprImpls(t *testing.T) {
	// All of the testing here actually happens at compile time. We dress it
	// up like run-time tests just because that makes this test visible in
	// the test results when we pass.
	tests := []Expr{
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
