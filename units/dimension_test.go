package units

import (
	"fmt"
	"testing"
)

func TestDimensionalityString(t *testing.T) {
	tests := []struct {
		Input Dimensionality
		Want  string
	}{
		{
			Dimensionality{},
			"",
		},
		{
			Dimensionality{
				Length: 1,
			},
			"[L]",
		},
		{
			Dimensionality{
				Length: 2,
			},
			"[L]²",
		},
		{
			Dimensionality{
				Length: -1,
			},
			"[L]⁻¹",
		},
		{
			Dimensionality{
				Mass:   1,
				Length: 1,
				Time:   -1,
			},
			"[M][L][T]⁻¹",
		},
		{
			Dimensionality{
				Length: 1,
				Time:   -2,
			},
			"[L][T]⁻²",
		},
		{
			Dimensionality{
				Angle: 1,
				Time:  -2,
			},
			"[angle][T]⁻²",
		},
		{
			Dimensionality{
				ElectricCurrent: 2,
			},
			"[I]²",
		},
		{
			Dimensionality{
				LuminousIntensity: 1234567890,
			},
			"[J]¹²³⁴⁵⁶⁷⁸⁹⁰",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			got := test.Input.String()
			if got != test.Want {
				t.Errorf("wrong result\ninput: %#v\ngot:   %s\nwant:  %s", test.Input, got, test.Want)
			}
		})
	}
}
