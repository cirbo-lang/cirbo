package eval

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/cbo"
	"github.com/cirbo-lang/cirbo/cbty"
)

func TestUnwrap(t *testing.T) {
	tests := []struct {
		Val  cbty.Value
		Want cbo.Any
	}{
		{
			cbty.PlaceholderVal,
			nil,
		},
		{
			cbty.UnknownVal(cbty.Bool),
			nil,
		},
		{
			cbty.True,
			true,
		},
		{
			cbty.False,
			false,
		},
		{
			cbty.StringVal("hello"),
			"hello",
		},
		{
			cbty.TypeTypeVal(cbty.String),
			cbty.String,
		},
		{
			deviceValue(&device{
				name: "Fred",
				attrs: StmtBlockAttrs{
					"required": {
						Type: cbty.String,
					},
					"optional": {
						Type:    cbty.Bool,
						Default: cbty.False,
					},
				},
			}),
			&cbo.Device{
				Name: "Fred",
				Attrs: cbo.AttributesDef{
					"required": {
						Type:     cbty.String,
						Required: true,
					},
					"optional": {
						Type:     cbty.Bool,
						Required: false,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Val.GoString(), func(t *testing.T) {
			unwr := &Unwrapper{}
			got := unwr.Unwrap(test.Val)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
