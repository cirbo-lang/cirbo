package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"path"
	"strings"

	"github.com/apparentlymart/go-textseg/textseg"
	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/source"
	"golang.org/x/tools/godoc/vfs"
)

var oneHundred = mustParseBigFloat("100")

type Parser struct {
	fs       vfs.FileSystem
	files    map[string]*ast.File
	packages map[string]*ast.Package
	diags    source.Diags
}

// NewParser creates a new parser that works with files under the given
// project root directory.
//
// A project directory should usually contain one or more project files
// (with the ".cb" filename extension) and may contain subdirectories
// that represent project-local packages. It may also contain a directory
// called "cirbo-pkg" which contains local copies of third-party packages,
// usually managed with cirbo's built-in package management tools.
//
// Internally, the given directory is used as the root of a virtual filesystem
// and all file paths are rooted in that filesystem. These project-relative
// paths appear, in particular, in returned diagnostics. Before displaying such
// paths to an end-user a caller should re-interpret these paths relative to
// the "real" filesystem to avoid confusion.
func NewParser(projectRoot string) *Parser {
	fs := vfs.OS(projectRoot)

	return &Parser{
		fs:       fs,
		files:    map[string]*ast.File{},
		packages: map[string]*ast.Package{},
	}
}

// ParsePackage parses all of the source files in a given package.
//
// If "from" is non-empty then it is a project-root-relative path to the file
// that is requesting this package, which enables the use of relative package
// paths.
func (p *Parser) ParsePackage(ppath string, from string) (*ast.Package, source.Diags) {
	pkg := &ast.Package{
		DefaultName: path.Base(ppath),
	}

	vfsPath := p.resolvePackagePath(ppath, from)
	if vfsPath == "" {
		return pkg, source.Diags{
			{
				Level:   source.Error,
				Summary: "Package not found",
				Detail:  fmt.Sprintf("The given path %q could not be resolved as a package path.", ppath),
			},
		}
	}

	if pkg := p.packages[vfsPath]; pkg != nil {
		return pkg, nil
	}

	entries, err := p.fs.ReadDir(vfsPath)
	if err != nil {
		return pkg, source.Diags{
			{
				Level:   source.Error,
				Summary: "Package not found",
				Detail:  fmt.Sprintf("The given path %q could not be resolved as a package path.", ppath),
			},
		}
	}

	p.packages[vfsPath] = pkg

	var diags source.Diags
	var files []*ast.File
	for _, i := range entries {
		if i.IsDir() {
			continue
		}
		name := i.Name()
		if !pathHasExtension(name, ".cbm") {
			continue
		}

		filePath := path.Join(vfsPath, name)
		file, fileDiags := p.ParseFile(filePath)
		diags = append(diags, fileDiags...)
		files = append(files, file)
	}

	pkg.Files = files

	return pkg, diags
}

// ParseAllProjectFiles finds all of the project files in the project root and
// parses them, returning a slice containing one entry for each.
//
// Most normal operations operate on a single file at a time, provided by the
// user on the command line. This method is provided for the rare commands that
// operate on an entire project directory, such as the package installer when
// it's looking for dependencies in all project files.
func (p *Parser) ParseAllProjectFiles() ([]*ast.File, source.Diags) {
	var files []*ast.File

	entries, err := p.fs.ReadDir("/")
	if err != nil {
		return files, source.Diags{
			{
				Level:   source.Error,
				Summary: "Invalid project root",
				Detail:  fmt.Sprintf("Failed to read from the project root: %s", err),
			},
		}
	}

	var diags source.Diags
	for _, i := range entries {
		if i.IsDir() {
			continue
		}
		name := i.Name()
		if !pathHasExtension(name, ".cb") {
			continue
		}

		filePath := path.Join("/", name)
		file, fileDiags := p.ParseFile(filePath)
		diags = append(diags, fileDiags...)
		files = append(files, file)
	}

	return files, diags
}

// ParseFile parses a single file.
//
// This is usually used for project files. Module files can be loaded via this
// method for purposes such as single-file validation and text editor support
// tools, but in normal use modules should be loaded as part of their packages
// using ParsePackage.
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

	tokens := scanTokens(src, fpath, source.StartPos, scanNormal)
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

func (p *Parser) resolvePackagePath(ppath string, from string) string {
	if len(ppath) == 0 {
		return ""
	}

	parts := strings.Split(ppath, "/")
	for i, part := range parts {
		if part == "" {
			// empty segments are never valid
			return ""
		}

		if i > 0 {
			if part == "." || part == ".." {
				// relative references only allowed in first segment
				return ""
			}
		}
	}

	if parts[0] == "." || parts[0] == ".." {
		// relative to requesting file

		if from == "" {
			// no requesting file, so relative references are not allowed
			return ""
		}

		// must always have at least one additional part to traverse after
		// the relative.
		if len(parts) < 2 {
			return ""
		}

		return path.Join(path.Dir(from), ppath)
	}

	return path.Join("/cirbo-pkg", ppath)
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
	return p.parseStmts(TokenEOF)
}

func (p *parser) parseStmts(endType TokenType) ([]ast.Node, source.Range, source.Diags) {
	var ret []ast.Node
	var diags source.Diags

	if p.Peek().Type == endType {
		// With an empty body we can't really produce a real range, so we'll
		// make a zero-length range that sits just before the next token.
		rng := p.PeekRange()
		rng.End = rng.Start
		return ret, rng, diags
	}

	startRange := p.PeekRange()
	endRange := p.PeekRange()

Statements:
	for p.Peek().Type != endType {
		if p.Peek().Type == TokenEOF {
			// It's the caller's responsibility to detect unclosed statement
			// blocks if the end type isn't TokenEOF
			break Statements
		}

		var node ast.Node
		var nodeDiags source.Diags

		nextKw := p.PeekKeyword()
		switch nextKw {

		case "import":
			node, nodeDiags = p.parseImport()

		case "board":
			node, nodeDiags = p.parseBoard()

		case "circuit":
			node, nodeDiags = p.parseCircuit()

		case "device":
			node, nodeDiags = p.parseDevice()

		case "land":
			node, nodeDiags = p.parseLand()

		case "pinout":
			node, nodeDiags = p.parsePinout()

		default:

			if p.keywordCanStartTerminalDecl(nextKw) {
				node, nodeDiags = p.parseTerminalDecl()
				break
			}

			if p.Peek().Type == TokenSemicolon {
				p.Read()
				continue Statements
			}

			// If not a keyword, then we should have either an assignment or
			// a connection statement.
			node, nodeDiags = p.parseAssignOrConnectStmt()
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

func (p *parser) parseStmtBlock() (*ast.StatementBlock, source.Diags) {
	var ret []ast.Node
	var diags source.Diags

	if p.Peek().Type != TokenOBrace {
		bad := p.Peek()
		if !p.recovering {
			switch p.Peek().Type {
			case TokenSemicolon:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Missing block",
					Detail:  "A brace-delimited statement block is required here.",
					Ranges:  []source.Range{p.PeekRange()},
				})
				p.recoverAfterSemicolon()
			default:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid block",
					Detail:  "A brace-delimited statement block is required here.",
					Ranges:  []source.Range{p.PeekRange()},
				})
				p.recoverAfterNextBlock()
			}
		}
		return &ast.StatementBlock{
			Statements: ret,
			WithRange: ast.WithRange{
				Range: bad.Range,
			},
		}, diags
	}

	open := p.Read() // eat {

	var rng source.Range
	ret, rng, diags = p.parseStmts(TokenCBrace)

	if p.Peek().Type != TokenCBrace {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Unclosed statement block",
				Detail:  "This statement block is not closed before the end of the file. Ensure that each opening brace \"{\" has a corresponding closing brace \"}\".",
				Ranges:  open.Range.List(),
			})
			// The only way we can get here is if we're at EOF, so calling
			// this is a little silly but it at least activates recovery mode
			// to inhibit any other errors that might otherwise be emitted
			// as we unwind our parsing stack.
			p.recoverAfterCurrentBlock()
		}
		return &ast.StatementBlock{
			Statements: ret,
			WithRange: ast.WithRange{
				Range: source.RangeBetween(open.Range, rng),
			},
		}, diags
	}

	close := p.Read() // eat }

	return &ast.StatementBlock{
		Statements: ret,
		WithRange: ast.WithRange{
			Range: source.RangeBetween(open.Range, close.Range),
		},
	}, diags
}

func (p *parser) parseAssignOrConnectStmt() (ast.Node, source.Diags) {
	switch p.Peek().Type {

	case TokenBarDashDash:
		// Left-pointing "not connected" symbol
		return p.parseConnectStmt(nil)

	default:
		expr, exprDiags := p.parseExpr()
		if exprDiags.HasErrors() {
			p.recoverAfterSemicolon()
			return &ast.Invalid{
				WithRange: ast.WithRange{
					Range: expr.SourceRange(),
				},
			}, exprDiags
		}

		// The token following the expression defines what kind of
		// expression it is.
		switch p.Peek().Type {
		case TokenDashDash, TokenDashDashBar:
			return p.parseConnectStmt(expr)
		case TokenAssign:
			return p.parseAssignStmt(expr)
		case TokenSemicolon:
			p.Read() // eat semicolon
			return expr, source.Diags{
				{
					Level:   source.Error,
					Summary: "Useless naked expression",
					Detail:  "An expression alone is not a valid statement.",
					Ranges:  expr.SourceRange().List(),
				},
			}
		default:
			var diags source.Diags
			if !p.recovering {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid characters after expression",
					Detail:  "The remaining characters cannot be interpreted as part of this expression.",
					Ranges:  p.PeekRange().List(),
				})
			}
			p.recoverAfterSemicolon()
			return expr, diags
		}
	}

}

func (p *parser) parseAssignStmt(lvalue ast.Node) (ast.Node, source.Diags) {
	var diags source.Diags

	// If the caller didn't already parse our lvalue expr, we'll do it
	// right here.
	if lvalue == nil {
		lvalue, diags = p.parseExpr()
		if diags.HasErrors() {
			p.recoverAfterSemicolon()
			return lvalue, diags
		}
	}

	var name string
	if varExpr, isVar := lvalue.(*ast.Variable); isVar {
		name = varExpr.Name
	} else {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Invalid assignment expression",
			Detail:  "Can only assign directly to a variable name.",
			Ranges:  lvalue.SourceRange().List(),
		})
	}

	if p.Peek().Type != TokenAssign {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Invalid characters after expression",
				Detail:  "Expected an equals (\"=\") symbol to assign a value.",
				Ranges:  p.PeekRange().List(),
			})
		}
		wrong := p.Peek()
		p.recoverAfterSemicolon()
		return &ast.Assign{
			Name: name,
			Value: &ast.Invalid{
				WithRange: ast.WithRange{
					Range: wrong.Range,
				},
			},
			WithRange: ast.WithRange{
				Range: source.RangeBetween(lvalue.SourceRange(), p.Peek().Range),
			},
		}, diags
	}

	p.Read() // eat assignment equals

	rhs, rhsDiags := p.parseExpr()
	diags = append(diags, rhsDiags...)

	var end Token
	if p.Peek().Type != TokenSemicolon {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Unterminated statement",
				Detail:  "This assignment statement must be terminated by a semicolon.",
				Ranges:  source.RangeBetween(lvalue.SourceRange(), rhs.SourceRange()).List(),
			})
		}
		end = p.Peek()
		p.recoverAfterSemicolon()
	} else {
		end = p.Read() // eat semicolon
	}

	return &ast.Assign{
		Name:  name,
		Value: rhs,

		WithRange: ast.WithRange{
			Range: source.RangeBetween(lvalue.SourceRange(), end.Range),
		},
	}, diags
}

func (p *parser) parseConnectStmt(first ast.Node) (ast.Node, source.Diags) {
	var diags source.Diags

	if p.Peek().Type == TokenBarDashDash {
		if first != nil {
			// this indicates a bug in the caller
			panic("parseConnectStmt start must be nil when next token is |--")
		}

		start := p.Read() // eat |-- token

		var expr ast.Node
		expr, diags = p.parseExpr()
		var end Token

		if p.Peek().Type != TokenSemicolon {
			if !p.recovering {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Unterminated statement",
					Detail:  "This connection statement must be terminated by a semicolon.",
					Ranges:  source.RangeBetween(start.Range, expr.SourceRange()).List(),
				})
			}
			end = p.Peek()
		} else {
			end = p.Read() // eat semicolon
		}

		if diags.HasErrors() {
			p.recoverAfterSemicolon()
		}

		return &ast.NoConnection{
			Terminal: expr,

			WithRange: ast.WithRange{
				Range: source.RangeBetween(start.Range, end.Range),
			},
		}, diags
	}

	if first == nil {
		// Only expect an initial expression if the caller didn't already
		// pass one in. (Callers may need to look past an expression before
		// they know they've found a connect statement.
		first, diags = p.parseExpr()

		if diags.HasErrors() {
			p.recoverAfterSemicolon()
			return &ast.Connection{
				Seq: nil,

				WithRange: ast.WithRange{
					Range: first.SourceRange(),
				},
			}, diags
		}
	}

	if p.Peek().Type == TokenDashDashBar {
		p.Read() // eat --| token

		var end Token
		if p.Peek().Type != TokenSemicolon {
			if !p.recovering {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Unterminated statement",
					Detail:  "This connection statement must be terminated by a semicolon.",
					Ranges:  []source.Range{p.PeekRange()},
				})
			}
			end = p.Peek()
			p.recoverAfterSemicolon()
		} else {
			end = p.Read() // eat semicolon
		}

		return &ast.NoConnection{
			Terminal: first,

			WithRange: ast.WithRange{
				Range: source.RangeBetween(first.SourceRange(), end.Range),
			},
		}, diags
	}

	var terminals []ast.Node
	terminals = append(terminals, first)
	for p.Peek().Type == TokenDashDash {
		p.Read() // eat "--" token

		if p.Peek().Type == TokenSemicolon {
			if !p.recovering {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Missing terminal expression",
					Detail:  "The \"--\" symbol must be followed by an expression for a connection terminal.",
					Ranges:  []source.Range{p.PeekRange()},
				})
			}
			break
		}

		expr, exprDiags := p.parseExpr()
		diags = append(diags, exprDiags...)
		terminals = append(terminals, expr)
	}

	var end Token
	if p.Peek().Type != TokenSemicolon {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Unterminated statement",
				Detail:  "This connection statement must be terminated by a semicolon.",
				Ranges:  []source.Range{p.PeekRange()},
			})
		}
		end = p.Peek()
		p.recoverAfterSemicolon()
	} else {
		end = p.Read() // eat semicolon
	}

	return &ast.Connection{
		Seq: terminals,

		WithRange: ast.WithRange{
			Range: source.RangeBetween(first.SourceRange(), end.Range),
		},
	}, diags
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
				Summary: "Unterminated statement",
				Detail:  "This import statement must be terminated by a semicolon.",
				Ranges:  source.RangeBetween(kw.Range, p.Peek().Range).List(),
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

func (p *parser) keywordCanStartTerminalDecl(kw string) bool {
	switch kw {
	case "terminal", "power", "input", "output", "bidi":
		return true
	default:
		return false
	}
}

func (p *parser) parseTerminalDecl() (ast.Node, source.Diags) {
	kw := p.PeekKeyword()
	if !p.keywordCanStartTerminalDecl(kw) {
		// Indicates a bug in the caller
		panic("parseTerminalDecl called with peeker not pointing at valid keyword")
	}

	first := p.Peek()
	last := first
	terminal := &ast.Terminal{}
	var diags source.Diags

	type DeclState int
	const (
		begin DeclState = iota
		afterPower
		afterBidi
		afterOutput
		optionalRole
		end
	)

	state := begin

Keywords:
	for state != end {
		kw = p.PeekKeyword()

		switch state {
		case begin:
			switch kw {
			case "terminal":
				terminal.Type = ast.Passive
				terminal.Dir = ast.Undirected
				state = end
			case "power":
				terminal.Type = ast.Power
				terminal.Dir = ast.Undirected
				state = afterPower
			case "input":
				terminal.Type = ast.Signal
				terminal.Dir = ast.Input
				state = optionalRole
			case "output":
				terminal.Type = ast.Signal
				terminal.Dir = ast.Output
				state = afterOutput
			case "bidi":
				terminal.Type = ast.Signal
				terminal.Dir = ast.Bidirectional
				state = afterBidi
			default:
				// Should never happen since the above should be exhaustive
				// of all the start keywords.
				panic("invalid initial keyword in terminal declaration")
			}
		case afterPower:
			switch kw {
			case "input":
				terminal.Dir = ast.Input
				state = optionalRole
			case "output":
				terminal.Dir = ast.Output
				state = afterOutput
			case "bidi":
				state = afterBidi
			default:
				break Keywords
			}
		case afterBidi:
			switch kw {
			case "leader":
				terminal.Role = ast.Leader
				state = end
			case "follower":
				terminal.Role = ast.Follower
				state = end
			default:
				// bidi must always be followed by a role keyword
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid bi-directional terminal role",
					Detail:  "The keyword \"bidi\" must be followed by either \"leader\" or \"follower\" to define the terminal's role relative to other terminals on its net.",
					Ranges:  p.PeekRange().List(),
				})
				p.setRecovering()

				// Placeholder value
				terminal.Dir = ast.Undirected
				break Keywords
			}
		case afterOutput:
			switch kw {
			case "emitter":
				terminal.OutputType = ast.OpenEmitter
				state = optionalRole
			case "collector":
				terminal.OutputType = ast.OpenCollector
				state = optionalRole
			case "tristate":
				terminal.OutputType = ast.Tristate
				state = optionalRole
			case "leader":
				terminal.Role = ast.Leader
				state = end
			case "follower":
				terminal.Role = ast.Follower
				state = end
			default:
				terminal.OutputType = ast.PushPull
				break Keywords
			}
		case optionalRole:
			switch kw {
			case "leader":
				terminal.Role = ast.Leader
				state = end
			case "follower":
				terminal.Role = ast.Follower
				state = end
			default:
				break Keywords
			}
		}

		// If we get down here without breaking out early, we consume a keyword.
		last = p.Read()
	}

	// The range so far: all of the keywords we visited in the state machine above
	terminal.WithRange.Range = source.RangeBetween(first.Range, last.Range)

	// If the next thing could potentially be an extraneous terminal type
	// keyword then we'll make a note of it and produce a different error
	// message below if the statement doesn't close as we expect.
	// (This is just to get a more helpful error message if the user
	// uses an inappropriate combination of keywords.)
	nextKw := p.PeekKeyword()
	var extraKw string
	var extraKwTok Token
	switch nextKw {
	case "terminal", "power", "input", "output", "bidi", "leader", "follower", "emitter", "collector", "tristate":
		extraKw = nextKw
		extraKwTok = p.Peek()
	}

	if p.Peek().Type != TokenIdent {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Missing or invalid terminal name",
				Detail:  "A terminal declaration must conclude with an identifer that names the terminal.",
				Ranges:  p.PeekRange().List(),
			})
			p.setRecovering()
		}
		p.recoverAfterSemicolon()
		return terminal, diags
	}

	nameTok := p.Read()
	terminal.Name = p.decodeIdentifierBytes(nameTok.Bytes)
	terminal.WithRange.Range = source.RangeBetween(first.Range, nameTok.Range)

	if p.Peek().Type != TokenSemicolon {
		if !p.recovering {
			// If the identifier we read above smelled like it could be an
			// extraneous decl keyword then we'll assume the user tried for
			// an invalid combination of keywords and produce a different
			// error message as a result.
			switch extraKw {
			case "terminal":
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid terminal declaration",
					Detail:  "The keyword \"terminal\" may not be used in conjunction with other terminal definition keywords.",
					Ranges:  extraKwTok.Range.List(),
				})
			case "power":
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid terminal declaration",
					Detail:  "The keyword \"power\" must always be the first keyword in a terminal definition.",
					Ranges:  extraKwTok.Range.List(),
				})
			case "leader", "follower":
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid terminal declaration",
					Detail:  fmt.Sprintf("The keyword %q must be the final terminal definition keyword.", extraKw),
					Ranges:  extraKwTok.Range.List(),
				})
			case "emitter", "collector", "tristate":
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid terminal declaration",
					Detail:  fmt.Sprintf("The keyword %q must appear either as the first keyword of a terminal definition or immediately after \"power\".", extraKw),
					Ranges:  extraKwTok.Range.List(),
				})
			case "input", "output", "bidi":
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid terminal declaration",
					Detail:  fmt.Sprintf("The keyword %q may only be used after the \"output\" keyword.", extraKw),
					Ranges:  extraKwTok.Range.List(),
				})
			default:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Unterminated statement",
					Detail:  "This terminal declaration must be terminated by a semicolon.",
					Ranges:  terminal.Range.List(),
				})
			}
			p.setRecovering()
		}
		p.recoverAfterSemicolon()
		return terminal, diags
	}

	close := p.Read() // eat semicolon
	terminal.WithRange.Range = source.RangeBetween(first.Range, close.Range)

	return terminal, diags
}

func (p *parser) parseNamedObjectBlock() (name string, params *ast.Arguments, body *ast.StatementBlock, headerRange source.Range, fullRange source.Range, diags source.Diags) {
	kwTok := p.Peek()
	if kwTok.Type != TokenIdent {
		// Should never happen because caller should've peeked ahead here
		panic("parseNamedObjectBlock called with peeker not pointing at ident")
	}

	kw := p.PeekKeyword()
	p.Read() // eat keyword

	// Start with reasonable values for all of our results that are valid
	// enough to return on an error, and then we'll fix these up to be more
	// useful as we go along.
	headerRange = kwTok.Range
	fullRange = kwTok.Range
	params = &ast.Arguments{
		WithRange: ast.WithRange{
			Range: source.Range{
				Start: kwTok.Range.End,
				End:   kwTok.Range.End,
			},
		},
	}
	body = &ast.StatementBlock{
		WithRange: ast.WithRange{
			Range: source.Range{
				Start: kwTok.Range.End,
				End:   kwTok.Range.End,
			},
		},
	}

	if p.Peek().Type != TokenIdent {
		if !p.recovering {
			switch p.Peek().Type {
			case TokenOBrace:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: fmt.Sprintf("Missing %s name", kw),
					Detail:  fmt.Sprintf("The %q keyword must be followed by a name for this %s.", kw, kw),
					Ranges:  []source.Range{p.PeekRange()},
				})
			default:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: fmt.Sprintf("Invalid %s name", kw),
					Detail:  fmt.Sprintf("The name of this %s must be an identifier.", kw),
					Ranges:  []source.Range{p.PeekRange()},
				})
			}
		}
		fullRange = source.RangeBetween(fullRange, p.PeekRange())
		p.recoverAfterNextBlock()
		return
	}

	nameTok := p.Read()
	name = p.decodeIdentifierBytes(nameTok.Bytes)

	var paramsDiags source.Diags
	params, paramsDiags = p.parseParameters()
	diags = append(diags, paramsDiags...)
	if len(params.Positional) == 0 {
		headerRange = source.RangeBetween(kwTok.Range, nameTok.Range)
	} else {
		headerRange = source.RangeBetween(kwTok.Range, params.Range)
	}

	var bodyDiags source.Diags
	body, bodyDiags = p.parseStmtBlock()
	diags = append(diags, bodyDiags...)
	fullRange = source.RangeBetween(headerRange, body.SourceRange())

	return

}

func (p *parser) parseBoard() (ast.Node, source.Diags) {
	kw := p.PeekKeyword()
	if kw != "board" {
		// Should never happen because caller should've peeked ahead here
		panic("parseBoard called with peeker not pointing at board keyword")
	}

	name, params, body, headerRange, fullRange, diags := p.parseNamedObjectBlock()

	if len(params.Positional) != 0 && !diags.HasErrors() {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Board may not have parameters",
			Detail:  "A \"board\" block does not accept parameters, so no parameter list is allowed.",
			Ranges:  params.SourceRange().List(),
		})
	}

	return &ast.Board{
		Name: name,
		Body: body,

		HeaderRange: headerRange,
		WithRange: ast.WithRange{
			Range: fullRange,
		},
	}, diags
}

func (p *parser) parseCircuit() (ast.Node, source.Diags) {
	kw := p.PeekKeyword()
	if kw != "circuit" {
		// Should never happen because caller should've peeked ahead here
		panic("parseCircuit called with peeker not pointing at circuit keyword")
	}

	name, params, body, headerRange, fullRange, diags := p.parseNamedObjectBlock()

	return &ast.Circuit{
		Name:   name,
		Params: params,
		Body:   body,

		HeaderRange: headerRange,
		WithRange: ast.WithRange{
			Range: fullRange,
		},
	}, diags
}

func (p *parser) parseDevice() (ast.Node, source.Diags) {
	kw := p.PeekKeyword()
	if kw != "device" {
		// Should never happen because caller should've peeked ahead here
		panic("parseDevice called with peeker not pointing at device keyword")
	}

	name, params, body, headerRange, fullRange, diags := p.parseNamedObjectBlock()

	return &ast.Device{
		Name:   name,
		Params: params,
		Body:   body,

		HeaderRange: headerRange,
		WithRange: ast.WithRange{
			Range: fullRange,
		},
	}, diags
}

func (p *parser) parseLand() (ast.Node, source.Diags) {
	kw := p.PeekKeyword()
	if kw != "land" {
		// Should never happen because caller should've peeked ahead here
		panic("parseLand called with peeker not pointing at land keyword")
	}

	name, params, body, headerRange, fullRange, diags := p.parseNamedObjectBlock()

	return &ast.Land{
		Name:   name,
		Params: params,
		Body:   body,

		HeaderRange: headerRange,
		WithRange: ast.WithRange{
			Range: fullRange,
		},
	}, diags
}

func (p *parser) parsePinout() (ast.Node, source.Diags) {
	kw := p.PeekKeyword()
	if kw != "pinout" {
		// Should never happen because caller should've peeked ahead here
		panic("parsePinout called with peeker not pointing at pinout keyword")
	}

	kwTok := p.Read() // eat keyword

	var diags source.Diags
	pinout := &ast.Pinout{
		Land: &ast.Invalid{
			WithRange: ast.WithRange{
				Range: source.Range{
					Start: kwTok.Range.End,
					End:   kwTok.Range.End,
				},
			},
		},
		Body: &ast.StatementBlock{
			WithRange: ast.WithRange{
				Range: source.Range{
					Start: kwTok.Range.End,
					End:   kwTok.Range.End,
				},
			},
		},

		HeaderRange: kwTok.Range,
		WithRange: ast.WithRange{
			Range: kwTok.Range,
		},
	}

	if p.Peek().Type != TokenIdent {
		if !p.recovering {
			switch p.Peek().Type {
			case TokenOBrace:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: fmt.Sprintf("Missing %s name", kw),
					Detail:  fmt.Sprintf("The %q keyword must be followed by a name for this %s.", kw, kw),
					Ranges:  []source.Range{p.PeekRange()},
				})
			default:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: fmt.Sprintf("Invalid %s name", kw),
					Detail:  fmt.Sprintf("The name of this %s must be an identifier.", kw),
					Ranges:  []source.Range{p.PeekRange()},
				})
			}
		}
		pinout.WithRange.Range = source.RangeBetween(kwTok.Range, p.PeekRange())
		p.recoverAfterNextBlock()
		return pinout, diags
	}

	nameTok := p.Read()
	pinout.Name = p.decodeIdentifierBytes(nameTok.Bytes)
	pinout.HeaderRange = source.RangeBetween(kwTok.Range, nameTok.Range)

	if p.PeekKeyword() == "from" {
		p.Read() // eat "from" keyword
		var exprDiags source.Diags
		pinout.Device, exprDiags = p.parseExpr()
		diags = append(diags, exprDiags...)
		pinout.HeaderRange = source.RangeBetween(kwTok.Range, pinout.Device.SourceRange())
		pinout.WithRange.Range = pinout.HeaderRange
		if diags.HasErrors() {
			p.recoverAfterNextBlock()
			return pinout, diags
		}
	}

	if p.PeekKeyword() == "to" {
		p.Read() // eat "to" keyword
		var exprDiags source.Diags
		pinout.Land, exprDiags = p.parseExpr()
		diags = append(diags, exprDiags...)
		pinout.HeaderRange = source.RangeBetween(kwTok.Range, pinout.Land.SourceRange())
		pinout.WithRange.Range = pinout.HeaderRange
		if diags.HasErrors() {
			p.recoverAfterNextBlock()
			return pinout, diags
		}
	} else {
		if !p.recovering {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Missing \"to\" clause",
				Detail:  "Pinout definition must include \"to\" keyword followed by a land expression.",
				Ranges:  pinout.HeaderRange.List(),
			})
		}
		p.setRecovering()
	}

	var bodyDiags source.Diags
	pinout.Body, bodyDiags = p.parseStmtBlock()
	diags = append(diags, bodyDiags...)
	pinout.WithRange.Range = source.RangeBetween(pinout.HeaderRange, pinout.Body.SourceRange())

	return pinout, diags
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

Trailers:
	for {
		next := p.Peek()
		switch next.Type {

		case TokenDot:
			dot := p.Read()
			if p.Peek().Type != TokenIdent {
				if !p.recovering {
					diags = append(diags, source.Diag{
						Level:   source.Error,
						Summary: "Invalid attribute name",
						Detail:  "Expected the name of an attribute to access.",
						Ranges:  p.Peek().Range.List(),
					})
				}
				p.setRecovering()
				// Still mark the place where an attribute is being accessed
				// for use in analysis for e.g. autocomplete.
				term = &ast.GetAttr{
					Name:   "",
					Source: term,

					WithRange: ast.WithRange{
						Range: source.RangeBetween(term.SourceRange(), dot.Range),
					},
				}
				return term, diags
			}

			ident := p.Read()
			name := p.decodeIdentifierBytes(ident.Bytes)

			term = &ast.GetAttr{
				Name:   name,
				Source: term,

				WithRange: ast.WithRange{
					Range: source.RangeBetween(term.SourceRange(), ident.Range),
				},
			}

		case TokenOBrack:
			p.Read() // eat open bracket
			idx, idxDiags := p.parseExpr()
			diags = append(diags, idxDiags...)

			if idxDiags.HasErrors() {
				p.recoverAfterClose(TokenCBrack)
			}

			if p.Peek().Type != TokenCBrack {
				if !p.recovering {
					diags = append(diags, source.Diag{
						Level:   source.Error,
						Summary: "Mismatched brackets",
						Detail:  "Expected a closing bracket \"]\" to mark the end of the index expression.",
						Ranges:  p.Peek().Range.List(),
					})
				}
				p.recoverAfterClose(TokenCBrack)
			}

			close := p.Read()

			term = &ast.GetIndex{
				Source: term,
				Index:  idx,

				WithRange: ast.WithRange{
					Range: source.RangeBetween(term.SourceRange(), close.Range),
				},
			}

		case TokenOParen:
			args, argDiags := p.parseArguments()
			diags = append(diags, argDiags...)

			// parseArguments tries itself to recover to after the closing
			// paren on errors, so we'll just continue and assume that the
			// peeker is already placed as best it can be.

			term = &ast.Call{
				Callee: term,
				Args:   args,

				WithRange: ast.WithRange{
					Range: source.RangeBetween(term.SourceRange(), args.SourceRange()),
				},
			}

		default:
			break Trailers

		}
	}

	return term, diags
}

func (p *parser) parseExpressionTerm() (ast.Node, source.Diags) {
	start := p.Peek()

	switch start.Type {
	case TokenOParen:
		open := p.Read()

		expr, diags := p.parseExpr()
		close := p.Peek()
		if close.Type != TokenCParen && !p.recovering {
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
		} else {
			p.Read() // eat closing paren
		}

		// We return a ParenExpr even in the case of errors, so that the
		// AST is complete as possible for syntax-only analyses such as
		// autocomplete.
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

	case TokenNumberLit:
		tok := p.Read()
		val, diags := p.decodeNumberLiteral(tok)

		next := p.Peek()
		switch next.Type {
		case TokenPercent:
			marker := p.Read()
			if val != nil {
				val.Quo(val, oneHundred)
			}
			return &ast.NumberLit{
				WithRange: ast.WithRange{
					Range: source.RangeBetween(tok.Range, marker.Range),
				},
				Value: val,
			}, diags
		case TokenIdent:
			kw := p.PeekKeyword()
			if ast.IsQuantityUnitKeyword(kw) {
				marker := p.Read()
				return &ast.NumberLit{
					WithRange: ast.WithRange{
						Range: source.RangeBetween(tok.Range, marker.Range),
					},
					Value: val,
					Unit:  kw,
				}, diags
			}
		}

		return &ast.NumberLit{
			WithRange: ast.WithRange{
				Range: tok.Range,
			},
			Value: val,
		}, diags

	case TokenBang:
		op := p.Read()
		// Important to use parseExpressionWithTrailers rather than
		// parseExpression here, or else we can capture a following binary
		// expression into our negation.
		operand, diags := p.parseExpressionWithTrailers()
		return &ast.ArithmeticUnary{
			Op:      ast.Not,
			Operand: operand,

			WithRange: ast.WithRange{
				Range: source.RangeBetween(op.Range, operand.SourceRange()),
			},
		}, diags

	case TokenMinus:
		op := p.Read()
		// Important to use parseExpressionWithTrailers rather than
		// parseExpression here, or else we can capture a following binary
		// expression into our negation.
		operand, diags := p.parseExpressionWithTrailers()
		return &ast.ArithmeticUnary{
			Op:      ast.Negate,
			Operand: operand,

			WithRange: ast.WithRange{
				Range: source.RangeBetween(op.Range, operand.SourceRange()),
			},
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

func (p *parser) parseParameters() (*ast.Arguments, source.Diags) {
	// parseParameters raturns an ast.Arguments that meets the constraints for
	// a parameter list: contains only positional arguments, and all of the
	// arguments are just direct variable references.
	//
	// If the result parses as an argument list but does _not_ meet the constraints
	// for parameter lists, the argument list is returned as parsed but error
	// diagnostics are emitted.

	if p.Peek().Type != TokenOParen {
		// Parameter lists are optional, so if there's no open then we'll
		// assume that the parameter list has been omitted and return
		// an empty one.
		rng := p.PeekRange()
		rng.End = rng.Start // make empty range just before the next token
		return &ast.Arguments{
			WithRange: ast.WithRange{
				Range: rng,
			},
		}, nil
	}

	args, diags := p.parseArguments()
	if diags.HasErrors() {
		return args, diags
	}

	for _, n := range args.Positional {
		if _, isVar := n.(*ast.Variable); !isVar {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Invalid parameter declaration",
				Detail:  "A parameter declaration must be just the parameter's name.",
				Ranges:  n.SourceRange().List(),
			})
		}
	}

	// Verify that this arguments object is a valid parameter declaration.
	if len(args.Named) != 0 {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Invalid parameter declaration",
			Detail:  "Default values may not be declared within the parameter list. Use 'attr' statements within the following block to define types and default values.",
			Ranges:  args.Named[0].SourceRange().List(),
		})
	}

	return args, diags
}

func (p *parser) parseArguments() (*ast.Arguments, source.Diags) {
	open := p.Read()
	if open.Type != TokenOParen {
		// indicates a bug in the caller
		panic("parseArguments called with peeker not pointing at TokenOParen")
	}

	var diags source.Diags
	ret := &ast.Arguments{}
	first := true

Arguments:
	for {
		if p.Peek().Type == TokenCParen {
			close := p.Read()
			ret.WithRange.Range = source.RangeBetween(open.Range, close.Range)
			break Arguments
		}

		if !first {
			if p.Peek().Type != TokenComma {
				if !p.recovering {
					diags = append(diags, source.Diag{
						Level:   source.Error,
						Summary: "Missing argument separator",
						Detail:  "Call arguments must be separated by commas.",
						Ranges:  p.Peek().Range.List(),
					})
				}
				ret.WithRange.Range = source.RangeBetween(open.Range, p.Peek().Range)
				p.recoverAfterClose(TokenCParen)
				break Arguments
			}

			p.Read() // eat comma

			if p.Peek().Type == TokenCParen {
				close := p.Read()
				ret.WithRange.Range = source.RangeBetween(open.Range, close.Range)
				break Arguments
			}
		}
		first = false

		var nameExpr ast.Node
		argExpr, argDiags := p.parseExpr()
		diags = append(diags, argDiags...)
		if argDiags.HasErrors() {
			ret.WithRange.Range = source.RangeBetween(open.Range, argExpr.SourceRange())
			p.recoverAfterClose(TokenCParen)
			break Arguments
		}

		if p.Peek().Type == TokenAssign {
			// A named argument

			p.Read() // eat equals sign

			nameExpr = argExpr
			argExpr, argDiags = p.parseExpr()
			diags = append(diags, argDiags...)
			if argDiags.HasErrors() {
				ret.WithRange.Range = source.RangeBetween(open.Range, argExpr.SourceRange())
				p.recoverAfterClose(TokenCParen)
				break Arguments
			}
		}

		if nameExpr == nil {
			// positional argument

			if len(ret.Named) != 0 {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Incorrect argument order",
					Detail:  "Positional arguments must all be listed before the first named argument.",
					Ranges:  argExpr.SourceRange().List(),
				})
			}

			ret.Positional = append(ret.Positional, argExpr)
		} else {
			var name string
			if varNode, isVar := nameExpr.(*ast.Variable); isVar {
				name = varNode.Name
			} else {
				if !p.recovering {
					diags = append(diags, source.Diag{
						Level:   source.Error,
						Summary: "Invalid parameter name",
						Detail:  "A parameter name must be a valid identifier.",
						Ranges:  argExpr.SourceRange().List(),
					})
				}
			}

			ret.Named = append(ret.Named, &ast.NamedArgument{
				Name:  name,
				Value: argExpr,

				WithRange: ast.WithRange{
					Range: source.RangeBetween(nameExpr.SourceRange(), argExpr.SourceRange()),
				},
			})
		}
	}

	return ret, diags
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

func (p *parser) decodeNumberLiteral(tok Token) (*big.Float, source.Diags) {
	if tok.Type != TokenNumberLit {
		panic("decodeNumberLiteral can only be used with TokenNumberLit tokens")
	}

	var diags source.Diags
	str := string(tok.Bytes)
	f := &big.Float{}
	_, _, err := f.Parse(str, 10)
	if err != nil {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Invalid number literal",
			Detail:  "The given number is invalid.",
			Ranges:  tok.Range.List(),
		})
	}

	return f, diags
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
			nest--
			if nest < 1 {
				return tok
			}
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

func mustParseBigFloat(str string) *big.Float {
	f, _, err := (&big.Float{}).Parse(str, 10)
	if err != nil {
		panic(err)
	}
	return f
}

// pathHasExtension checks if the given path has the given extension (suffix)
// while also ignoring files that have names starting with "." or "_" that
// are presumed to be temporary files created by editors or other tools.
func pathHasExtension(path, ext string) bool {
	if !strings.HasSuffix(path, ext) {
		return false
	}

	if strings.HasPrefix(path, ".") || strings.HasPrefix(path, "_") {
		return false
	}

	return true
}
