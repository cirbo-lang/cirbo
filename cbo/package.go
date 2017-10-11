package cbo

import (
	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/source"
)

// Package is the result of compiling a package directory containing one or more
// Cirbo module files (.cbm files).
type Package struct {
	block eval.StmtBlock
}

func NewPackage(Block eval.StmtBlock) *Package {
	return &Package{
		block: Block,
	}
}

// PackagesImported returns a slice of the package paths imported by the
// receiving package.
//
// These are the packages that must be present in the otherPackages map
// when calling the ExportedValue method.
func (p *Package) PackagesImported() []eval.PackageRef {
	return p.block.PackagesImported()
}

// ExportedValue evaluates the package's declaration statements and returns
// the value exported by the package.
//
// The otherPackages map must include the export value of each of the package
// paths returned by the PackagesImported method, or else the result is
// undefined. It's the caller's responsibility to build the package dependency
// graph and resolve modules in an appropriate order to satisfy each package's
// dependencies.
func (p *Package) ExportedValue(otherPackages map[string]cbty.Value) (cbty.Value, source.Diags) {
	result, diags := p.block.Execute(eval.StmtBlockExecute{
		Context:  eval.GlobalContext(),
		Packages: otherPackages,
	})

	if result.ExportValue != cbty.NilValue {
		return result.ExportValue, diags
	}

	// If the package didn't explicitly export a value, we'll create a synthetic
	// one using the symbols it defined in its scope.
	impliedAttrs := map[string]cbty.Value{}
	for sym := range p.block.ImplicitExports() {
		impliedAttrs[sym.DeclaredName()] = result.Context.Value(sym)
	}
	return cbty.ObjectVal(impliedAttrs), diags
}
