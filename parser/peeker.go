package parser

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
