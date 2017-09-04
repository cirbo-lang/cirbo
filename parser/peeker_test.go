package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cirbo-lang/cirbo/source"
	"github.com/kylelemons/godebug/pretty"
)

func TestTokenPeeker(t *testing.T) {
	input := `
hello world // foo
bar /* baz */ quux
`
	tokens := scanTokens([]byte(input), "", source.StartPos, scanNormal)
	it := newTokenIterator(tokens)

	p := &tokenPeeker{
		Iter: it,
	}
	var got Tokens
	for i := 0; i < 6; i++ {
		got = append(got, p.Read())
	}
	want := Tokens{
		{
			Type:  TokenIdent,
			Bytes: []byte("hello"),
			Range: source.Range{
				Start: source.Pos{Line: 2, Column: 1, Byte: 1},
				End:   source.Pos{Line: 2, Column: 6, Byte: 6},
			},
		},
		{
			Type:  TokenIdent,
			Bytes: []byte("world"),
			Range: source.Range{
				Start: source.Pos{Line: 2, Column: 7, Byte: 7},
				End:   source.Pos{Line: 2, Column: 12, Byte: 12},
			},
		},
		{
			Type:  TokenIdent,
			Bytes: []byte("bar"),
			Range: source.Range{
				Start: source.Pos{Line: 3, Column: 1, Byte: 20},
				End:   source.Pos{Line: 3, Column: 4, Byte: 23},
			},
		},
		{
			Type:  TokenIdent,
			Bytes: []byte("quux"),
			Range: source.Range{
				Start: source.Pos{Line: 3, Column: 15, Byte: 34},
				End:   source.Pos{Line: 3, Column: 19, Byte: 38},
			},
		},
		{
			Type:  TokenEOF,
			Bytes: []byte{},
			Range: source.Range{
				Start: source.Pos{Line: 4, Column: 1, Byte: 39},
				End:   source.Pos{Line: 4, Column: 1, Byte: 39},
			},
		},
		// TokenEOF is repeated indefinitely at the end
		{
			Type:  TokenEOF,
			Bytes: []byte{},
			Range: source.Range{
				Start: source.Pos{Line: 4, Column: 1, Byte: 39},
				End:   source.Pos{Line: 4, Column: 1, Byte: 39},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		prettyConfig := &pretty.Config{
			Diffable:          true,
			IncludeUnexported: true,
			PrintStringers:    true,
		}
		diff := prettyConfig.Compare(want, got)
		if strings.TrimSpace(diff) == "" {
			// Sometimes stringers obscure differences
			prettyConfig.PrintStringers = false
			diff = prettyConfig.Compare(want, got)
		}
		t.Errorf("wrong result\ndiff: %s", diff)
	}
}
