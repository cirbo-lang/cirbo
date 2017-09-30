package eval

import (
	"testing"

	"github.com/cirbo-lang/cirbo/cty"
)

func TestGlobalScope(t *testing.T) {
	tests := []struct {
		Name   string
		Wanted bool
	}{
		{
			"",
			false,
		},
		{
			"BlahBlahBaz",
			false,
		},

		// We do not exhaustively test all of the symbols in our global table
		// because that would just redundantly re-define the table from the
		// globals package, but we test a few just to verify that the scope-
		// building mechanism is working.
		{
			"Length",
			true,
		},
		{
			"Area",
			true,
		},
		{
			"String",
			true,
		},
		{
			"Object",
			true,
		},
		{
			"Type",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			declared := GlobalScope().Declared(test.Name)
			symbol := GlobalScope().Get(test.Name)

			if test.Wanted {
				if !declared {
					t.Errorf("symbol %q not declared; should have been", test.Name)
				}
				if symbol == nil {
					t.Errorf("symbol %q is nil; should be non-nil", test.Name)
				}
			} else {
				if declared {
					t.Errorf("symbol %q is declared; should not have been", test.Name)
				}
				if symbol != nil {
					t.Errorf("symbol %q is non-nil; should be nil", test.Name)
				}
			}
		})
	}
}

func TestGlobalContext(t *testing.T) {
	tests := []struct {
		Name     string
		WantType cty.Type
	}{
		{
			"",
			cty.NilType,
		},
		{
			"BlahBlahBaz",
			cty.NilType,
		},

		// We do not exhaustively test all of the symbols in our global table
		// because that would just redundantly re-define the table from the
		// globals package, but we test a few just to verify that the scope-
		// building mechanism is working.
		{
			"Length",
			cty.TypeType,
		},
		{
			"Area",
			cty.TypeType,
		},
		{
			"String",
			cty.TypeType,
		},
		{
			"Object",
			cty.Function(&cty.CallSignature{
				Parameters:           map[string]cty.CallParameter{},
				AcceptsVariadicNamed: true,
				Result:               cty.TypeType,
			}),
		},
		{
			"Type",
			cty.TypeType,
		},
	}

	scope := GlobalScope()
	ctx := GlobalContext()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			symbol := scope.Get(test.Name)
			if symbol == nil {
				if test.WantType == cty.NilType {
					return
				}
				t.Fatalf("symbol %q not declared; should have been", test.Name)
			}

			defined := ctx.Defined(symbol)
			value := ctx.Value(symbol)

			if test.WantType != cty.NilType {
				if !defined {
					t.Errorf("symbol %q not defined; should have been", test.Name)
				}
				if !value.Type().Same(test.WantType) {
					t.Errorf("symbol %q is %#v; should be %#v", test.Name, value.Type(), test.WantType)
				}
			} else {
				if defined {
					t.Errorf("symbol %q is defined; should not have been", test.Name)
				}
				if value.Type() == cty.NilType {
					t.Errorf("symbol %q is %#v; should be %#v", test.Name, value, test.WantType)
				}
			}
		})
	}
}
