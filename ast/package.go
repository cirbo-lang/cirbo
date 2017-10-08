package ast

// Package is the top-level object in the AST, representing the set of files
// associated with a particular package.
//
// A Package is not actually a Node.
type Package []*File

func (p Package) Files() []*File {
	return []*File(p)
}

// VisitAll is a helper that calls the top-level VisitAll for each of the
// files in turn.
func (p Package) VisitAll(cb VisitFunc) {
	for _, f := range p {
		VisitAll(f, cb)
	}
}

// Walk is a helper that calls the top-level Walk for each of the files
// in turn.
func (p Package) Walk(w Walker) {
	for _, f := range p {
		Walk(f, w)
	}
}
