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
