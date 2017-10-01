package compiler

import (
	"testing"

	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/parser"
	"github.com/cirbo-lang/cirbo/source"
	"github.com/cirbo-lang/cirbo/units"
)

func TestCompileExpr(t *testing.T) {
	tests := []struct {
		Source string
		Want   cty.Value
		Diags  int
	}{
		{
			"true",
			cty.True,
			0,
		},
		{
			"(true)",
			cty.True,
			0,
		},
		{
			"0",
			cty.Zero,
			0,
		},
		{
			"50%",
			cty.NumberValFloat(0.5),
			0,
		},
		{
			"5mm",
			cty.QuantityVal(units.MakeQuantityInt(5, units.ByName("mm"))),
			0,
		},
		{
			"foo",
			cty.StringVal("foo"),
			0,
		},
		{
			"bar",
			cty.StringVal("bar"),
			0,
		},
		{
			"baz",
			cty.PlaceholderVal,
			1, // variable not declared in this scope
		},
		{
			"notDefined",
			cty.PlaceholderVal,
			1, // symbol not defined (if this happens then it's a Cirbo bug rather than user error, but we want to still catch it gracefully)
		},
		{
			"1 + 2",
			cty.NumberValInt(3),
			0,
		},
		{
			"1m + 50cm",
			cty.QuantityVal(units.MakeQuantityFloat(1.5, units.ByName("m"))),
			0,
		},
		{
			"4 - 1",
			cty.NumberValInt(3),
			0,
		},
		{
			"2 * 6",
			cty.NumberValInt(12),
			0,
		},
		{
			"16 / 2",
			cty.NumberValInt(8),
			0,
		},
		{
			"true == true",
			cty.True,
			0,
		},
		{
			"true == false",
			cty.False,
			0,
		},
		{
			"true != false",
			cty.True,
			0,
		},
		{
			"true || false",
			cty.True,
			0,
		},
		{
			"true && false",
			cty.False,
			0,
		},
		{
			"!true",
			cty.False,
			0,
		},
		{
			"blah blah",
			cty.PlaceholderVal,
			2, // unknown variable "blah", unexpected characters after expression
		},
	}

	scope1 := eval.GlobalScope().NewChild()
	fooSym := scope1.Declare("foo")
	scope2 := scope1.NewChild()
	barSym := scope2.Declare("bar")
	upperSym := scope1.Declare("upper")
	scope1.Declare("notDefined")

	ctx := eval.GlobalContext().NewChild()
	ctx.Define(fooSym, cty.StringVal("foo"))
	ctx.Define(barSym, cty.StringVal("bar"))
	ctx.Define(upperSym, cty.FunctionVal(cty.FunctionImpl{
		Signature: &cty.CallSignature{
			Parameters: map[string]cty.CallParameter{
				"str": {
					Type:     cty.String,
					Required: true,
				},
			},
			Result: cty.String,
		},
		Callback: func(args cty.CallArgs) (cty.Value, source.Diags) {
			// TODO: Implement this and test with it below once the evaluator
			// knows how to handle calls.
			panic("upper function not yet implemented")
		},
	}))

	for _, test := range tests {
		t.Run(test.Source, func(t *testing.T) {
			var diags source.Diags

			node, parseDiags := parser.ParseExpr([]byte(test.Source))
			diags = append(diags, parseDiags...)

			expr, compileDiags := CompileExpr(node, scope2)
			diags = append(diags, compileDiags...)

			val, evalDiags := eval.Value(expr, ctx)
			diags = append(diags, evalDiags...)

			if len(diags) != test.Diags {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.Diags)
				for _, diag := range diags {
					t.Logf("- %s", diag)
				}
			}

			if !val.Same(test.Want) {
				t.Errorf("wrong result\nsrc:  %s\ngot:  %#v\nwant: %#v", test.Source, val, test.Want)
			}
		})
	}
}
