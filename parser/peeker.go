package parser

import "github.com/cirbo-lang/cirbo/source"

type tokenPeeker struct {
	Iter   *tokenIterator
	Peeked Token
}

func (p *tokenPeeker) Peek() Token {
	for p.Peeked.Type == TokenNil || p.Peeked.Type == TokenWhitespace || p.Peeked.Type == TokenComment {
		p.Peeked = p.Iter.Next()
	}
	return p.Peeked
}

func (p *tokenPeeker) Read() Token {
	ret := p.Peek()
	p.Peeked.Type = TokenNil
	return ret
}

// PeekIdent checks if the next token is an unquoted identifier, and if so it
// returns the identifier's name as a string. If the next token is not an
// identifier, or if it is quoted with backticks, it returns an empty string.
func (p *tokenPeeker) PeekKeyword() string {
	next := p.Peek()
	if next.Type != TokenIdent {
		return ""
	}
	got := string(next.Bytes)
	if got[0] == '`' {
		// this is a `quoted` ident, which means it can never be a keyword
		// (user can write e.g. `import` to escape the special meaning of
		// the import keyword.)
		return ""
	}
	return got
}

func (p *tokenPeeker) PeekRange() source.Range {
	return p.Peek().Range
}

func (p *tokenPeeker) EOF() bool {
	return p.Peek().Type == TokenEOF
}

type tokenIterator struct {
	Tokens Tokens
	Pos    int
}

func newTokenIterator(tokens Tokens) *tokenIterator {
	return &tokenIterator{
		Tokens: tokens,
		Pos:    0,
	}
}

func (i *tokenIterator) Next() Token {
	ret := i.Tokens[i.Pos]
	if i.Pos < (len(i.Tokens) - 1) {
		// When we reach the end we will just keep returning the final token
		// forever, since we assume it will be a TokenEOF.
		i.Pos++
	}
	return ret
}
