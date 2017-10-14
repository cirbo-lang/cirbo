package compiler

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/source"
)

func compilePositionalParams(decls []ast.Node, attrNames map[string]*eval.Symbol) (eval.PosParameters, source.Diags) {
	// We expect "decls" here to be a list of ast.Variable nodes, as would be
	// produced by the parser for a parameter list. However, since we want to
	// allow partial analysis of invalid source, we'll tolerate and ignore
	// invalid stuff under the assumption that diagnostics were already
	// reported during parsing, since the parser _will_ pass through other
	// AST nodes in a parameter list if they are present.

	// attrNames may be nil to just extract the given names, without validating
	// them at all. If passed, however, diagnostics will be generated for each
	// parameter that doesn't have a corresponding attribute.

	var params eval.PosParameters
	var diags source.Diags

	for _, param := range decls {
		vn, isVar := param.(*ast.Variable)
		if !isVar {
			// Not valid in a parameter list, but we'll ignore for the
			// reasons noted above.
			continue
		}

		if attrNames != nil {
			_, declared := attrNames[vn.Name]
			if !declared {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Parameter without corresponding attribute",
					Detail:  fmt.Sprintf("The parameter name %q does not have a corresponding attribute in the block body. Define such an attribute in order to specify this parameter's required type.", vn.Name),
					Ranges:  param.SourceRange().List(),
				})
				continue
			}
		}

		params = append(params, eval.Parameter{
			Name:        vn.Name,
			SourceRange: vn.SourceRange(),
		})
	}

	return params, diags
}
