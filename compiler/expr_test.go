package compiler

import (
	"testing"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/parser"
	"github.com/cirbo-lang/cirbo/source"
	"github.com/cirbo-lang/cirbo/units"
)

func TestCompileExpr(t *testing.T) {
	tests := []struct {
		Source string
		Want   cbty.Value
		Diags  int
	}{
		{
			"true",
			cbty.True,
			0,
		},
		{
			"(true)",
			cbty.True,
			0,
		},
		{
			"0",
			cbty.Zero,
			0,
		},
		{
			"50%",
			cbty.NumberValFloat(0.5),
			0,
		},
		{
			"5mm",
			cbty.QuantityVal(units.MakeQuantityInt(5, units.ByName("mm"))),
			0,
		},
		{
			"foo",
			cbty.StringVal("foo"),
			0,
		},
		{
			"bar",
			cbty.StringVal("bar"),
			0,
		},
		{
			"baz",
			cbty.PlaceholderVal,
			1, // variable not declared in this scope
		},
		{
			"notDefined",
			cbty.PlaceholderVal,
			1, // symbol not defined (if this happens then it's a Cirbo bug rather than user error, but we want to still catch it gracefully)
		},
		{
			"1 + 2",
			cbty.NumberValInt(3),
			0,
		},
		{
			"1m + 50cm",
			cbty.QuantityVal(units.MakeQuantityFloat(1.5, units.ByName("m"))),
			0,
		},
		{
			"4 - 1",
			cbty.NumberValInt(3),
			0,
		},
		{
			"2 * 6",
			cbty.NumberValInt(12),
			0,
		},
		{
			"16 / 2",
			cbty.NumberValInt(8),
			0,
		},
		{
			"true == true",
			cbty.True,
			0,
		},
		{
			"true == false",
			cbty.False,
			0,
		},
		{
			"true != false",
			cbty.True,
			0,
		},
		{
			"true || false",
			cbty.True,
			0,
		},
		{
			"true && false",
			cbty.False,
			0,
		},
		{
			"!true",
			cbty.False,
			0,
		},
		{
			"blah blah",
			cbty.PlaceholderVal,
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
	ctx.DefineLiteral(fooSym, cbty.StringVal("foo"))
	ctx.DefineLiteral(barSym, cbty.StringVal("bar"))
	ctx.DefineLiteral(upperSym, cbty.FunctionVal(cbty.FunctionImpl{
		Signature: &cbty.CallSignature{
			Parameters: map[string]cbty.CallParameter{
				"str": {
					Type:     cbty.String,
					Required: true,
				},
			},
			Result: cbty.String,
		},
		Callback: func(args cbty.CallArgs) (cbty.Value, source.Diags) {
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

			val, evalDiags := expr.Value(ctx)
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
