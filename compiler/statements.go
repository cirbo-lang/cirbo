package compiler

import (
	"fmt"
	"strconv"

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

	// With all of the explicit definitions dealt with, we also need to go
	// hunting for implicit definitions where the user calls a non-existant
	// function whose name looks like a component reference for a simple,
	// intrinsic device like a resistor.
	impliedStmts, swap, promDiags := promoteIntrinsicDecls(nodes, scope)
	diags = append(diags, promDiags...)
	if len(impliedStmts) > 0 {
		for _, node := range impliedStmts {
			decl := declForNode(node)
			if decl.Name == "" {
				continue
			}
			scope.Declare(decl.Name)
			declRange[decl.Name] = decl.Range
		}

		newNodes := make([]ast.Node, 0, len(nodes)+len(impliedStmts))
		newNodes = append(newNodes, nodes...)
		newNodes = append(newNodes, impliedStmts...)
		nodes = newNodes
	}

	var stmts []eval.Stmt
	for _, node := range nodes {
		stmt, stmtDiags := compileStatement(node, scope, swap)
		diags = append(diags, stmtDiags...)
		stmts = append(stmts, stmt)
	}

	block, blockDiags := eval.MakeStmtBlock(scope, stmts)
	diags = append(diags, blockDiags...)

	return block, diags
}

func compileStatement(node ast.Node, scope *eval.Scope, swap ast.SwapTable) (eval.Stmt, source.Diags) {
	switch tn := node.(type) {
	case *ast.Assign:
		expr, diags := compileExpr(tn.Value, scope, swap)
		sym := scope.Get(tn.Name)
		return eval.AssignStmt(sym, expr, tn.SourceRange()), diags
	case *ast.Import:
		sym := scope.Get(tn.SymbolName())
		return eval.ImportStmt(tn.Package, sym, tn.SourceRange()), nil
	case *ast.Export:
		expr, diags := compileExpr(tn.Value, scope, swap)
		return eval.ExportStmt(expr, tn.SourceRange()), diags
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
			// Not using the swap table here because we're compiling in the
			// parent scope, and so the swap won't be relevant here.
			defVal, diags := CompileExpr(tn.Value, scope)
			return eval.AttrStmtDefault(sym, defVal, tn.SourceRange()), diags
		case tn.Type != nil:
			// Not using the swap table here because we're compiling in the
			// parent scope, and so the swap won't be relevant here.
			typeExpr, diags := CompileExpr(tn.Type, scope)
			return eval.AttrStmt(sym, typeExpr, tn.SourceRange()), diags
		default:
			// should never happen
			panic("invalid *ast.Attr: neither Value nor Type is set")
		}
	case *ast.Designator:
		expr, diags := compileExpr(tn.Value, scope, swap)
		return eval.DesignatorStmt(expr, tn.SourceRange()), diags
	case *ast.Device:
		sym := scope.Get(tn.Name)
		block, diags := compileStatements(tn.Body.Statements, scope)
		params, paramDiags := compilePositionalParams(tn.Params.Positional, block.AttributeNames())
		diags = append(diags, paramDiags...)
		return eval.DeviceStmt(sym, params, block, tn.SourceRange()), diags
	default:
		panic(fmt.Errorf("%T cannot be compiled to a statement", node))
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

type promotionWalker struct {
	scope *eval.Scope
	swap  ast.SwapTable
	diags source.Diags
	stmts []ast.Node
	ours  map[*eval.Symbol]source.Range
}

func promoteIntrinsicDecls(nodes []ast.Node, scope *eval.Scope) ([]ast.Node, ast.SwapTable, source.Diags) {
	walker := &promotionWalker{
		scope: scope,
	}

	for _, node := range nodes {
		ast.Walk(node, walker)
	}

	return walker.stmts, walker.swap, walker.diags
}

func (w *promotionWalker) EnterNode(node ast.Node) bool {
	callExpr, isCall := node.(*ast.Call)
	if !isCall {
		return true
	}

	callee := callExpr.Callee
	varExpr, isVar := callee.(*ast.Variable)
	if !isVar {
		return true
	}

	name := varExpr.Name
	if sym := w.scope.Get(name); sym != nil {
		// If the name is already defined, then it can't be an implicit
		// declaration.

		// If it's a symbol we defined though, that's an error because the
		// user is attempting to implicitly define the same name in two
		// places.
		if rng, ours := w.ours[sym]; ours {
			w.diags = append(w.diags, source.Diag{
				Level:   source.Error,
				Summary: "Duplicate implicit declaration",
				Detail: fmt.Sprintf(
					"The name %q was already implicitly defined at %s. If terminals of this device must be connected multiple times, assign it to a name explicitly and connect using that name.",
					name, rng,
				),
				Ranges: varExpr.Range.List(),
			})
		}

		return true
	}

	if len(name) < 2 {
		// If it's not at least two characters long then it can't possibly
		// be a component-reference-style name.
		return true
	}

	var prefix, num string

	// The full set of reference designators we accept as implicit decls here
	// are the following:
	//     R    Resistor
	//     C    Capacitor (unpolarized)
	//     CP   Capacitor (polarized)
	//     D    Diode
	//     L    Inductor
	//     Z    Zener Diode
	//     F    Fuse
	//     FB   Ferrite Bead
	//
	// This implicit definition mechanism only makes sense for components that
	// have just "IN" and "OUT" terminals, since the primary purpose of it
	// is to conveniently instantiate simple components within connection
	// statements:
	//      dev.FOO -- R1(1kohm) -- GND; // pull-down resistor

	switch name[0] {
	case 'R', 'C', 'D', 'L', 'Z', 'F':
		// all possible implied definition prefixes
		prefix, num = name[0:1], name[1:]
	default:
		return true
	}

	// Deal with a few exceptions where the designator is longer
	switch {
	// Polarized capacitors (CP) and ferrite beads (FB)
	case prefix == "C" && num[0] == 'P' || prefix == "F" && num[0] == 'B':
		prefix, num = name[0:2], name[2:]
		if len(num) == 0 {
			// can't be a CPn name, then
			return true
		}
	}

	// The number portion must be all digits to be valid. We'll check that
	// by trying to parse it as an integer.
	if _, err := strconv.Atoi(num); err != nil {
		return true
	}

	// If we got though all of the above (phew!), then we have an implicit
	// definition to record.
	if w.swap == nil {
		w.swap = ast.SwapTable{}
	}

	// Replace the original call with a reference to the new name
	newRef := &ast.Variable{
		Name: name,
		WithRange: ast.WithRange{
			Range: callExpr.Range,
		},
	}
	w.swap.Add(node, newRef)

	// Construct a synthetic assignment statement that calls to the global
	// device constructor named after our prefix.
	newCall := &ast.Call{
		Callee: &ast.Variable{
			// If the name has been redefined in a child scope then we'll use
			// the overridden version. This isn't really the ideal, but we'll
			// accept it for implementation simplicity and simply suggest that
			// redefining the intrinsic device constructors is a bad idea.
			Name: prefix,
			WithRange: ast.WithRange{
				Range: source.Range{
					Filename: varExpr.Range.Filename,
					Start:    varExpr.Range.Start,
					End: source.Pos{
						Line:   varExpr.Range.Start.Line,
						Column: varExpr.Range.Start.Column + len(prefix),
						Byte:   varExpr.Range.Start.Byte + len(prefix),
					},
				},
			},
		},
		Args:      callExpr.Args,
		WithRange: callExpr.WithRange,
	}
	w.stmts = append(w.stmts, &ast.Assign{
		Name:  name,
		Value: newCall,
		WithRange: ast.WithRange{
			Range: callExpr.Range,
		},
	})

	// Still recurse into the child nodes of the call, since there might be
	// more implicit definitions lurking inside.
	return true
}

func (w *promotionWalker) ExitNode(node ast.Node) {
	// nothing to do
}
