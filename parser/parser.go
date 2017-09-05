package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/apparentlymart/go-textseg/textseg"
	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/source"
	"golang.org/x/tools/godoc/vfs"
)

type Parser struct {
	fs       vfs.FileSystem
	files    map[string]*ast.File
	packages map[string]*ast.Package
	diags    source.Diags
}

func NewParser(fs vfs.FileSystem) *Parser {
	return &Parser{
		files:    map[string]*ast.File{},
		packages: map[string]*ast.Package{},
	}
}

// ParsePackage parses all of the source files in a given package.
func (p *Parser) ParsePackage(ppath string) (*ast.Package, source.Diags) {
	if pkg := p.packages[ppath]; pkg != nil {
		return pkg, nil
	}

	pkg := &ast.Package{
		Path: ppath,
	}

	entries, err := p.fs.ReadDir(ppath)
	if err != nil {
		return pkg, source.Diags{
			{
				Level:   source.Error,
				Summary: "Package not found",
				Detail:  fmt.Sprintf("The given path %q could not be resolved as a package path.", ppath),
			},
		}
	}

	p.packages[ppath] = pkg

	var diags source.Diags
	var files []*ast.File
	for _, i := range entries {
		if i.IsDir() {
			continue
		}
		name := i.Name()
		if !strings.HasSuffix(name, ".cb") {
			continue
		}

		filePath := path.Join(ppath, name)
		file, fileDiags := p.ParseFile(filePath)
		diags = append(diags, fileDiags...)
		files = append(files, file)
	}

	pkg.Files = files

	return pkg, diags
}

// ParseFile parses a single file. Most callers should use LoadPackage.
func (p *Parser) ParseFile(fpath string) (*ast.File, source.Diags) {
	if ret := p.files[fpath]; ret != nil {
		return ret, nil
	}

	ret := &ast.File{
		WithRange: ast.WithRange{
			source.Range{
				Filename: fpath,
				Start:    source.StartPos,
				End:      source.StartPos,
			},
		},
	}

	f, err := p.fs.Open(fpath)
	if err != nil {
		return ret, source.Diags{
			{
				Level:   source.Error,
				Summary: "Failed to read file",
				Detail:  fmt.Sprintf("The file %q could not be read: %s.", fpath, err),
			},
		}
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return ret, source.Diags{
			{
				Level:   source.Error,
				Summary: "Failed to read file",
				Detail:  fmt.Sprintf("The file %q could not be read: %s.", fpath, err),
			},
		}
	}

	ret.Source = src

	p.files[fpath] = ret

	tokens := scanTokens(src, "", source.StartPos, scanNormal)
	it := newTokenIterator(tokens)
	ip := &parser{
		tokenPeeker: tokenPeeker{
			Iter: it,
		},
	}

	topLevel, rng, diags := ip.ParseTopLevel()
	ret.TopLevel = topLevel
	ret.Range = rng
	p.diags = append(p.diags, diags...)

	return ret, diags
}

// ScanFile loads a single file and returns the tokens found within it.
//
// This entrypoint does not actually do any parsing, and thus doesn't produce
// diagnostics. An error is returned only if the given file cannot be read.
func (p *Parser) ScanFile(fpath string) (Tokens, error) {
	f, err := p.fs.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	tokens := scanTokens(src, fpath, source.StartPos, scanNormal)
	return tokens, nil
}

func (p *Parser) Diagnostics() source.Diags {
	return p.diags
}

// ParseExpr parses a standalone expression.
func ParseExpr(src []byte) (ast.Node, source.Diags) {
	tokens := scanTokens(src, "", source.StartPos, scanNormal)
	it := newTokenIterator(tokens)
	ip := &parser{
		tokenPeeker: tokenPeeker{
			Iter: it,
		},
	}
	return ip.ParseExpr()
}

type parser struct {
	tokenPeeker
	recovering bool
}

func (p *parser) ParseTopLevel() ([]ast.Node, source.Range, source.Diags) {
	return p.parseTopLevel()
}

func (p *parser) ParseExpr() (ast.Node, source.Diags) {
	expr, diags := p.parseExpr()

	// We tolerate leftover characters in the presence of errors because
	// we may have aborted parsing early due to the error.
	if !p.EOF() && !diags.HasErrors() {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Extra characters after expression",
			Detail:  "The remaining characters cannot be interpreted as part of this expression.",
			Ranges:  p.PeekRange().List(),
		})
	}

	return expr, diags
}

func (p *parser) parseTopLevel() ([]ast.Node, source.Range, source.Diags) {
	var ret []ast.Node
	var diags source.Diags

	startRange := p.PeekRange()
	endRange := p.PeekRange()

Statements:
	for !p.EOF() {
		var node ast.Node
		var nodeDiags source.Diags

		nextKw := p.PeekKeyword()
		switch nextKw {

		case "import":
			node, nodeDiags = p.parseImport()

		default:
			// TODO: try to parse either assignment or connection statement
			break Statements
		}

		if node != nil {
			ret = append(ret, node)
			endRange = node.SourceRange()
		}
		diags = append(diags, nodeDiags...)
	}

	rng := source.RangeBetween(startRange, endRange)

	return ret, rng, diags
}

func (p *parser) parseImport() (ast.Node, source.Diags) {
	kw := p.Read()
	if kw.Type != TokenIdent {
		// Should never happen because caller should've peeked ahead here
		panic("parseImport called with peeker not pointing at ident")
	}

	var diags source.Diags

	if p.Peek().Type != TokenStringLit {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Invalid import path",
				Detail:  "An import path must be a quoted string.",
				Ranges:  []source.Range{p.PeekRange()},
			})
		}
		p.recoverAfterSemicolon()
		return nil, diags
	}

	pathTok := p.Read()
	path, diags := p.decodeStringLiteral(pathTok)
	if diags.HasErrors() {
		return nil, diags
	}

	var name string

	if p.PeekKeyword() == "as" {
		p.Read() // eat the "as" keyword
		if p.Peek().Type != TokenIdent {
			if !p.recovering {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid import name",
					Detail:  "The name for an import must be an identifier.",
					Ranges:  []source.Range{p.PeekRange()},
				})
			}
			p.recoverAfterSemicolon()
			return nil, diags
		}

		nameTok := p.Read()
		name = p.decodeIdentifierBytes(nameTok.Bytes)
	}

	if p.Peek().Type != TokenSemicolon {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Invalid import name",
				Detail:  "The name for an import must be an identifier.",
				Ranges:  []source.Range{p.PeekRange()},
			})
		}
		p.recoverAfterSemicolon()
		return nil, diags
	}

	semicolon := p.Read()

	return &ast.Import{
		Package: path,
		Name:    name,

		PackageRange: pathTok.Range,

		WithRange: ast.WithRange{
			Range: source.RangeBetween(kw.Range, semicolon.Range),
		},
	}, diags
}

func (p *parser) parseExpr() (ast.Node, source.Diags) {
	return p.parseTernaryConditional()
}

func (p *parser) parseTernaryConditional() (ast.Node, source.Diags) {
	// TODO: implement conditional
	return p.parseBinaryOps(binaryOps)
}

// parseBinaryOps calls itself recursively to work through all of the
// operator precedence groups, and then eventually calls parseExpressionWithTrailers
// for each operand.
func (p *parser) parseBinaryOps(ops []map[TokenType]ast.ArithmeticOp) (ast.Node, source.Diags) {
	if len(ops) == 0 {
		// We've run out of operators, so now we'll just try to parse a term.
		return p.parseExpressionWithTrailers()
	}

	thisLevel := ops[0]
	remaining := ops[1:]

	var lhs, rhs ast.Node
	var operation ast.ArithmeticOp
	var diags source.Diags

	// Parse a term that might be the first operand of a binary
	// operation or it might just be a standalone term.
	// We won't know until we've parsed it and can look ahead
	// to see if there's an operator token for this level.
	lhs, lhsDiags := p.parseBinaryOps(remaining)
	diags = append(diags, lhsDiags...)
	if p.recovering && lhsDiags.HasErrors() {
		return lhs, diags
	}

	// We'll keep eating up operators until we run out, so that operators
	// with the same precedence will combine in a left-associative manner:
	// a+b+c => (a+b)+c, not a+(b+c)
	//
	// Should we later want to have right-associative operators, a way
	// to achieve that would be to call back up to ParseExpression here
	// instead of iteratively parsing only the remaining operators.
	for {
		next := p.Peek()
		var newOp ast.ArithmeticOp
		var ok bool
		if newOp, ok = thisLevel[next.Type]; !ok {
			break
		}

		// Are we extending an expression started on the previous iteration?
		if operation != ast.ArithmeticOpNil {
			lhs = &ast.ArithmeticBinary{
				LHS: lhs,
				Op:  operation,
				RHS: rhs,

				WithRange: ast.WithRange{
					Range: source.RangeBetween(lhs.SourceRange(), rhs.SourceRange()),
				},
			}
		}

		operation = newOp
		p.Read() // eat operator token
		var rhsDiags source.Diags
		rhs, rhsDiags = p.parseBinaryOps(remaining)
		diags = append(diags, rhsDiags...)
		if p.recovering && rhsDiags.HasErrors() {
			return lhs, diags
		}
	}

	if operation == ast.ArithmeticOpNil {
		return lhs, diags
	}

	return &ast.ArithmeticBinary{
		LHS: lhs,
		Op:  operation,
		RHS: rhs,

		WithRange: ast.WithRange{
			Range: source.RangeBetween(lhs.SourceRange(), rhs.SourceRange()),
		},
	}, diags
}

func (p *parser) parseExpressionWithTrailers() (ast.Node, source.Diags) {
	term, diags := p.parseExpressionTerm()

	// TODO: actually parse trailers (GetAttr, GetIndex, Call)

	return term, diags
}

func (p *parser) parseExpressionTerm() (ast.Node, source.Diags) {
	start := p.Peek()

	switch start.Type {
	case TokenOParen:
		open := p.Read()

		expr, diags := p.parseExpr()
		close := p.Peek()
		if close.Type != TokenCParen {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Unbalanced parentheses",
				Detail:  "Expected a closing parenthesis to terminate the expression.",
				Ranges:  close.Range.List(),
			})
		}

		if diags.HasErrors() {
			// attempt to place the peeker after our closing paren
			// before we return, so that the next parser has some
			// chance of finding a valid expression.
			p.recoverAfterClose(TokenCParen)
			return expr, diags
		}

		p.Read() // eat closing paren

		return &ast.ParenExpr{
			WithRange: ast.WithRange{
				Range: source.RangeBetween(open.Range, close.Range),
			},
			Content: expr,
		}, diags

	case TokenIdent:
		kw := p.PeekKeyword()
		tok := p.Read()

		switch kw {
		case "true":
			return &ast.BooleanLit{
				Value: true,
				WithRange: ast.WithRange{
					Range: tok.Range,
				},
			}, nil
		case "false":
			return &ast.BooleanLit{
				Value: false,
				WithRange: ast.WithRange{
					Range: tok.Range,
				},
			}, nil
		default:
			return &ast.Variable{
				Name: p.decodeIdentifierBytes(tok.Bytes),
				WithRange: ast.WithRange{
					Range: tok.Range,
				},
			}, nil
		}

	case TokenStringLit:
		tok := p.Read()
		val, diags := p.decodeStringLiteral(tok)

		return &ast.StringLit{
			WithRange: ast.WithRange{
				Range: tok.Range,
			},
			Value: val,
		}, diags

	default:
		var diags source.Diags
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Invalid expression",
				Detail:  "Expected the start of an expression, but found invalid characters.",
				Ranges:  start.Range.List(),
			})
		}
		p.setRecovering()

		return &ast.Invalid{
			WithRange: ast.WithRange{
				Range: start.Range,
			},
		}, diags
	}
}

func (p *parser) decodeIdentifierBytes(src []byte) string {
	if len(src) == 0 {
		// should never happen, but we'll catch it to avoid a panic below
		return ""
	}

	if src[0] == '`' {
		// Trim off the leading and trailing ` characters that quote the sequence
		src = src[1 : len(src)-1]
	}

	return string(src)
}

func (p *parser) decodeStringLiteral(tok Token) (string, source.Diags) {
	var quoted bool
	src := tok.Bytes
	switch tok.Type {
	case TokenStringLit:
		quoted = true
		src = src[1 : len(src)-1]
	default:
		panic("decodeStringLiteral can only be used with TokenStringLit tokens")
	}
	var diags source.Diags

	ret := make([]byte, 0, len(src))
	var esc []byte

	sc := bufio.NewScanner(bytes.NewReader(src))
	sc.Split(textseg.ScanGraphemeClusters)

	pos := tok.Range.Start
	newPos := pos
Character:
	for sc.Scan() {
		pos = newPos
		ch := sc.Bytes()

		// Adjust position based on our new character.
		// \r\n is considered to be a single character in text segmentation,
		if (len(ch) == 1 && ch[0] == '\n') || (len(ch) == 2 && ch[1] == '\n') {
			newPos.Line++
			newPos.Column = 0
		} else {
			newPos.Column++
		}
		newPos.Byte += len(ch)

		if len(esc) > 0 {
			switch esc[0] {
			case '\\':
				if len(ch) == 1 {
					switch ch[0] {

					// TODO: numeric character escapes with \uXXXX

					case 'n':
						ret = append(ret, '\n')
						esc = esc[:0]
						continue Character
					case 'r':
						ret = append(ret, '\r')
						esc = esc[:0]
						continue Character
					case 't':
						ret = append(ret, '\t')
						esc = esc[:0]
						continue Character
					case '"':
						ret = append(ret, '"')
						esc = esc[:0]
						continue Character
					case '\\':
						ret = append(ret, '\\')
						esc = esc[:0]
						continue Character
					}
				}

				var detail string
				switch {
				case len(ch) == 1 && (ch[0] == '$' || ch[0] == '!'):
					detail = fmt.Sprintf(
						"The characters \"\\%s\" do not form a recognized escape sequence. To escape a \"%s{\" template sequence, use \"%s%s{\".",
						ch, ch, ch, ch,
					)
				default:
					detail = fmt.Sprintf("The characters \"\\%s\" do not form a recognized escape sequence.", ch)
				}

				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid escape sequence",
					Detail:  detail,
					Ranges: []source.Range{
						{
							Filename: tok.Range.Filename,
							Start: source.Pos{
								Line:   pos.Line,
								Column: pos.Column - 1, // safe because we know the previous character must be a backslash
								Byte:   pos.Byte - 1,
							},
							End: source.Pos{
								Line:   pos.Line,
								Column: pos.Column + 1, // safe because we know the previous character must be a backslash
								Byte:   pos.Byte + len(ch),
							},
						},
					},
				})
				ret = append(ret, ch...)
				esc = esc[:0]
				continue Character

			case '$', '!':
				switch len(esc) {
				case 1:
					if len(ch) == 1 && ch[0] == esc[0] {
						esc = append(esc, ch[0])
						continue Character
					}

					// Any other character means this wasn't an escape sequence
					// after all.
					ret = append(ret, esc...)
					ret = append(ret, ch...)
					esc = esc[:0]
				case 2:
					if len(ch) == 1 && ch[0] == '{' {
						// successful escape sequence
						ret = append(ret, esc[0])
					} else {
						// not an escape sequence, so just output literal
						ret = append(ret, esc...)
					}
					ret = append(ret, ch...)
					esc = esc[:0]
				default:
					// should never happen
					panic("have invalid escape sequence >2 characters")
				}

			}
		} else {
			if len(ch) == 1 {
				switch ch[0] {
				case '\\':
					if quoted { // ignore backslashes in unquoted mode
						esc = append(esc, '\\')
						continue Character
					}
				case '$':
					esc = append(esc, '$')
					continue Character
				case '!':
					esc = append(esc, '!')
					continue Character
				}
			}
			ret = append(ret, ch...)
		}
	}

	return string(ret), diags
}

func (p *parser) setRecovering() {
	p.recovering = true
}

// recoverAfterClose seeks forward in the token stream until it finds TokenType
// "end", then returns with the peeker pointed at the following token.
//
// If the given token type is a bracketer, this function will additionally
// count nested instances of the brackets to try to leave the peeker at
// the end of the _current_ instance of that bracketer, skipping over any
// nested instances. This is a best-effort operation and may have
// unpredictable results on input with bad bracketer nesting.
func (p *parser) recoverAfterClose(end TokenType) Token {
	start := p.oppositeBracket(end)
	p.recovering = true

	nest := 0
	for {
		tok := p.Read()
		ty := tok.Type

		switch ty {
		case start:
			nest++
		case end:
			if nest < 1 {
				return tok
			}

			nest--
		case TokenEOF:
			return tok
		}
	}
}

// recoverAfterBlock sets the recovery flag and then tries to place the
// peeker just after the brace that closes the current block.
//
// This should be called when the peeker has already read the opening
// brace for the current block. If the peeker is at or before the brace,
// use recoverAfterNextBlock.
//
// Recovery is not an exact science, so the peeker may be left in a strange
// place that will lead to more errors. The recovery flag should be used to
// suppress "invalid token"-type errors and abort early to reduce the risk
// of reporting a chain of compounding errors to the user.
func (p *parser) recoverAfterCurrentBlock() Token {
	return p.recoverAfterClose(TokenCBrace)
}

// recoverAfterNextBlock is like recoverAfterCurrentBlock except that it
// first seeks forward to locate the next opening brace, and then places
// the peeker after its corresponding closing brace.
func (p *parser) recoverAfterNextBlock() Token {
	for p.Peek().Type != TokenOBrace && p.Peek().Type != TokenEOF {
		p.Read()
	}

	// After we've located an open brace or an EOF, seek forward one more
	// time so we're placed _after_ the brace. (If it was EOF, we'll just
	// get the same EOF again.)
	p.Read()

	return p.recoverAfterCurrentBlock()
}

// recoverAfterSemicolon sets the recovery flag and then tries to place the
// peeker just after the semicolon that closes the current statement.
//
// Recovery is not an exact science, so the peeker may be left in a strange
// place that will lead to more errors. The recovery flag should be used to
// suppress "invalid token"-type errors and abort early to reduce the risk
// of reporting a chain of compounding errors to the user.
func (p *parser) recoverAfterSemicolon() {
	p.recovering = true
	braceCount := 0

	for {
		next := p.Read()

		switch next.Type {
		case TokenEOF:
			return
		case TokenOBrace:
			braceCount++
		case TokenCBrace:
			braceCount--
		case TokenSemicolon:
			// Only semicolons that are not inside braces are considered,
			// since we don't want to stop too early if there's a nested
			// set of braces (e.g. an object expression) in our path to the
			// semicolon.
			if braceCount <= 0 {
				return
			}
		}
	}
}

// oppositeBracket finds the bracket that opposes the given bracketer, or
// NilToken if the given token isn't a bracketer.
//
// "Bracketer", for the sake of this function, is one end of a matching
// open/close set of tokens that establish a bracketing context.
func (p *parser) oppositeBracket(ty TokenType) TokenType {
	switch ty {

	case TokenOBrace:
		return TokenCBrace
	case TokenOBrack:
		return TokenCBrack
	case TokenOParen:
		return TokenCParen
	case TokenOPoint:
		return TokenCPoint

	case TokenCBrace:
		return TokenOBrace
	case TokenCBrack:
		return TokenOBrack
	case TokenCParen:
		return TokenOParen
	case TokenCPoint:
		return TokenOPoint

	default:
		return TokenNil
	}
}

// invalidTokenDiags takes a source.Diags and returns it only if the parser
// is not in recovery mode. If it _is_ in recovery mode, nil is returned.
//
// We tend to skip returning "invalid token"-type messages when in recovery
// mode because we want to avoid returning many cascading failures in the
// presence of a missing token but yet still parse as much of the file as
// we can manage.
func (p *parser) invalidTokenDiags(diags source.Diags) source.Diags {
	if p.recovering {
		return nil
	}
	return diags
}
