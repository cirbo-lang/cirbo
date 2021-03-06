package parser

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/cbo"
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

		{
			`export "baz";`,
			[]ast.Node{
				&ast.Export{
					Value: &ast.StringLit{
						Value: "baz",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},
					Name: "",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			0,
		},
		{
			`export false;`,
			[]ast.Node{
				&ast.Export{
					Value: &ast.BooleanLit{
						Value: false,

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},
					Name: "",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			0,
		},
		{
			`export "baz" as thing;`,
			[]ast.Node{
				&ast.Export{
					Value: &ast.StringLit{
						Value: "baz",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},
					Name: "thing",
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 23, Byte: 22},
						},
					},
				},
			},
			0,
		},

		{
			`designator "R";`,
			[]ast.Node{
				&ast.Designator{
					Value: &ast.StringLit{
						Value: "R",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 15, Byte: 14},
							},
						},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 16, Byte: 15},
						},
					},
				},
			},
			0,
		},

		{
			`attr value resistance;`,
			[]ast.Node{
				&ast.Attr{
					Name: "value",
					Type: &ast.Variable{
						Name: "resistance",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 22, Byte: 21},
							},
						},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 23, Byte: 22},
						},
					},
				},
			},
			0,
		},
		{
			`attr value = 1ohm;`,
			[]ast.Node{
				&ast.Attr{
					Name: "value",
					Value: &ast.NumberLit{
						Value: mustParseBigFloat("1"),
						Unit:  "ohm",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 14, Byte: 13},
								End:   source.Pos{Line: 1, Column: 18, Byte: 17},
							},
						},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 19, Byte: 18},
						},
					},
				},
			},
			0,
		},
		{
			`attr value;`,
			[]ast.Node{
				&ast.Attr{
					Name: "value",
					Type: &ast.Invalid{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 12, Byte: 11},
							},
						},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 12, Byte: 11},
						},
					},
				},
			},
			1, // invalid attribute definition (missing type or value)
		},

		{
			`a = true;`,
			[]ast.Node{
				&ast.Assign{
					Name: "a",
					Value: &ast.BooleanLit{
						Value: true,

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 5, Byte: 4},
								End:   source.Pos{Line: 1, Column: 9, Byte: 8},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 10, Byte: 9},
						},
					},
				},
			},
			0,
		},
		{
			`a = true indeed;`,
			[]ast.Node{
				&ast.Assign{
					Name: "a",
					Value: &ast.BooleanLit{
						Value: true,

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 5, Byte: 4},
								End:   source.Pos{Line: 1, Column: 9, Byte: 8},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 16, Byte: 15},
						},
					},
				},
			},
			1, // unterminated statement
		},
		{
			`a = true indeed; b = true;`,
			[]ast.Node{
				&ast.Assign{
					Name: "a",
					Value: &ast.BooleanLit{
						Value: true,

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 5, Byte: 4},
								End:   source.Pos{Line: 1, Column: 9, Byte: 8},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 16, Byte: 15},
						},
					},
				},
				// We should recover from the error in the first statement
				// and then successfully parse the second, below.
				&ast.Assign{
					Name: "b",
					Value: &ast.BooleanLit{
						Value: true,

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 22, Byte: 21},
								End:   source.Pos{Line: 1, Column: 26, Byte: 25},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 18, Byte: 17},
							End:   source.Pos{Line: 1, Column: 27, Byte: 26},
						},
					},
				},
			},
			1, // unterminated statement
		},
		{
			`false = true;`,
			[]ast.Node{
				&ast.Assign{
					Name: "",
					Value: &ast.BooleanLit{
						Value: true,

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 9, Byte: 8},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			1, // invalid assignment expression (can't assign to boolean literal)
		},
		{
			`true;`,
			[]ast.Node{
				&ast.BooleanLit{
					Value: true,

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 5, Byte: 4},
						},
					},
				},
			},
			1, // Useless naked expression
		},

		{
			`|-- GND;`,
			[]ast.Node{
				&ast.NoConnection{
					Terminal: &ast.Variable{
						Name: "GND",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 5, Byte: 4},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 9, Byte: 8},
						},
					},
				},
			},
			0,
		},
		{
			`GND --|;`,
			[]ast.Node{
				&ast.NoConnection{
					Terminal: &ast.Variable{
						Name: "GND",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 1, Byte: 0},
								End:   source.Pos{Line: 1, Column: 4, Byte: 3},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 9, Byte: 8},
						},
					},
				},
			},
			0,
		},
		{
			`GND -- PGND;`,
			[]ast.Node{
				&ast.Connection{
					Seq: []ast.Node{
						&ast.Variable{
							Name: "GND",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 1, Byte: 0},
									End:   source.Pos{Line: 1, Column: 4, Byte: 3},
								},
							},
						},
						&ast.Variable{
							Name: "PGND",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 8, Byte: 7},
									End:   source.Pos{Line: 1, Column: 12, Byte: 11},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 13, Byte: 12},
						},
					},
				},
			},
			0,
		},
		{
			`GND -- PGND -- AGND;`,
			[]ast.Node{
				&ast.Connection{
					Seq: []ast.Node{
						&ast.Variable{
							Name: "GND",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 1, Byte: 0},
									End:   source.Pos{Line: 1, Column: 4, Byte: 3},
								},
							},
						},
						&ast.Variable{
							Name: "PGND",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 8, Byte: 7},
									End:   source.Pos{Line: 1, Column: 12, Byte: 11},
								},
							},
						},
						&ast.Variable{
							Name: "AGND",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 16, Byte: 15},
									End:   source.Pos{Line: 1, Column: 20, Byte: 19},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 21, Byte: 20},
						},
					},
				},
			},
			0,
		},
		{
			`MODE -- R1(1kohm) -- VCC;`,
			[]ast.Node{
				&ast.Connection{
					Seq: []ast.Node{
						&ast.Variable{
							Name: "MODE",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 1, Byte: 0},
									End:   source.Pos{Line: 1, Column: 5, Byte: 4},
								},
							},
						},
						&ast.Call{
							Callee: &ast.Variable{
								Name: "R1",

								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 9, Byte: 8},
										End:   source.Pos{Line: 1, Column: 11, Byte: 10},
									},
								},
							},
							Args: &ast.Arguments{
								Positional: []ast.Node{
									&ast.NumberLit{
										Value: mustParseBigFloat("1"),
										Unit:  "kohm",

										WithRange: ast.WithRange{
											Range: source.Range{
												Start: source.Pos{Line: 1, Column: 12, Byte: 11},
												End:   source.Pos{Line: 1, Column: 17, Byte: 16},
											},
										},
									},
								},

								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 11, Byte: 10},
										End:   source.Pos{Line: 1, Column: 18, Byte: 17},
									},
								},
							},

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 9, Byte: 8},
									End:   source.Pos{Line: 1, Column: 18, Byte: 17},
								},
							},
						},
						&ast.Variable{
							Name: "VCC",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 22, Byte: 21},
									End:   source.Pos{Line: 1, Column: 25, Byte: 24},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 26, Byte: 25},
						},
					},
				},
			},
			0,
		},
		{
			`GND --;`,
			[]ast.Node{
				&ast.Connection{
					Seq: []ast.Node{
						&ast.Variable{
							Name: "GND",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 1, Byte: 0},
									End:   source.Pos{Line: 1, Column: 4, Byte: 3},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 8, Byte: 7},
						},
					},
				},
			},
			1, // missing terminal expression
		},

		{
			`circuit foo {}`,
			[]ast.Node{
				&ast.Circuit{
					Name: "foo",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 13, Byte: 12},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 13, Byte: 12},
								End:   source.Pos{Line: 1, Column: 15, Byte: 14},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 12, Byte: 11},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 15, Byte: 14},
						},
					},
				},
			},
			0,
		},
		{
			`circuit foo(bar, baz) {}`,
			[]ast.Node{
				&ast.Circuit{
					Name: "foo",
					Params: &ast.Arguments{
						Positional: []ast.Node{
							&ast.Variable{
								Name: "bar",
								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 13, Byte: 12},
										End:   source.Pos{Line: 1, Column: 16, Byte: 15},
									},
								},
							},
							&ast.Variable{
								Name: "baz",
								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 18, Byte: 17},
										End:   source.Pos{Line: 1, Column: 21, Byte: 20},
									},
								},
							},
						},

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 22, Byte: 21},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 23, Byte: 22},
								End:   source.Pos{Line: 1, Column: 25, Byte: 24},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 22, Byte: 21},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 25, Byte: 24},
						},
					},
				},
			},
			0,
		},
		{
			`circuit foo { import "baz"; }`, // import not semantically valid here, but okay syntax-wise
			[]ast.Node{
				&ast.Circuit{
					Name: "foo",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 13, Byte: 12},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},
					Body: &ast.StatementBlock{
						Statements: []ast.Node{
							&ast.Import{
								Package: "baz",
								PackageRange: source.Range{
									Start: source.Pos{Line: 1, Column: 22, Byte: 21},
									End:   source.Pos{Line: 1, Column: 27, Byte: 26},
								},
								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 15, Byte: 14},
										End:   source.Pos{Line: 1, Column: 28, Byte: 27},
									},
								},
							},
						},

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 13, Byte: 12},
								End:   source.Pos{Line: 1, Column: 30, Byte: 29},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 12, Byte: 11},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 30, Byte: 29},
						},
					},
				},
			},
			0,
		},
		{
			"circuit `foo` {}",
			[]ast.Node{
				&ast.Circuit{
					Name: "foo",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 15, Byte: 14},
								End:   source.Pos{Line: 1, Column: 15, Byte: 14},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 15, Byte: 14},
								End:   source.Pos{Line: 1, Column: 17, Byte: 16},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 14, Byte: 13},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 17, Byte: 16},
						},
					},
				},
			},
			0,
		},
		{
			`circuit "foo" {}`,
			[]ast.Node{
				&ast.Circuit{
					Name: "",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 8, Byte: 7},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			1, // circuit name must be an identifier
		},
		{
			`circuit "foo" {} circuit bar {}`,
			[]ast.Node{
				&ast.Circuit{
					Name: "",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 8, Byte: 7},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},

				// should recover from error in first circuit and then parse the second
				&ast.Circuit{
					Name: "bar",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 30, Byte: 29},
								End:   source.Pos{Line: 1, Column: 30, Byte: 29},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 30, Byte: 29},
								End:   source.Pos{Line: 1, Column: 32, Byte: 31},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 18, Byte: 17},
						End:   source.Pos{Line: 1, Column: 29, Byte: 28},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 18, Byte: 17},
							End:   source.Pos{Line: 1, Column: 32, Byte: 31},
						},
					},
				},
			},
			1, // circuit name must be an identifier
		},
		{
			`circuit foo {`,
			[]ast.Node{
				&ast.Circuit{
					Name: "foo",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 13, Byte: 12},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 13, Byte: 12},
								End:   source.Pos{Line: 1, Column: 14, Byte: 13},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 12, Byte: 11},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			1, // unclosed statement block
		},
		{
			`circuit {}`,
			[]ast.Node{
				&ast.Circuit{
					Name: "",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 8, Byte: 7},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 8, Byte: 7},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 10, Byte: 9},
						},
					},
				},
			},
			1, // missing circuit name
		},
		{
			`circuit foo(bar, "a") {}`,
			[]ast.Node{
				&ast.Circuit{
					Name: "foo",
					Params: &ast.Arguments{
						Positional: []ast.Node{
							&ast.Variable{
								Name: "bar",
								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 13, Byte: 12},
										End:   source.Pos{Line: 1, Column: 16, Byte: 15},
									},
								},
							},
							&ast.StringLit{
								Value: "a",
								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 18, Byte: 17},
										End:   source.Pos{Line: 1, Column: 21, Byte: 20},
									},
								},
							},
						},

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 22, Byte: 21},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 23, Byte: 22},
								End:   source.Pos{Line: 1, Column: 25, Byte: 24},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 22, Byte: 21},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 25, Byte: 24},
						},
					},
				},
			},
			1, // invalid parameter declaration (can't use string literal)
		},
		{
			`circuit foo(bar, a=1) {}`,
			[]ast.Node{
				&ast.Circuit{
					Name: "foo",
					Params: &ast.Arguments{
						Positional: []ast.Node{
							&ast.Variable{
								Name: "bar",
								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 13, Byte: 12},
										End:   source.Pos{Line: 1, Column: 16, Byte: 15},
									},
								},
							},
						},
						Named: []*ast.NamedArgument{
							{
								Name: "a",
								Value: &ast.NumberLit{
									Value: mustParseBigFloat("1"),

									WithRange: ast.WithRange{
										Range: source.Range{
											Start: source.Pos{Line: 1, Column: 20, Byte: 19},
											End:   source.Pos{Line: 1, Column: 21, Byte: 20},
										},
									},
								},

								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 18, Byte: 17},
										End:   source.Pos{Line: 1, Column: 21, Byte: 20},
									},
								},
							},
						},

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 22, Byte: 21},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 23, Byte: 22},
								End:   source.Pos{Line: 1, Column: 25, Byte: 24},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 22, Byte: 21},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 25, Byte: 24},
						},
					},
				},
			},
			1, // invalid parameter declaration (can't use named argument)
		},

		{
			`board foo {}`,
			[]ast.Node{
				&ast.Board{
					Name: "foo",
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 11, Byte: 10},
								End:   source.Pos{Line: 1, Column: 13, Byte: 12},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 10, Byte: 9},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 13, Byte: 12},
						},
					},
				},
			},
			0,
		},
		{
			`board foo() {}`,
			[]ast.Node{
				&ast.Board{
					Name: "foo",
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 13, Byte: 12},
								End:   source.Pos{Line: 1, Column: 15, Byte: 14},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 10, Byte: 9},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 15, Byte: 14},
						},
					},
				},
			},
			0,
		},
		{
			`board foo(bar) {}`,
			[]ast.Node{
				&ast.Board{
					Name: "foo",
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 16, Byte: 15},
								End:   source.Pos{Line: 1, Column: 18, Byte: 17},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 15, Byte: 14},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 18, Byte: 17},
						},
					},
				},
			},
			1, // no parameter list is allowed
		},

		{
			`device foo {}`,
			[]ast.Node{
				&ast.Device{
					Name: "foo",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 12, Byte: 11},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 14, Byte: 13},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 11, Byte: 10},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			0,
		},

		{
			`land foo {}`,
			[]ast.Node{
				&ast.Land{
					Name: "foo",
					Params: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 10, Byte: 9},
								End:   source.Pos{Line: 1, Column: 10, Byte: 9},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 10, Byte: 9},
								End:   source.Pos{Line: 1, Column: 12, Byte: 11},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 9, Byte: 8},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 12, Byte: 11},
						},
					},
				},
			},
			0,
		},

		{
			`pinout foo to bar {}`,
			[]ast.Node{
				&ast.Pinout{
					Name: "foo",
					Land: &ast.Variable{
						Name: "bar",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 15, Byte: 14},
								End:   source.Pos{Line: 1, Column: 18, Byte: 17},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 19, Byte: 18},
								End:   source.Pos{Line: 1, Column: 21, Byte: 20},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 18, Byte: 17},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 21, Byte: 20},
						},
					},
				},
			},
			0,
		},
		{
			`pinout foo from baz to bar {}`,
			[]ast.Node{
				&ast.Pinout{
					Name: "foo",
					Device: &ast.Variable{
						Name: "baz",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 17, Byte: 16},
								End:   source.Pos{Line: 1, Column: 20, Byte: 19},
							},
						},
					},
					Land: &ast.Variable{
						Name: "bar",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 24, Byte: 23},
								End:   source.Pos{Line: 1, Column: 27, Byte: 26},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 28, Byte: 27},
								End:   source.Pos{Line: 1, Column: 30, Byte: 29},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 27, Byte: 26},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 30, Byte: 29},
						},
					},
				},
			},
			0,
		},
		{
			`pinout foo {}`,
			[]ast.Node{
				&ast.Pinout{
					Name: "foo",
					Land: &ast.Invalid{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 7, Byte: 6},
								End:   source.Pos{Line: 1, Column: 7, Byte: 6},
							},
						},
					},
					Body: &ast.StatementBlock{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 12, Byte: 11},
								End:   source.Pos{Line: 1, Column: 14, Byte: 13},
							},
						},
					},

					HeaderRange: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 11, Byte: 10},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			1, // missing "to" clause
		},

		{
			`terminal foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Passive,
					Dir:  cbo.Undirected,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
			},
			0,
		},
		{
			`power foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Power,
					Dir:  cbo.Undirected,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 11, Byte: 10},
						},
					},
				},
			},
			0,
		},
		{
			`input foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Signal,
					Dir:  cbo.Input,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 11, Byte: 10},
						},
					},
				},
			},
			0,
		},
		{
			`output foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name:       "foo",
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 12, Byte: 11},
						},
					},
				},
			},
			0,
		},
		{
			`output tristate foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name:       "foo",
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.Tristate,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 21, Byte: 20},
						},
					},
				},
			},
			0,
		},
		{
			`output emitter foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name:       "foo",
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.OpenEmitter,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 20, Byte: 19},
						},
					},
				},
			},
			0,
		},
		{
			`output collector foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name:       "foo",
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.OpenCollector,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 22, Byte: 21},
						},
					},
				},
			},
			0,
		},
		{
			`bidi leader foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Signal,
					Dir:  cbo.Bidirectional,
					Role: cbo.Leader,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 17, Byte: 16},
						},
					},
				},
			},
			0,
		},
		{
			`bidi follower foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Signal,
					Dir:  cbo.Bidirectional,
					Role: cbo.Follower,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 19, Byte: 18},
						},
					},
				},
			},
			0,
		},
		{
			`input follower foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Signal,
					Dir:  cbo.Input,
					Role: cbo.Follower,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 20, Byte: 19},
						},
					},
				},
			},
			0,
		},
		{
			`output leader foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Signal,
					Dir:  cbo.Output,
					Role: cbo.Leader,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 19, Byte: 18},
						},
					},
				},
			},
			0,
		},
		{
			`output tristate leader foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name:       "foo",
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					Role:       cbo.Leader,
					OutputType: cbo.Tristate,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 28, Byte: 27},
						},
					},
				},
			},
			0,
		},
		{
			`bidi foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Signal,
					Dir:  cbo.Undirected,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 10, Byte: 9},
						},
					},
				},
			},
			1, // missing bidirectional terminal role
		},
		{
			`power input foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Power,
					Dir:  cbo.Input,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 17, Byte: 16},
						},
					},
				},
			},
			0,
		},
		{
			`power output foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name:       "foo",
					Type:       cbo.Power,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 18, Byte: 17},
						},
					},
				},
			},
			0,
		},
		{
			`power foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "foo",
					Type: cbo.Power,
					Dir:  cbo.Undirected,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 11, Byte: 10},
						},
					},
				},
			},
			0,
		},
		{
			`power power foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "power",
					Type: cbo.Power,
					Dir:  cbo.Undirected,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 12, Byte: 11},
						},
					},
				},
			},
			1, // invalid terminal declaration
		},
		{
			`input power foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "power",
					Type: cbo.Signal,
					Dir:  cbo.Input,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 12, Byte: 11},
						},
					},
				},
			},
			1, // invalid terminal declaration
		},
		{
			`input bidi foo;`,
			[]ast.Node{
				&ast.Terminal{
					Name: "bidi",
					Type: cbo.Signal,
					Dir:  cbo.Input,
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 11, Byte: 10},
						},
					},
				},
			},
			1, // invalid terminal declaration
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
			&ast.NumberLit{
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
			&ast.NumberLit{
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
			&ast.NumberLit{
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
			`("hello"`,
			&ast.ParenExpr{
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 9, Byte: 8},
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
			1, // expected a closing parenthesis
		},
		{
			`("hello" world!`,
			&ast.ParenExpr{
				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 15, Byte: 14},
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
			1, // expected a closing parenthesis
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

		{
			`foo.bar`,
			&ast.GetAttr{
				Source: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Name: "bar",

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
			`foo.bar.baz`,
			&ast.GetAttr{
				Source: &ast.GetAttr{
					Source: &ast.Variable{
						Name: "foo",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 1, Byte: 0},
								End:   source.Pos{Line: 1, Column: 4, Byte: 3},
							},
						},
					},
					Name: "bar",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 8, Byte: 7},
						},
					},
				},
				Name: "baz",

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 12, Byte: 11},
					},
				},
			},
			0,
		},
		{
			`foo.bar + baz`,
			&ast.ArithmeticBinary{
				LHS: &ast.GetAttr{
					Source: &ast.Variable{
						Name: "foo",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 1, Byte: 0},
								End:   source.Pos{Line: 1, Column: 4, Byte: 3},
							},
						},
					},
					Name: "bar",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 8, Byte: 7},
						},
					},
				},
				RHS: &ast.Variable{
					Name: "baz",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 11, Byte: 10},
							End:   source.Pos{Line: 1, Column: 14, Byte: 13},
						},
					},
				},
				Op: ast.Add,

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 14, Byte: 13},
					},
				},
			},
			0,
		},
		{
			`foo. + bar`,
			&ast.GetAttr{
				Source: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Name: "", // empty to indicate that it was invalid

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 5, Byte: 4},
					},
				},
			},
			1, // Invalid attribute name
		},

		{
			`foo[bar]`,
			&ast.GetIndex{
				Source: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Index: &ast.Variable{
					Name: "bar",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 5, Byte: 4},
							End:   source.Pos{Line: 1, Column: 8, Byte: 7},
						},
					},
				},

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
			`foo[bar][baz]`,
			&ast.GetIndex{
				Source: &ast.GetIndex{
					Source: &ast.Variable{
						Name: "foo",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 1, Byte: 0},
								End:   source.Pos{Line: 1, Column: 4, Byte: 3},
							},
						},
					},
					Index: &ast.Variable{
						Name: "bar",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 5, Byte: 4},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 9, Byte: 8},
						},
					},
				},
				Index: &ast.Variable{
					Name: "baz",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 10, Byte: 9},
							End:   source.Pos{Line: 1, Column: 13, Byte: 12},
						},
					},
				},

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 14, Byte: 13},
					},
				},
			},
			0,
		},
		{
			`foo[bar[baz]]`,
			&ast.GetIndex{
				Source: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 14, Byte: 13},
					},
				},
				Index: &ast.GetIndex{
					Source: &ast.Variable{
						Name: "bar",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 5, Byte: 4},
								End:   source.Pos{Line: 1, Column: 8, Byte: 7},
							},
						},
					},
					Index: &ast.Variable{
						Name: "baz",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 9, Byte: 8},
								End:   source.Pos{Line: 1, Column: 12, Byte: 11},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 5, Byte: 4},
							End:   source.Pos{Line: 1, Column: 13, Byte: 12},
						},
					},
				},
			},
			0,
		},

		{
			`foo()`,
			&ast.Call{
				Callee: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Args: &ast.Arguments{
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 4, Byte: 3},
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
			`foo(true)`,
			&ast.Call{
				Callee: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Args: &ast.Arguments{
					Positional: []ast.Node{
						&ast.BooleanLit{
							Value: true,

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 5, Byte: 4},
									End:   source.Pos{Line: 1, Column: 9, Byte: 8},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 4, Byte: 3},
							End:   source.Pos{Line: 1, Column: 10, Byte: 9},
						},
					},
				},

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 10, Byte: 9},
					},
				},
			},
			0,
		},
		{
			`foo(true, "?")`,
			&ast.Call{
				Callee: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Args: &ast.Arguments{
					Positional: []ast.Node{
						&ast.BooleanLit{
							Value: true,

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 5, Byte: 4},
									End:   source.Pos{Line: 1, Column: 9, Byte: 8},
								},
							},
						},
						&ast.StringLit{
							Value: "?",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 11, Byte: 10},
									End:   source.Pos{Line: 1, Column: 14, Byte: 13},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 4, Byte: 3},
							End:   source.Pos{Line: 1, Column: 15, Byte: 14},
						},
					},
				},

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 15, Byte: 14},
					},
				},
			},
			0,
		},
		{
			`foo(good=true)`,
			&ast.Call{
				Callee: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Args: &ast.Arguments{
					Named: []*ast.NamedArgument{
						{
							Name: "good",
							Value: &ast.BooleanLit{
								Value: true,

								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 10, Byte: 9},
										End:   source.Pos{Line: 1, Column: 14, Byte: 13},
									},
								},
							},

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 5, Byte: 4},
									End:   source.Pos{Line: 1, Column: 14, Byte: 13},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 4, Byte: 3},
							End:   source.Pos{Line: 1, Column: 15, Byte: 14},
						},
					},
				},

				WithRange: ast.WithRange{
					Range: source.Range{
						Start: source.Pos{Line: 1, Column: 1, Byte: 0},
						End:   source.Pos{Line: 1, Column: 15, Byte: 14},
					},
				},
			},
			0,
		},
		{
			`foo(bar, good=true)`,
			&ast.Call{
				Callee: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Args: &ast.Arguments{
					Positional: []ast.Node{
						&ast.Variable{
							Name: "bar",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 5, Byte: 4},
									End:   source.Pos{Line: 1, Column: 8, Byte: 7},
								},
							},
						},
					},
					Named: []*ast.NamedArgument{
						{
							Name: "good",
							Value: &ast.BooleanLit{
								Value: true,

								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 15, Byte: 14},
										End:   source.Pos{Line: 1, Column: 19, Byte: 18},
									},
								},
							},

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 10, Byte: 9},
									End:   source.Pos{Line: 1, Column: 19, Byte: 18},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 4, Byte: 3},
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
		{
			`foo(good=true, bar)`,
			&ast.Call{
				Callee: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Args: &ast.Arguments{
					Positional: []ast.Node{
						&ast.Variable{
							Name: "bar",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 16, Byte: 15},
									End:   source.Pos{Line: 1, Column: 19, Byte: 18},
								},
							},
						},
					},
					Named: []*ast.NamedArgument{
						{
							Name: "good",
							Value: &ast.BooleanLit{
								Value: true,

								WithRange: ast.WithRange{
									Range: source.Range{
										Start: source.Pos{Line: 1, Column: 10, Byte: 9},
										End:   source.Pos{Line: 1, Column: 14, Byte: 13},
									},
								},
							},

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 5, Byte: 4},
									End:   source.Pos{Line: 1, Column: 14, Byte: 13},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 4, Byte: 3},
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
			1, // incorrect argument order
		},
		{
			`foo()(a)`,
			&ast.Call{
				Callee: &ast.Call{
					Callee: &ast.Variable{
						Name: "foo",

						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 1, Byte: 0},
								End:   source.Pos{Line: 1, Column: 4, Byte: 3},
							},
						},
					},
					Args: &ast.Arguments{
						WithRange: ast.WithRange{
							Range: source.Range{
								Start: source.Pos{Line: 1, Column: 4, Byte: 3},
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
				Args: &ast.Arguments{
					Positional: []ast.Node{
						&ast.Variable{
							Name: "a",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 7, Byte: 6},
									End:   source.Pos{Line: 1, Column: 8, Byte: 7},
								},
							},
						},
					},
					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 6, Byte: 5},
							End:   source.Pos{Line: 1, Column: 9, Byte: 8},
						},
					},
				},

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
			`foo(a`,
			&ast.Call{
				Callee: &ast.Variable{
					Name: "foo",

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 1, Byte: 0},
							End:   source.Pos{Line: 1, Column: 4, Byte: 3},
						},
					},
				},
				Args: &ast.Arguments{
					Positional: []ast.Node{
						&ast.Variable{
							Name: "a",

							WithRange: ast.WithRange{
								Range: source.Range{
									Start: source.Pos{Line: 1, Column: 5, Byte: 4},
									End:   source.Pos{Line: 1, Column: 6, Byte: 5},
								},
							},
						},
					},

					WithRange: ast.WithRange{
						Range: source.Range{
							Start: source.Pos{Line: 1, Column: 4, Byte: 3},
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
			1, // missing argument separator
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
