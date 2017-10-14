package compiler

import (
	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/source"
)

// CompilePackage compiles the given package and returns the value exported
// by it.
//
// If the package has an explicit "export" statement then its argument is
// returned. Otherwise, symbols from the module's symbol table are used to
// construct and instantiate an object type.
func CompilePackage(pkg ast.Package) (*eval.Package, source.Diags) {
	var diags source.Diags

	// First we'll make sure the top-level statements in all of our files
	// are permitted at the module level.
	for _, file := range pkg {
		for _, node := range file.TopLevel {
			switch tn := node.(type) {
			case *ast.Assign, *ast.Import, *ast.Export, *ast.Circuit, *ast.Device, *ast.Land, *ast.Pinout:
				// allowed
			case *ast.Connection:
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid top-level statement",
					Detail:  "Connection statements may not appear in the top-level scope of a module.",
					Ranges:  tn.SourceRange().List(),
				})
			default:
				// Generic error message for all other types
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Invalid top-level statement",
					Detail:  "This statement is not allowed in the top-level scope of a module.",
					Ranges:  tn.SourceRange().List(),
				})
			}
		}
	}

	var nodes []ast.Node
	for _, file := range pkg {
		nodes = append(nodes, file.TopLevel...)
	}
	block, compileDiags := compileStatements(nodes, eval.GlobalScope())
	diags = append(diags, compileDiags...)

	return eval.NewPackage(block), diags
}
