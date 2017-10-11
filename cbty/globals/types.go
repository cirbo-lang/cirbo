package globals

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

//go:generate go run generate/generate_types_gen.go

var String = cbty.TypeTypeVal(cbty.String)
var Bool = cbty.TypeTypeVal(cbty.Bool)
var Type = cbty.TypeTypeVal(cbty.TypeType)
var Object cbty.Value

func init() {
	Object = cbty.FunctionVal(cbty.FunctionImpl{
		Signature: &cbty.CallSignature{
			Parameters:           map[string]cbty.CallParameter{},
			Result:               cbty.TypeType,
			AcceptsVariadicNamed: true,
		},
		Callback: func(args cbty.CallArgs) (cbty.Value, source.Diags) {
			atys := map[string]cbty.Type{}
			var diags source.Diags

			for n, v := range args.NamedVariadic {
				if !v.Type().Same(cbty.TypeType) {
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

			return cbty.TypeTypeVal(cbty.Object(atys)), diags
		},
	})
}
