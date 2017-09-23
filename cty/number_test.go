package cty

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/cirbo-lang/cirbo/units"
)

func TestNumberEqual(t *testing.T) {
	tests := []struct {
		A, B Value
		Want Value
	}{
		{
			testNumber("2", ""),
			testNumber("2", ""),
			True,
		},
		{
			testNumber("2", ""),
			testNumber("3", ""),
			False,
		},
		{
			testNumber("2", ""),
			testNumber("2", "cm"),
			False,
		},
		{
			testNumber("20", "mm"),
			testNumber("2", "cm"),
			True,
		},
		{
			UnknownVal(Length),
			testNumber("3", "m"),
			UnknownVal(Bool),
		},
		{
			testNumber("3", "m"),
			UnknownVal(Length),
			UnknownVal(Bool),
		},
		{
			UnknownVal(Length),
			UnknownVal(Length),
			UnknownVal(Bool),
		},
		{
			UnknownVal(Area),
			UnknownVal(Length),
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
		{
			UnknownVal(Length),
			testNumber("3", "m"),
			UnknownVal(Length),
		},
		{
			UnknownVal(Length),
			testNumber("3", "m"),
			UnknownVal(Length),
		},
		{
			testNumber("3", "m"),
			UnknownVal(Length),
			UnknownVal(Length),
		},
		{
			UnknownVal(Length),
			UnknownVal(Length),
			UnknownVal(Length),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Add(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Add(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
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
		{
			UnknownVal(Length),
			testNumber("3", "m"),
			UnknownVal(Length),
		},
		{
			testNumber("3", "m"),
			UnknownVal(Length),
			UnknownVal(Length),
		},
		{
			UnknownVal(Length),
			UnknownVal(Length),
			UnknownVal(Length),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Subtract(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Subtract(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
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
		{
			UnknownVal(Current),
			UnknownVal(Resistance),
			UnknownVal(Voltage),
		},
		{
			testNumber("15", "A"),
			UnknownVal(Resistance),
			UnknownVal(Voltage),
		},
		{
			UnknownVal(Resistance),
			testNumber("15", "A"),
			UnknownVal(Voltage),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Multiply(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Multiply(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
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
		{
			UnknownVal(Length),
			UnknownVal(Time),
			UnknownVal(Speed),
		},
		{
			testNumber("1", "m"),
			UnknownVal(Time),
			UnknownVal(Speed),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v.Divide(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Divide(test.B)
			if !got.Same(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
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
	return QuantityVal(q)
}

func testNumberU(n string, unit *units.Unit) Value {
	num, _, err := (&big.Float{}).Parse(n, 10)
	if err != nil {
		panic(err)
	}

	q := units.MakeQuantity(num, unit)
	return QuantityVal(q)
}
