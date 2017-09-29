package units

import (
	"fmt"
	"testing"
)

func TestUnitCommensurableWith(t *testing.T) {
	tests := []struct {
		A    string
		B    string
		Want bool
	}{
		{"<nil>", "<nil>", true},
		{"<nil>", "m", false},
		{"m", "<nil>", false},
		{"", "", true},

		{"kg", "kg", true},
		{"kg", "lb", true},

		{"m", "m", true},
		{"m", "mm", true},
		{"mm", "in", true},

		{"deg", "deg", true},
		{"deg", "rad", true},
		{"rad", "deg", true},

		{"s", "s", true},
		{"s", "ms", true},

		{"A", "A", true},

		{"cd", "cd", true},

		{"ohm", "kohm", true},

		{"V", "kV", true},

		{"Hz", "MHz", true},

		{"kg", "m", false},
		{"V", "W", false},
		{"N", "W", false},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s to %s", test.A, test.B)
		t.Run(name, func(t *testing.T) {
			a := unitByName[test.A]
			b := unitByName[test.B]

			if a == nil && test.A != "<nil>" {
				t.Fatalf("no unit named %q", test.A)
			}
			if b == nil && test.B != "<nil>" {
				t.Fatalf("no unit named %q", test.B)
			}
			got := a.CommensurableWith(b)
			if got != test.Want {
				t.Errorf(
					"wrong result\nA: %#v\nB: %#v\ngot:  %#v\nwant: %#v",
					a, b, got, test.Want,
				)
			}
		})
	}
}

func TestUnitReciprocal(t *testing.T) {
	tests := []struct {
		Input *Unit
		Want  string
	}{
		{
			dimless,
			"",
		},
		{
			unitByName["kg"],
			"kg⁻¹",
		},
		{
			unitByName["m"],
			"m⁻¹",
		},
		{
			unitByName["in"],
			"in⁻¹",
		},
		{
			unitByName["deg"],
			"deg⁻¹",
		},
		{
			unitByName["s"],
			"Hz",
		},
		{
			unitByName["Hz"],
			"s",
		},
		{
			unitByName["A"],
			"A⁻¹",
		},
		{
			unitByName["cd"],
			"cd⁻¹",
		},
		{
			unitByName["m"].Multiply(unitByName["s"]),
			"m⁻¹ s⁻¹",
		},
		{
			unitByName["V"],
			"kg⁻¹ m⁻² A s³",
		},
		{
			unitByName["lx"],
			"m² cd⁻¹",
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%#v", test.Input)
		t.Run(name, func(t *testing.T) {
			got := test.Input.Reciprocal().String()

			if got != test.Want {
				t.Errorf(
					"wrong result\ninput: %#v\ngot:  %s\nwant: %s",
					test.Input, got, test.Want,
				)
			}
		})
	}

}

func TestUnitString(t *testing.T) {
	tests := []struct {
		Input *Unit
		Want  string
	}{
		{
			dimless,
			"",
		},
		{
			unitByName["kg"],
			"kg",
		},
		{
			unitByName["m"],
			"m",
		},
		{
			unitByName["V"],
			"V",
		},
		{
			unitByName["kohm"],
			"kohm",
		},
		{
			&Unit{
				Dimensionality{Mass: 1, Time: 1},
				baseUnits{Mass: kilogram, Time: second},
				0,
			},
			"kg s",
		},
		{
			&Unit{
				Dimensionality{Length: 1, Time: -2},
				baseUnits{Length: meter, Time: second},
				0,
			},
			"m s⁻²",
		},
		{
			&Unit{
				Dimensionality{Angle: 1, Time: -2},
				baseUnits{Angle: degree, Time: second},
				0,
			},
			"deg s⁻²",
		},
		{
			&Unit{
				Dimensionality{Length: 1, Time: -1},
				baseUnits{Length: inch, Time: microsecond},
				0,
			},
			"in us⁻¹",
		},
		{
			&Unit{
				Dimensionality{Length: 2},
				baseUnits{Length: centimeter},
				0,
			},
			"cm²",
		},
		{
			&Unit{
				Dimensionality{ElectricCurrent: 12},
				baseUnits{ElectricCurrent: ampere},
				0,
			},
			"A¹²",
		},
		{
			&Unit{
				Dimensionality{LuminousIntensity: -20},
				baseUnits{LuminousIntensity: candela},
				0,
			},
			"cd⁻²⁰",
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%#v", test.Input)
		t.Run(name, func(t *testing.T) {
			got := test.Input.String()

			if got != test.Want {
				t.Errorf(
					"wrong result\ninput: %#v\ngot:  %s\nwant: %s",
					test.Input, got, test.Want,
				)
			}
		})
	}
}
