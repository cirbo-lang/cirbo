package cty

import (
	"fmt"
	"testing"
)

func TestCallSignatureSame(t *testing.T) {
	tests := []struct {
		A, B *CallSignature
		Want bool
	}{
		{
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,
			},
			true,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Result: String,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Result: String,
			},
			true,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Positional: []string{"enabled"},
				Result:     String,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Positional: []string{"enabled"},
				Result:     String,
			},
			true,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Positional: []string{"enabled"},
				Result:     String,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Result: String,
			},
			false,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Result: String,
			},
			false,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Result: String,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: String,
					},
				},
				Result: String,
			},
			false,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Result: String,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{
					"enabled": {
						Type: Bool,
					},
				},
				Result: Bool,
			},
			false,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicNamed: true,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicNamed: true,
			},
			true,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicNamed: false,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicNamed: true,
			},
			false,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicPositional: true,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicPositional: true,
			},
			true,
		},
		{
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicPositional: false,
			},
			&CallSignature{
				Parameters: map[string]CallParameter{},
				Result:     String,

				AcceptsVariadicPositional: true,
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("(%#v).Same(%#v)", test.A, test.B), func(t *testing.T) {
			got := test.A.Same(test.B)
			if got != test.Want {
				t.Logf("A: %#v", test.A)
				t.Logf("B: %#v", test.B)
				t.Errorf("wrong result %#v; want %#v", got, test.Want)
			}
		})
	}
}
