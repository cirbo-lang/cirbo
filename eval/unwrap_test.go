package eval

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/cbty"
)

func TestUnwrap(t *testing.T) {
	tests := []struct {
		Val  cbty.Value
		Want Unwrapped
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
	}

	for _, test := range tests {
		t.Run(test.Val.GoString(), func(t *testing.T) {
			got := Unwrap(test.Val)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
