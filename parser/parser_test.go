package parser

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/source"
	"github.com/kylelemons/godebug/pretty"
)

func TestParseTopLevel(t *testing.T) {
	tests := []struct {
		Input     string
		Want      []ast.Node
		DiagCount int
	}{
		{
			"",
			nil,
			0,
		},
		{
			"    ",
			nil,
			0,
		},
		{
			"\n\n\n\n",
			nil,
			0,
		},

		{
			`import "baz";`,
			[]ast.Node{
				&ast.Import{
					Package: "baz",
					Name:    "",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
					PackageRange: source.Range{
						Start: source.Pos{Line: 1, Column: 8, Byte: 7},
						End:   source.Pos{Line: 1, Column: 13, Byte: 12},
					},
				},
			},
			0,
		},
		{
			`import "baz" as foo;`,
			[]ast.Node{
				&ast.Import{
					Package: "baz",
					Name:    "foo",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 21, Byte: 20},
						},
					},
					PackageRange: source.Range{
						Start: source.Pos{Line: 1, Column: 8, Byte: 7},
						End:   source.Pos{Line: 1, Column: 13, Byte: 12},
					},
				},
			},
			0,
		},
		{
			`import invalid;`,
			nil,
			1, // import path must be quoted string
		},
		{
			`import "valid1"; import invalid; import "valid2";`,
			[]ast.Node{
				&ast.Import{
					Package: "valid1",
					Name:    "",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 17, Byte: 16},
						},
					},
					PackageRange: source.Range{
						Start: source.Pos{Line: 1, Column: 8, Byte: 7},
						End:   source.Pos{Line: 1, Column: 16, Byte: 15},
					},
				},
				&ast.Import{
					Package: "valid2",
					Name:    "",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 34, Byte: 33},
							End:   source.Pos{Line: 1, Column: 50, Byte: 49},
						},
					},
					PackageRange: source.Range{
						Start: source.Pos{Line: 1, Column: 41, Byte: 40},
						End:   source.Pos{Line: 1, Column: 49, Byte: 48},
					},
				},
			},
			1, // import path must be quoted string
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			tokens := scanTokens([]byte(test.Input), "", source.StartPos, scanNormal)
			it := newTokenIterator(tokens)
			ip := &parser{
				tokenPeeker: tokenPeeker{
					Iter: it,
				},
			}
			got, _, diags := ip.ParseTopLevel()

			if len(diags) != test.DiagCount {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.DiagCount)
				for _, diag := range diags {
					t.Logf("- %s", diag.String())
				}
			}

			prettyConfig := &pretty.Config{
				Diffable:          true,
				IncludeUnexported: true,
				PrintStringers:    false,
			}

			if !reflect.DeepEqual(got, test.Want) {
				diff := prettyConfig.Compare(test.Want, got)
				t.Errorf("wrong result\ninput:\n%s\n\ndiff: %s", test.Input, diff)
			}
		})
	}
}
func TestParseExpression(t *testing.T) {
	tests := []struct {
		Input     string
		Want      ast.Node
		DiagCount int
	}{
		{
			"",
			&ast.Invalid{
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 1, Byte: 0},
					},
				},
			},
			1, // expected start of expression
		},
		{
			"    ",
			&ast.Invalid{
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 5, Byte: 4},
						End:   source.Pos{Line: 1, Column: 5, Byte: 4},
					},
				},
			},
			1, // expected start of expression
		},
		{
			"\n\n\n\n",
			&ast.Invalid{
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 5, Column: 1, Byte: 4},
						End:   source.Pos{Line: 5, Column: 1, Byte: 4},
					},
				},
			},
			1, // expected start of expression
		},

		{
			`"hello"`,
			&ast.StringLit{
				Value: "hello",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 8, Byte: 7},
					},
				},
			},
			0,
		},
		{
			`"he\nlo"`,
			&ast.StringLit{
				Value: "he\nlo",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 9, Byte: 8},
					},
				},
			},
			0,
		},
		{
			`"\q"`,
			&ast.StringLit{
				Value: "q",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 5, Byte: 4},
					},
				},
			},
			1, // invalid escape sequence
		},
		{
			`"hello" true`,
			&ast.StringLit{
				Value: "hello",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 8, Byte: 7},
					},
				},
			},
			1, // extra junk after expression
		},

		{
			`true`,
			&ast.BooleanLit{
				Value: true,
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 5, Byte: 4},
					},
				},
			},
			0,
		},
		{
			`false`,
			&ast.BooleanLit{
				Value: false,
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 6, Byte: 5},
					},
				},
			},
			0,
		},
		{
			`foo`,
			&ast.Variable{
				Name: "foo",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 4, Byte: 3},
					},
				},
			},
			0,
		},
		{
			"`foo`",
			&ast.Variable{
				Name: "foo",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 6, Byte: 5},
					},
				},
			},
			0,
		},
		{
			"`true`",
			&ast.Variable{
				Name: "true",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 7, Byte: 6},
					},
				},
			},
			0,
		},
		{
			"`false`",
			&ast.Variable{
				Name: "false",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 8, Byte: 7},
					},
				},
			},
			0,
		},

		{
			`1`,
			&ast.NumberLit{
				Value: mustParseBigFloat("1"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 2, Byte: 1},
					},
				},
			},
			0,
		},
		{
			`1.2`,
			&ast.NumberLit{
				Value: mustParseBigFloat("1.2"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 4, Byte: 3},
					},
				},
			},
			0,
		},
		{
			`1.0`,
			&ast.NumberLit{
				Value: mustParseBigFloat("1.0"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 4, Byte: 3},
					},
				},
			},
			0,
		},
		{
			`1%`,
			&ast.NumberLit{
				Value: mustParseBigFloat("0.01"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 3, Byte: 2},
					},
				},
			},
			0,
		},
		{
			`1.5%`,
			&ast.NumberLit{
				Value: mustParseBigFloat("0.015"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 5, Byte: 4},
					},
				},
			},
			0,
		},
		{
			`50%`,
			&ast.NumberLit{
				Value: mustParseBigFloat("0.5"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 4, Byte: 3},
					},
				},
			},
			0,
		},
		{
			`100%`,
			&ast.NumberLit{
				Value: mustParseBigFloat("1"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 5, Byte: 4},
					},
				},
			},
			0,
		},
		{
			`150%`,
			&ast.NumberLit{
				Value: mustParseBigFloat("1.5"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 5, Byte: 4},
					},
				},
			},
			0,
		},
		{
			`1m`,
			&ast.QuantityLit{
				Value: mustParseBigFloat("1"),
				Unit:  "m",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 3, Byte: 2},
					},
				},
			},
			0,
		},
		{
			`1kV`,
			&ast.QuantityLit{
				Value: mustParseBigFloat("1"),
				Unit:  "kV",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 4, Byte: 3},
					},
				},
			},
			0,
		},
		{
			`1 ohm`,
			&ast.QuantityLit{
				Value: mustParseBigFloat("1"),
				Unit:  "ohm",
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 6, Byte: 5},
					},
				},
			},
			0,
		},
		{
			`1nonunit`,
			&ast.NumberLit{
				Value: mustParseBigFloat("1"),
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 2, Byte: 1},
					},
				},
			},
			1, // extra characters after expression
		},

		{
			`("hello")`,
			&ast.ParenExpr{
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 10, Byte: 9},
					},
				},
				Content: &ast.StringLit{
					Value: "hello",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 2, Byte: 1},
							End:   source.Pos{Line: 1, Column: 9, Byte: 8},
						},
					},
				},
			},
			0,
		},

		{
			`-1`,
			&ast.ArithmeticUnary{
				Op: ast.Negate,
				Operand: &ast.NumberLit{
					Value: mustParseBigFloat("1"),
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 2, Byte: 1},
							End:   source.Pos{Line: 1, Column: 3, Byte: 2},
						},
					},
				},

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 3, Byte: 2},
					},
				},
			},
			0,
		},
		{
			`-1 + 2`,
			&ast.ArithmeticBinary{
				Op: ast.Add,
				LHS: &ast.ArithmeticUnary{
					Op: ast.Negate,
					Operand: &ast.NumberLit{
						Value: mustParseBigFloat("1"),
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 2, Byte: 1},
								End:   source.Pos{Line: 1, Column: 3, Byte: 2},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 3, Byte: 2},
						},
					},
				},
				RHS: &ast.NumberLit{
					Value: mustParseBigFloat("2"),
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 6, Byte: 5},
							End:   source.Pos{Line: 1, Column: 7, Byte: 6},
						},
					},
				},
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 7, Byte: 6},
					},
				},
			},
			0,
		},

		{
			`!true`,
			&ast.ArithmeticUnary{
				Op: ast.Not,
				Operand: &ast.BooleanLit{
					Value: true,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 2, Byte: 1},
							End:   source.Pos{Line: 1, Column: 6, Byte: 5},
						},
					},
				},

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 6, Byte: 5},
					},
				},
			},
			0,
		},

		{
			`"hello " .. "world"`,
			&ast.ArithmeticBinary{
				Op: ast.Concat,
				LHS: &ast.StringLit{
					Value: "hello ",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 9, Byte: 8},
						},
					},
				},
				RHS: &ast.StringLit{
					Value: "world",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 13, Byte: 12},
							End:   source.Pos{Line: 1, Column: 20, Byte: 19},
						},
					},
				},
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 20, Byte: 19},
					},
				},
			},
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, diags := ParseExpr([]byte(test.Input))

			if len(diags) != test.DiagCount {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.DiagCount)
				for _, diag := range diags {
					t.Logf("- %s", diag.String())
				}
			}

			prettyConfig := &pretty.Config{
				Diffable:          true,
				IncludeUnexported: true,
				PrintStringers:    false,
			}

			if !reflect.DeepEqual(got, test.Want) {
				diff := prettyConfig.Compare(test.Want, got)
				t.Errorf("wrong result\ninput:\n%s\n\ndiff: %s", test.Input, diff)
			}
		})
	}
}
