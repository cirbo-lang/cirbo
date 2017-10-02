package compiler

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/source"
)

func compileStatements(nodes []ast.Node, parentScope *eval.Scope) (eval.StmtBlock, source.Diags) {
	var diags source.Diags
	declRange := map[string]source.Range{}
	scope := parentScope.NewChild()

	for _, node := range nodes {
		decl := declForNode(node)
		if decl.Name == "" {
			continue
		}

		if rng, exists := declRange[decl.Name]; exists {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Duplicate declaration",
				Detail:  fmt.Sprintf("The name %q was already used in the declaration at %s.", decl.Name, rng),
				Ranges:  decl.Range.List(),
			})
			continue
		}

		scope.Declare(decl.Name)
		declRange[decl.Name] = decl.Range
	}

	// TODO: Deal with any implicit device instantiations that may appear
	// in connection statements. These need to be added as additional
	// symbols with implied assign statements and then have some special
	// treatment applied so that when we compile the expr we'll refer to
	// these additional symbols.

	var stmts []eval.Stmt
	for _, node := range nodes {
		stmt, stmtDiags := compileStatement(node, scope)
		diags = append(diags, stmtDiags...)
		stmts = append(stmts, stmt)
	}

	block, blockDiags := eval.MakeStmtBlock(scope, stmts)
	diags = append(diags, blockDiags...)

	return block, diags
}

func compileStatement(node ast.Node, scope *eval.Scope) (eval.Stmt, source.Diags) {
	switch tn := node.(type) {
	case *ast.Assign:
		expr, diags := CompileExpr(tn.Value, scope)
		sym := scope.Get(tn.Name)
		return eval.AssignStmt(sym, expr, tn.SourceRange()), diags
	default:
		panic(fmt.Errorf("%#v cannot be compiled to a statement", node))
	}
}

func declForNode(node ast.Node) symbolDecl {
	switch tn := node.(type) {
	case *ast.Assign:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.SourceRange(),
		}
	case *ast.Import:
		return symbolDecl{
			Name:  tn.SymbolName(),
			Range: tn.SourceRange(),
		}
	case *ast.Attr:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.SourceRange(),
		}
	case *ast.Terminal:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.SourceRange(),
		}
	case *ast.Circuit:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.DeclRange(),
		}
	case *ast.Device:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.DeclRange(),
		}
	case *ast.Land:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.DeclRange(),
		}
	case *ast.Board:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.DeclRange(),
		}
	case *ast.Pinout:
		return symbolDecl{
			Name:  tn.Name,
			Range: tn.DeclRange(),
		}
	default:
		return symbolDecl{}
	}
}

type symbolDecl struct {
	Name  string
	Range source.Range
}
