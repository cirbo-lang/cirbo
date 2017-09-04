package parser

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

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

	p.files[fpath] = ret

	var diags source.Diags

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

	tokens := scanTokens(src, fpath, source.Pos{Line: 1, Column: 1}, scanNormal)
	return tokens, nil
}

func (p *Parser) Diagnostics() source.Diags {
	return p.diags
}
