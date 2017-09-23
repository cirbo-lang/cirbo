package cty

import (
	"fmt"
	"testing"
)

func TestStringEqual(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			StringVal(""),
			StringVal(""),
			True,
		},
		{
			StringVal("hello"),
			StringVal("hello"),
			True,
		},
		{
			StringVal("hello"),
			StringVal("world"),
			False,
		},
		{
			UnknownVal(String),
			StringVal("hello"),
			UnknownVal(Bool),
		},
		{
			StringVal("hello"),
			UnknownVal(String),
			UnknownVal(Bool),
		},
		{
			UnknownVal(String),
			UnknownVal(String),
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

func TestConcatString(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			StringVal(""),
			StringVal(""),
			StringVal(""),
		},
		{
			StringVal("a"),
			StringVal(""),
			StringVal("a"),
		},
		{
			StringVal("a"),
			StringVal("b"),
			StringVal("ab"),
		},
		{
			UnknownVal(String),
			StringVal("b"),
			UnknownVal(String),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Concat(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Concat(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
