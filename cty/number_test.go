package cty

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/cirbo-lang/cirbo/units"
)

func TestAddNumber(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			testNumber("2", ""),
			testNumber("3", ""),
			testNumber("5", ""),
		},
		{
			testNumber("20", "cm"),
			testNumber("3", "m"),
			testNumber("3.2", "m"),
		},
	}

	for _, test := range tests {
		aq := test.A.v.(units.Quantity)
		bq := test.B.v.(units.Quantity)
		t.Run(fmt.Sprintf("%s + %s", aq, bq), func(t *testing.T) {
			gotVal := test.A.Add(test.B)
			if got, want := gotVal.Type(), test.A.Type(); !got.Same(want) {
				t.Fatalf("got %#v value; want %#v", got, want)
			}
			got := gotVal.v.(units.Quantity)
			want := test.Want.v.(units.Quantity)
			if !got.Same(want) {
				t.Errorf("wrong result\ngot:  %s\nwant: %s", got, want)
			}
		})
	}
}

func testNumber(n, u string) Value {
	num, _, err := (&big.Float{}).Parse(n, 10)
	if err != nil {
		panic(err)
	}

	unit := units.ByName(u)
	if unit == nil {
		panic(fmt.Errorf("no unit named %q", u))
	}

	q := units.MakeQuantity(num, unit)
	return NumberVal(q)
}
