package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cirbo-lang/cirbo/source"
	"github.com/kylelemons/godebug/pretty"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		input string
		want  []Token
	}{
		// Empty input
		{
			``,
			[]Token{
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 0, Line: 1, Column: 1},
					},
				},
			},
		},

		{
			`foo`,
			[]Token{
				{
					Type:  TokenIdent,
					Bytes: []byte(`foo`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 3, Line: 1, Column: 4},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 3, Line: 1, Column: 4},
						End:   source.Pos{Byte: 3, Line: 1, Column: 4},
					},
				},
			},
		},
		{
			`~foo`,
			[]Token{
				{
					Type:  TokenIdent,
					Bytes: []byte(`~foo`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},
		{
			`foo~baz`,
			[]Token{
				{
					Type:  TokenIdent,
					Bytes: []byte(`foo~baz`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 7, Line: 1, Column: 8},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 7, Line: 1, Column: 8},
						End:   source.Pos{Byte: 7, Line: 1, Column: 8},
					},
				},
			},
		},
		{
			"`foo`",
			[]Token{
				{
					Type:  TokenIdent,
					Bytes: []byte("`foo`"),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 5, Line: 1, Column: 6},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 5, Line: 1, Column: 6},
						End:   source.Pos{Byte: 5, Line: 1, Column: 6},
					},
				},
			},
		},

		{
			`1234`,
			[]Token{
				{
					Type:  TokenNumberLit,
					Bytes: []byte(`1234`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},
		{
			`12.4`,
			[]Token{
				{
					Type:  TokenNumberLit,
					Bytes: []byte(`12.4`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},
		{
			`12e1`,
			[]Token{
				{
					Type:  TokenNumberLit,
					Bytes: []byte(`12e1`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},
		{
			`1e-1`,
			[]Token{
				{
					Type:  TokenNumberLit,
					Bytes: []byte(`1e-1`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},
		{
			`1e+1`,
			[]Token{
				{
					Type:  TokenNumberLit,
					Bytes: []byte(`1e+1`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},
		{
			`100%`,
			[]Token{
				{
					Type:  TokenNumberLit,
					Bytes: []byte(`100`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 3, Line: 1, Column: 4},
					},
				},
				{
					Type:  TokenPercent,
					Bytes: []byte(`%`),
					Range: source.Range{
						Start: source.Pos{Byte: 3, Line: 1, Column: 4},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},
		{
			`10mm`,
			[]Token{
				{
					Type:  TokenNumberLit,
					Bytes: []byte(`10`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 2, Line: 1, Column: 3},
					},
				},
				{
					Type:  TokenIdent,
					Bytes: []byte(`mm`),
					Range: source.Range{
						Start: source.Pos{Byte: 2, Line: 1, Column: 3},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
			},
		},

		{
			`"hello"`,
			[]Token{
				{
					Type:  TokenStringLit,
					Bytes: []byte(`"hello"`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 7, Line: 1, Column: 8},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 7, Line: 1, Column: 8},
						End:   source.Pos{Byte: 7, Line: 1, Column: 8},
					},
				},
			},
		},

		{
			`foo -- bar`,
			[]Token{
				{
					Type:  TokenIdent,
					Bytes: []byte(`foo`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 3, Line: 1, Column: 4},
					},
				},
				{
					Type:  TokenWhitespace,
					Bytes: []byte(` `),
					Range: source.Range{
						Start: source.Pos{Byte: 3, Line: 1, Column: 4},
						End:   source.Pos{Byte: 4, Line: 1, Column: 5},
					},
				},
				{
					Type:  TokenDashDash,
					Bytes: []byte(`--`),
					Range: source.Range{
						Start: source.Pos{Byte: 4, Line: 1, Column: 5},
						End:   source.Pos{Byte: 6, Line: 1, Column: 7},
					},
				},
				{
					Type:  TokenWhitespace,
					Bytes: []byte(` `),
					Range: source.Range{
						Start: source.Pos{Byte: 6, Line: 1, Column: 7},
						End:   source.Pos{Byte: 7, Line: 1, Column: 8},
					},
				},
				{
					Type:  TokenIdent,
					Bytes: []byte(`bar`),
					Range: source.Range{
						Start: source.Pos{Byte: 7, Line: 1, Column: 8},
						End:   source.Pos{Byte: 10, Line: 1, Column: 11},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 10, Line: 1, Column: 11},
						End:   source.Pos{Byte: 10, Line: 1, Column: 11},
					},
				},
			},
		},

		{
			"// bigly\nhello",
			[]Token{
				{
					Type:  TokenComment,
					Bytes: []byte("// bigly\n"),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 9, Line: 2, Column: 1},
					},
				},
				{
					Type:  TokenIdent,
					Bytes: []byte(`hello`),
					Range: source.Range{
						Start: source.Pos{Byte: 9, Line: 2, Column: 1},
						End:   source.Pos{Byte: 14, Line: 2, Column: 6},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 14, Line: 2, Column: 6},
						End:   source.Pos{Byte: 14, Line: 2, Column: 6},
					},
				},
			},
		},
		{
			"/* hi */\nhello",
			[]Token{
				{
					Type:  TokenComment,
					Bytes: []byte(`/* hi */`),
					Range: source.Range{
						Start: source.Pos{Byte: 0, Line: 1, Column: 1},
						End:   source.Pos{Byte: 8, Line: 1, Column: 9},
					},
				},
				{
					Type:  TokenWhitespace,
					Bytes: []byte("\n"),
					Range: source.Range{
						Start: source.Pos{Byte: 8, Line: 1, Column: 9},
						End:   source.Pos{Byte: 9, Line: 2, Column: 1},
					},
				},
				{
					Type:  TokenIdent,
					Bytes: []byte(`hello`),
					Range: source.Range{
						Start: source.Pos{Byte: 9, Line: 2, Column: 1},
						End:   source.Pos{Byte: 14, Line: 2, Column: 6},
					},
				},
				{
					Type:  TokenEOF,
					Bytes: []byte{},
					Range: source.Range{
						Start: source.Pos{Byte: 14, Line: 2, Column: 6},
						End:   source.Pos{Byte: 14, Line: 2, Column: 6},
					},
				},
			},
		},
	}

	prettyConfig := &pretty.Config{
		Diffable:          true,
		IncludeUnexported: true,
		PrintStringers:    true,
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := scanTokens([]byte(test.input), "", source.Pos{Byte: 0, Line: 1, Column: 1}, scanNormal)

			if !reflect.DeepEqual(got, test.want) {
				diff := prettyConfig.Compare(test.want, got)
				if strings.TrimSpace(diff) == "" {
					// Sometimes stringers obscure differences
					prettyConfig.PrintStringers = false
					diff = prettyConfig.Compare(test.want, got)
					prettyConfig.PrintStringers = true
				}
				t.Errorf(
					"wrong result\ninput: %s\ndiff:  %s",
					test.input, diff,
				)
			}
		})
	}
}
