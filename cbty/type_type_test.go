package cbty

import (
	"fmt"
	"testing"
)

func TestTypeTypeEqual(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			TypeTypeVal(Bool),
			TypeTypeVal(Bool),
			True,
		},
		{
			TypeTypeVal(Bool),
			TypeTypeVal(String),
			False,
		},
		{
			TypeTypeVal(TypeType),
			TypeTypeVal(TypeType),
			True,
		},
		{
			TypeTypeVal(TypeType),
			TypeTypeVal(TypeType).TypeValue(),
			True,
		},
		{
			UnknownVal(TypeType),
			TypeTypeVal(Bool),
			UnknownVal(Bool),
		},
		{
			UnknownVal(TypeType),
			UnknownVal(TypeType),
			UnknownVal(Bool),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Equal(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Equal(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
