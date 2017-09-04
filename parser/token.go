package parser

import (
	"fmt"

	"github.com/apparentlymart/go-textseg/textseg"
	"github.com/cirbo-lang/cirbo/source"
)

// Token represents a sequence of bytes from some zcl code that has been
// tagged with a type and its range within the source file.
type Token struct {
	Type  TokenType
	Bytes []byte
	Range source.Range
}

// Tokens is a slice of Token.
type Tokens []Token

// TokenType is an enumeration used for the Type field on Token.
type TokenType rune

const (
	// Single-character tokens are represented by their own character, for
	// convenience in producing these within the scanner. However, the values
	// are otherwise arbitrary and just intended to be mnemonic for humans
	// who might see them in debug output.

	TokenOBrace TokenType = '{'
	TokenCBrace TokenType = '}'
	TokenOBrack TokenType = '['
	TokenCBrack TokenType = ']'
	TokenOParen TokenType = '('
	TokenCParen TokenType = ')'

	TokenStar    TokenType = '*'
	TokenSlash   TokenType = '/'
	TokenPlus    TokenType = '+'
	TokenMinus   TokenType = '-'
	TokenPercent TokenType = '%'
	TokenCaret   TokenType = '^'

	TokenAssign        TokenType = '='
	TokenEqual         TokenType = '≔'
	TokenNotEqual      TokenType = '≠'
	TokenLessThan      TokenType = '<'
	TokenLessThanEq    TokenType = '≤'
	TokenGreaterThan   TokenType = '>'
	TokenGreaterThanEq TokenType = '≥'

	TokenDashDash TokenType = '—'
	TokenDotDot   TokenType = '…'

	TokenAnd  TokenType = '∧'
	TokenOr   TokenType = '∨'
	TokenBang TokenType = '!'

	TokenDot   TokenType = '.'
	TokenComma TokenType = ','

	TokenQuestion  TokenType = '?'
	TokenColon     TokenType = ':'
	TokenSemicolon TokenType = ';'

	TokenStringLit TokenType = 'S'
	TokenNumberLit TokenType = 'N'
	TokenIdent     TokenType = 'I'

	TokenComment    TokenType = 'C'
	TokenWhitespace TokenType = ' '

	TokenEOF TokenType = '␄'

	// The following are not used in the language but recognized by the scanner
	// so we can generate good diagnostics in the parser when users try to write
	// things that might work in other languages they are familiar with, or
	// simply make incorrect assumptions about the cirbo language.

	TokenBitwiseAnd TokenType = '&'
	TokenBitwiseOr  TokenType = '|'
	TokenBitwiseNot TokenType = '~'
	TokenInvalid    TokenType = '�'
	TokenBadUTF8    TokenType = '💩'

	// TokenNil is a placeholder for when a token is required but none is
	// available, e.g. when reporting errors. The scanner will never produce
	// this as part of a token stream.
	TokenNil TokenType = '\x00'
)

func (t TokenType) GoString() string {
	return fmt.Sprintf("parser.%s", t.String())
}

type scanMode int

const (
	scanNormal scanMode = iota
	scanTemplate
)

type tokenAccum struct {
	Filename string
	Bytes    []byte
	Pos      source.Pos
	Tokens   []Token
}

func (f *tokenAccum) emitToken(ty TokenType, startOfs, endOfs int) {
	// Walk through our buffer to figure out how much we need to adjust
	// the start pos to get our end pos.

	start := f.Pos
	start.Column += startOfs - f.Pos.Byte // Safe because only ASCII whitespace can be in the offset (we count a tab character as one column)
	start.Byte = startOfs

	end := start
	end.Byte = endOfs
	b := f.Bytes[startOfs:endOfs]
	for len(b) > 0 {
		advance, seq, _ := textseg.ScanGraphemeClusters(b, true)
		if len(seq) == 1 && seq[0] == '\n' {
			end.Line++
			end.Column = 1
		} else {
			end.Column++
		}
		b = b[advance:]
	}

	f.Pos = end

	f.Tokens = append(f.Tokens, Token{
		Type:  ty,
		Bytes: f.Bytes[startOfs:endOfs],
		Range: source.Range{
			Filename: f.Filename,
			Start:    start,
			End:      end,
		},
	})
}

type heredocInProgress struct {
	Marker      []byte
	StartOfLine bool
}
