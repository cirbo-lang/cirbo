package cty

import (
	"fmt"
	"testing"
)

func TestBoolEqual(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			False,
			False,
			True,
		},
		{
			False,
			True,
			False,
		},
		{
			True,
			False,
			False,
		},
		{
			True,
			True,
			True,
		},
		{
			UnknownVal(Bool),
			True,
			UnknownVal(Bool),
		},
		{
			UnknownVal(Bool),
			False,
			UnknownVal(Bool),
		},
		{
			UnknownVal(Bool),
			UnknownVal(Bool),
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
