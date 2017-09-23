package cty

import (
	"fmt"
	"testing"
)

func TestObjectTypeName(t *testing.T) {
	tests := []struct {
		Type Type
		Want string
	}{
		{
			EmptyObject,
			"Object()",
		},
		{
			Object(map[string]Type{"foo": Bool}),
			"Object(foo=Bool)",
		},
		{
			Object(map[string]Type{"foo": Bool, "bar": String}),
			"Object(bar=String, foo=Bool)",
		},
		{
			Object(map[string]Type{"foo": Bool, "bar": EmptyObject}),
			"Object(bar=Object(), foo=Bool)",
		},
		{
			Object(map[string]Type{"foo": Bool, "bar": Object(map[string]Type{"baz": Number})}),
			"Object(bar=Object(baz=Number), foo=Bool)",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Name()", test.Type), func(t *testing.T) {
			got := test.Type.Name()
			if got != test.Want {
				t.Errorf("wrong result\ntype: %#v\ngot:  %#v\nwant: %#v", test.Type, got, test.Want)
			}
		})
	}

}

func TestObjectEqual(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			EmptyObjectVal,
			EmptyObjectVal,
			True,
		},
		{
			EmptyObjectVal,
			ObjectVal(map[string]Value{"foo": True}),
			False,
		},
		{
			ObjectVal(map[string]Value{"foo": True}),
			ObjectVal(map[string]Value{"foo": True}),
			True,
		},
		{
			UnknownVal(EmptyObject),
			ObjectVal(map[string]Value{"foo": True}),
			False,
		},
		{
			UnknownVal(EmptyObject),
			EmptyObjectVal,
			UnknownVal(Bool),
		},
		{
			ObjectVal(map[string]Value{"foo": UnknownVal(Bool)}),
			ObjectVal(map[string]Value{"foo": True}),
			UnknownVal(Bool),
		},
		{
			ObjectVal(map[string]Value{"foo": UnknownVal(String)}),
			ObjectVal(map[string]Value{"foo": True}),
			False,
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
