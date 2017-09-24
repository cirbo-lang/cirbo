package globals

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/source"
)

//go:generate go run generate/generate_types_gen.go

var String = cty.TypeTypeVal(cty.String)
var Bool = cty.TypeTypeVal(cty.Bool)
var Type = cty.TypeTypeVal(cty.TypeType)
var Object cty.Value

func init() {
	Object = cty.FunctionVal(cty.FunctionImpl{
		Signature: &cty.CallSignature{
			Parameters:           map[string]cty.CallParameter{},
			Result:               cty.TypeType,
			AcceptsVariadicNamed: true,
		},
		Callback: func(args cty.CallArgs) (cty.Value, source.Diags) {
			atys := map[string]cty.Type{}
			var diags source.Diags

			for n, v := range args.NamedVariadic {
				if !v.Type().Same(cty.TypeType) {
					diags = append(diags, source.Diag{
						Level:   source.Error,
						Summary: "Invalid object attribute type",
						Detail:  fmt.Sprintf("Attribute %q defined with a value of type %s. A Type is required.", n, v.Type()),
						// Ranges will be set on the way out through the evaluator, if possible.
					})
				}

				aty := v.UnwrapType()
				atys[n] = aty
			}

			return cty.TypeTypeVal(cty.Object(atys)), diags
		},
	})
}
