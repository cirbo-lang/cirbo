package parser

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/source"
	"github.com/kylelemons/godebug/pretty"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func TestParsePackage(t *testing.T) {
	tests := map[string]struct {
		Files     map[string]string
		PPath     string
		From      string
		Want      *ast.Package
		DiagCount int
	}{
		"empty package": {
			map[string]string{
				"cirbo-pkg/test/placeholder/.exists-but-empty": "",
			},
			"test/placeholder",
			"",
			&ast.Package{
				DefaultName: "placeholder",
			},
			0,
		},
		"missing package": {
			map[string]string{
				"cirbo-pkg/test/placeholder/.exists-but-empty": "",
			},
			"test/nonexist",
			"",
			&ast.Package{
				DefaultName: "nonexist",
			},
			1, // could not be resolved as a package path
		},
		"single-file package": {
			map[string]string{
				"cirbo-pkg/test/pkg/test.cbm": `import "test/other";`,
			},
			"test/pkg",
			"",
			&ast.Package{
				DefaultName: "pkg",

				Files: []*ast.File{
					{
						Source: []byte(`import "test/other";`),

						TopLevel: []ast.Node{
							&ast.Import{
								Package: "test/other",
								Name:    "",

								PackageRange: source.Range{
									Filename: "/cirbo-pkg/test/pkg/test.cbm",
									Start:    source.Pos{Line: 1, Column: 8, Byte: 7},
									End:      source.Pos{Line: 1, Column: 20, Byte: 19},
								},

								WithRange: ast.WithRange{
									Range: source.Range{
										Filename: "/cirbo-pkg/test/pkg/test.cbm",
										Start:    source.StartPos,
										End:      source.Pos{Line: 1, Column: 21, Byte: 20},
									},
								},
							},
						},

						WithRange: ast.WithRange{
							Range: source.Range{
								Filename: "/cirbo-pkg/test/pkg/test.cbm",
								Start:    source.StartPos,
								End:      source.Pos{Line: 1, Column: 21, Byte: 20},
							},
						},
					},
				},
			},
			0,
		},
		"package with ignored files": {
			map[string]string{
				"cirbo-pkg/test/pkg/test.cbm":   ``,
				"cirbo-pkg/test/pkg/.test.cbm":  ``,
				"cirbo-pkg/test/pkg/_test.cbm":  ``,
				"cirbo-pkg/test/pkg/project.cb": ``,
				"cirbo-pkg/test/pkg/README.md":  ``,
			},
			"test/pkg",
			"",
			&ast.Package{
				DefaultName: "pkg",

				Files: []*ast.File{
					{
						Source: []byte(``),

						WithRange: ast.WithRange{
							Range: source.Range{
								Filename: "/cirbo-pkg/test/pkg/test.cbm",
								Start:    source.StartPos,
								End:      source.StartPos,
							},
						},
					},
				},
			},
			0,
		},
		"multi-file package": {
			map[string]string{
				"cirbo-pkg/test/pkg/aaaa.cbm":       `import "test/other";`,
				"cirbo-pkg/test/pkg/zzzz.cbm":       `import "test/foooo";`,
				"cirbo-pkg/test/irrelevant/foo.cbm": ``,
			},
			"test/pkg",
			"",
			&ast.Package{
				DefaultName: "pkg",

				Files: []*ast.File{
					{
						Source: []byte(`import "test/other";`),

						TopLevel: []ast.Node{
							&ast.Import{
								Package: "test/other",
								Name:    "",

								PackageRange: source.Range{
									Filename: "/cirbo-pkg/test/pkg/aaaa.cbm",
									Start:    source.Pos{Line: 1, Column: 8, Byte: 7},
									End:      source.Pos{Line: 1, Column: 20, Byte: 19},
								},

								WithRange: ast.WithRange{
									Range: source.Range{
										Filename: "/cirbo-pkg/test/pkg/aaaa.cbm",
										Start:    source.StartPos,
										End:      source.Pos{Line: 1, Column: 21, Byte: 20},
									},
								},
							},
						},

						WithRange: ast.WithRange{
							Range: source.Range{
								Filename: "/cirbo-pkg/test/pkg/aaaa.cbm",
								Start:    source.StartPos,
								End:      source.Pos{Line: 1, Column: 21, Byte: 20},
							},
						},
					},
					{
						Source: []byte(`import "test/foooo";`),

						TopLevel: []ast.Node{
							&ast.Import{
								Package: "test/foooo",
								Name:    "",

								PackageRange: source.Range{
									Filename: "/cirbo-pkg/test/pkg/zzzz.cbm",
									Start:    source.Pos{Line: 1, Column: 8, Byte: 7},
									End:      source.Pos{Line: 1, Column: 20, Byte: 19},
								},

								WithRange: ast.WithRange{
									Range: source.Range{
										Filename: "/cirbo-pkg/test/pkg/zzzz.cbm",
										Start:    source.StartPos,
										End:      source.Pos{Line: 1, Column: 21, Byte: 20},
									},
								},
							},
						},

						WithRange: ast.WithRange{
							Range: source.Range{
								Filename: "/cirbo-pkg/test/pkg/zzzz.cbm",
								Start:    source.StartPos,
								End:      source.Pos{Line: 1, Column: 21, Byte: 20},
							},
						},
					},
				},
			},
			0,
		},
		"dot relative path": {
			map[string]string{
				"project.cb":     ``,
				"local/test.cbm": ``,
			},
			"./local",
			"/project.cb",
			&ast.Package{
				DefaultName: "local",

				Files: []*ast.File{
					{
						Source: []byte(``),

						WithRange: ast.WithRange{
							Range: source.Range{
								Filename: "/local/test.cbm",
								Start:    source.StartPos,
								End:      source.StartPos,
							},
						},
					},
				},
			},
			0,
		},
		"dotdot relative path": {
			map[string]string{
				"a/test.cbm": ``,
				"b/test.cbm": ``,
			},
			"../b",
			"/a/test.cbm",
			&ast.Package{
				DefaultName: "b",

				Files: []*ast.File{
					{
						Source: []byte(``),

						WithRange: ast.WithRange{
							Range: source.Range{
								Filename: "/b/test.cbm",
								Start:    source.StartPos,
								End:      source.StartPos,
							},
						},
					},
				},
			},
			0,
		},
		"dotdot through root": {
			map[string]string{
				"test.cb": ``,
			},
			"../b",
			"test.cb",
			&ast.Package{
				DefaultName: "b",
			},
			1, // could not be resolved as a package path
		},
		"slash-prefixed module path": {
			map[string]string{
				// This shouldn't match because it would be "./invalid"
				"invalid/test.cb": ``,

				// This also shouldn't match because it would be just "invalid"
				"cirbo-pkg/invalid/test/cb": ``,
			},
			"/invalid",
			"",
			&ast.Package{
				DefaultName: "invalid",
			},
			1, // could not be resolved as a package path
		},
		"file with errors": {
			map[string]string{
				"cirbo-pkg/test-pkg/foo.cbm": `import ?`,
			},
			"test-pkg",
			"/project.cb",
			&ast.Package{
				DefaultName: "test-pkg",

				Files: []*ast.File{
					{
						Source: []byte(`import ?`),

						WithRange: ast.WithRange{
							Range: source.Range{
								Filename: "/cirbo-pkg/test-pkg/foo.cbm",
								Start:    source.StartPos,
								End:      source.Pos{Line: 1, Column: 7, Byte: 6},
							},
						},
					},
				},
			},
			1, // import path must be a quoted string
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := inMemParser(test.Files)
			got, diags := p.ParsePackage(test.PPath, test.From)

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
				t.Errorf("wrong result\n%s", diff)
			}
		})
	}
}

func TestParseAllProjectFiles(t *testing.T) {
	tests := map[string]struct {
		Files     map[string]string
		Want      []*ast.File
		DiagCount int
	}{
		"single file": {
			map[string]string{
				"main.cb":   ``,
				"README.md": `# Hello World`,

				"cirbo-pkg/test-pkg/public.cbm": ``,
			},
			[]*ast.File{
				{
					Source: []byte{},

					WithRange: ast.WithRange{
						Range: source.Range{
							Filename: "/main.cb",
							Start:    source.StartPos,
							End:      source.StartPos,
						},
					},
				},
			},
			0,
		},
		"multiple files": {
			map[string]string{
				"a.cb":  ``,
				"b.cb":  ``,
				"c.cb":  ``,
				"_c.cb": ``,
				".c.cb": ``,
			},
			[]*ast.File{
				{
					Source: []byte{},

					WithRange: ast.WithRange{
						Range: source.Range{
							Filename: "/a.cb",
							Start:    source.StartPos,
							End:      source.StartPos,
						},
					},
				},
				{
					Source: []byte{},

					WithRange: ast.WithRange{
						Range: source.Range{
							Filename: "/b.cb",
							Start:    source.StartPos,
							End:      source.StartPos,
						},
					},
				},
				{
					Source: []byte{},

					WithRange: ast.WithRange{
						Range: source.Range{
							Filename: "/c.cb",
							Start:    source.StartPos,
							End:      source.StartPos,
						},
					},
				},
			},
			0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := inMemParser(test.Files)
			got, diags := p.ParseAllProjectFiles()

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
				t.Errorf("wrong result\n%s", diff)
			}
		})
	}

}

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

func inMemParser(files map[string]string) *Parser {
	fs := mapfs.New(files)

	return &Parser{
		fs:       fs,
		files:    map[string]*ast.File{},
		packages: map[string]*ast.Package{},
	}
}