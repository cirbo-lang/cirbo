package units

import (
	"fmt"
	"math/big"
	"testing"
)

func TestQuantityConvert(t *testing.T) {
	tests := []struct {
		Q    Quantity
		U    *Unit
		Want string
	}{
		{
			q("1", unitByName["lb"]),
			unitByName["kg"],
			"0.45359237 kg",
		},
		{
			q("1", unitByName["lb"]),
			unitByName["st"],
			"14 st",
		},
		{
			q("1", unitByName["m"]),
			unitByName["cm"],
			"100 cm",
		},
		{
			q("1", unitByName["in"]),
			unitByName["cm"],
			"2.54 cm",
		},
		{
			q("1", unitByName["in"]),
			unitByName["mil"],
			"1000 mil",
		},
		{
			q("1", unitByName["ft"]),
			unitByName["in"],
			"12 in",
		},
		{
			q("1", unitByName["yd"]),
			unitByName["in"],
			"36 in",
		},
		{
			q("1", &Unit{Dimensionality{Length: 2}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: 2}, baseUnits{Length: centimeter}, 0},
			"10000 cm²",
		},
		{
			q("1", &Unit{Dimensionality{Length: 3}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: 3}, baseUnits{Length: centimeter}, 0},
			"1000000 cm³",
		},
		{
			q("1", &Unit{Dimensionality{Length: -2}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: -2}, baseUnits{Length: centimeter}, 0},
			"0.0001 cm⁻²",
		},
		{
			q("1", &Unit{Dimensionality{Length: -3}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: -3}, baseUnits{Length: centimeter}, 0},
			"1e-06 cm⁻³",
		},
		{
			q("1", unitByName["MHz"]),
			unitByName["Hz"],
			"1000000 Hz",
		},
		{
			q("1", unitByName["ohm"]),
			unitByName["kohm"],
			"0.001 kohm",
		},
		{
			q("1", unitByName["kohm"]),
			unitByName["ohm"],
			"1000 ohm",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s to %s", test.Q, test.U), func(t *testing.T) {
			got := test.Q.Convert(test.U)
			gotStr := got.String()
			if gotStr != test.Want {
				t.Errorf("wrong result\nquant: %s\nunit:  %s\ngot:   %s\nwant:  %s", test.Q, test.U, gotStr, test.Want)
			}
		})
	}
}

func TestQuantityString(t *testing.T) {
	tests := []struct {
		Input Quantity
		Want  string
	}{
		{
			q("1", unitByName["cm"]),
			"1 cm",
		},
		{
			q("1", nil),
			"1",
		},
		{
			q("2.89", &Unit{
				Dimensionality{Length: 1, Time: -2},
				baseUnits{Length: meter, Time: second},
				0,
			}),
			"2.89 m s⁻²",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Input), func(t *testing.T) {
			got := test.Input.String()
			if got != test.Want {
				t.Errorf("wrong result\ninput: %#v\ngot:   %s\nwant:  %s", test.Input, got, test.Want)
			}
		})
	}
}

func q(v string, u *Unit) Quantity {
	f, _, err := (&big.Float{}).Parse(v, 10)
	if err != nil {
		panic(err)
	}

	return MakeQuantity(f, u)
}
