package cty

import (
	"fmt"
	"testing"
)

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
	}

	for _, test := range tests {
		as := test.A.v.(string)
		bs := test.B.v.(string)
		t.Run(fmt.Sprintf("%q .. %q", as, bs), func(t *testing.T) {
			gotVal := test.A.Concat(test.B)
			if got, want := gotVal.Type(), test.Want.Type(); !got.Same(want) {
				t.Fatalf("got %#v value; want %#v", got, want)
			}
			got := gotVal.v.(string)
			want := test.Want.v.(string)
			if got != want {
				t.Errorf("wrong result\ngot:  %s\nwant: %s", got, want)
			}
		})
	}
}
