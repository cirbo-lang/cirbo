package cbty

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

func TestBoolNot(t *testing.T) {
	tests := []struct {
		V    Value
		Want Value
	}{
		{
			False,
			True,
		},
		{
			True,
			False,
		},
		{
			UnknownVal(Bool),
			UnknownVal(Bool),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Not()", test.V), func(t *testing.T) {
			got := test.V.Not()
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestBoolAnd(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			False,
			False,
			False,
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
		t.Run(fmt.Sprintf("%#v.And(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.And(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestBoolOr(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			False,
			False,
			False,
		},
		{
			False,
			True,
			True,
		},
		{
			True,
			False,
			True,
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
		t.Run(fmt.Sprintf("%#v.Or(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Or(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
