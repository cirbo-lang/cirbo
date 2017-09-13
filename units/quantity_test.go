package units

import (
	"fmt"
	"math/big"
	"testing"
)

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
