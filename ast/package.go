package ast

// Package is the top-level object in the AST, representing the set of files
// associated with a particular package.
//
// A Package is not actually a Node.
type Package struct {
	DefaultName string
	Files       []*File
}

// VisitAll is a helper that calls the top-level VisitAll for each of the
// files in turn.
func (p *Package) VisitAll(cb VisitFunc) {
	for _, f := range p.Files {
		VisitAll(f, cb)
	}
}

// Walk is a helper that calls the top-level Walk for each of the files
// in turn.
func (p *Package) Walk(w Walker) {
	for _, f := range p.Files {
		Walk(f, w)
	}
}
