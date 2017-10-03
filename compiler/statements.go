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
	case *ast.Import:
		sym := scope.Get(tn.SymbolName())
		return eval.ImportStmt(tn.Package, sym, tn.SourceRange()), nil
	case *ast.Attr:
		sym := scope.Get(tn.Name)

		// attr statements are special in that they always compile in the
		// parent scope. This prevents an attribute's default value or type
		// from depending on something within the body, which would prevent
		// us from successfully identifying the required types before
		// execution.
		scope = scope.Parent()
		if scope == nil {
			// should never happen
			panic("attempt to compile attr statement in the global scope")
		}

		switch {
		case tn.Value != nil:
			defVal, diags := CompileExpr(tn.Value, scope)
			return eval.AttrStmtDefault(sym, defVal, tn.SourceRange()), diags
		case tn.Type != nil:
			typeExpr, diags := CompileExpr(tn.Type, scope)
			return eval.AttrStmt(sym, typeExpr, tn.SourceRange()), diags
		default:
			// should never happen
			panic("invalid *ast.Attr: neither Value nor Type is set")
		}
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
