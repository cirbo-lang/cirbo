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
			if got, want := gotVal.Type(), test.Want.Type(); !got.Same(want) {
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

func TestSubtractNumber(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			testNumber("2", ""),
			testNumber("3", ""),
			testNumber("-1", ""),
		},
		{
			testNumber("3", "m"),
			testNumber("20", "cm"),
			testNumber("2.8", "m"),
		},
	}

	for _, test := range tests {
		aq := test.A.v.(units.Quantity)
		bq := test.B.v.(units.Quantity)
		t.Run(fmt.Sprintf("%s - %s", aq, bq), func(t *testing.T) {
			gotVal := test.A.Subtract(test.B)
			if got, want := gotVal.Type(), test.Want.Type(); !got.Same(want) {
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

func TestMultiplyNumber(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			testNumber("2", ""),
			testNumber("3", ""),
			testNumber("6", ""),
		},
		{
			testNumber("3", "m"),
			testNumber("2", ""),
			testNumber("6", "m"),
		},
		{
			testNumber("3", "m"),
			testNumber("2", "m"),
			testNumberU("6", units.Meter.ToPower(2)),
		},
		{
			testNumber("3", "m"),
			testNumberU("2", units.Meter.ToPower(2)),
			testNumberU("6", units.Meter.ToPower(3)),
		},
		{
			testNumber("3", "A"),
			testNumber("2", "ohm"),
			testNumber("6", "V"),
		},
		{
			testNumber("3", "A"),
			testNumber("2", "V"),
			testNumber("6", "W"),
		},
		{
			testNumberU("3", units.Meter.Multiply(units.Second.ToPower(-1))),
			testNumberU("10", units.Second),
			testNumber("30", "m"),
		},
	}

	for _, test := range tests {
		aq := test.A.v.(units.Quantity)
		bq := test.B.v.(units.Quantity)
		t.Run(fmt.Sprintf("%s × %s", aq, bq), func(t *testing.T) {
			gotVal := test.A.Multiply(test.B)
			if got, want := gotVal.Type(), test.Want.Type(); !got.Same(want) {
				t.Fatalf("wrong result type\ngot:  %#v\nwant: %#v", got, want)
			}
			got := gotVal.v.(units.Quantity)
			want := test.Want.v.(units.Quantity)
			if !got.Same(want) {
				t.Errorf("wrong result\ngot:  %s\nwant: %s", got, want)
			}
		})
	}
}

func TestDivideNumber(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			testNumber("3", ""),
			testNumber("2", ""),
			testNumber("1.5", ""),
		},
		{
			testNumber("18", "m"),
			testNumber("3", ""),
			testNumber("6", "m"),
		},
		{
			testNumber("3", "m"),
			testNumber("2", "m"),
			testNumber("1.5", ""),
		},
		{
			testNumber("3", "m"),
			testNumberU("2", units.Meter.ToPower(2)),
			testNumberU("1.5", units.Meter.ToPower(-1)),
		},
		{
			testNumber("3", "V"),
			testNumber("2", "ohm"),
			testNumber("1.5", "A"),
		},
		{
			testNumber("3", "W"),
			testNumber("2", "V"),
			testNumber("1.5", "A"),
		},
		{
			testNumber("3", "m"),
			testNumber("10", "s"),
			testNumberU("0.3", units.Meter.Multiply(units.Second.ToPower(-1))),
		},
	}

	for _, test := range tests {
		aq := test.A.v.(units.Quantity)
		bq := test.B.v.(units.Quantity)
		t.Run(fmt.Sprintf("%s ÷ %s", aq, bq), func(t *testing.T) {
			gotVal := test.A.Divide(test.B)
			if got, want := gotVal.Type(), test.Want.Type(); !got.Same(want) {
				t.Fatalf("wrong result type\ngot:  %#v\nwant: %#v", got, want)
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

func testNumberU(n string, unit *units.Unit) Value {
	num, _, err := (&big.Float{}).Parse(n, 10)
	if err != nil {
		panic(err)
	}

	q := units.MakeQuantity(num, unit)
	return NumberVal(q)
}
