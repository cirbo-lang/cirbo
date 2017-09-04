package parser

type tokenPeeker struct {
	Iter   *tokenIterator
	Peeked Token
}

func (p *tokenPeeker) Peek() Token {
	if p.Peeked.Type == TokenNil {
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
		i.Pos++
	}
	return ret
}
