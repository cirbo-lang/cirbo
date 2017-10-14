package eval

import (
	"github.com/cirbo-lang/cirbo/source"
)

// Parameter represents a single parameter declared for something that can be
// called.
type Parameter struct {
	Name        string
	SourceRange source.Range
}

// PosParameters represents a sequence of positional parameters.
type PosParameters []Parameter

// NamedParameters represents a set of named parameters.
type NamedParameters map[string]Parameter

// Parameters represents a collection of positional and named parameters for
// something that can be called.
type Parameters struct {
	Positional PosParameters
	Named      NamedParameters
}
