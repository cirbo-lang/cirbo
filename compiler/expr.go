package compiler

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/units"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/source"
)

// CompileExpr transforms the given node (and its child nodes, if any) into
// an evaluatable expression for the given scope.
//
// Only the subset of AST nodes that represent expression components can be
// passed here. If an unsupported node is passed, this function will panic.
// Callers should generally only pass AST subtrees that are documented as
// representing expressions.
func CompileExpr(node ast.Node, scope *eval.Scope) (eval.Expr, source.Diags) {
	switch tn := node.(type) {
	case *ast.ParenExpr:
		cn, diags := CompileExpr(tn.Content, scope)
		return eval.PassthroughExpr(cn, tn.SourceRange()), diags
	case *ast.BooleanLit:
		tv := cty.BoolVal(tn.Value)
		return eval.LiteralExpr(tv, tn.SourceRange()), nil
	case *ast.StringLit:
		tv := cty.StringVal(tn.Value)
		return eval.LiteralExpr(tv, tn.SourceRange()), nil
	case *ast.NumberLit:
		unitName := tn.Unit
		if unitName == "" {
			tv := cty.QuantityVal(units.MakeDimensionless(tn.Value))
			return eval.LiteralExpr(tv, tn.SourceRange()), nil
		}
		unit := units.ByName(tn.Unit)
		if unit == nil {
			suggestion := nameSuggestion(tn.Unit, units.AllNames())
			if suggestion != "" {
				suggestion = fmt.Sprintf(" Did you mean %q?", suggestion)
			}
			return placeholderExpr(tn.SourceRange()), source.Diags{
				{
					Level:   source.Error,
					Summary: "Invalid unit",
					Detail:  fmt.Sprintf("The name %q is not a known quantity unit.%s", tn.Unit, suggestion),
					Ranges:  tn.SourceRange().List(),
				},
			}
		}
		tv := cty.QuantityVal(units.MakeQuantity(tn.Value, unit))
		return eval.LiteralExpr(tv, tn.SourceRange()), nil
	case *ast.Variable:
		sym := scope.Get(tn.Name)
		if sym == nil {
			suggestion := nameSuggestion(tn.Name, scope.AllNames())
			if suggestion != "" {
				suggestion = fmt.Sprintf(" Did you mean %q?", suggestion)
			}
			return placeholderExpr(tn.SourceRange()), source.Diags{
				{
					Level:   source.Error,
					Summary: "Undeclared variable",
					Detail:  fmt.Sprintf("The name %q has not been declared as a variable in this scope.%s", tn.Name, suggestion),
					Ranges:  tn.SourceRange().List(),
				},
			}
		}
		return eval.SymbolExpr(sym, tn.SourceRange()), nil
	case *ast.ArithmeticBinary:
		var diags source.Diags
		lhs, lhsDiags := CompileExpr(tn.LHS, scope)
		rhs, rhsDiags := CompileExpr(tn.RHS, scope)
		diags = append(diags, lhsDiags...)
		diags = append(diags, rhsDiags...)

		switch tn.Op {
		case ast.Equal:
			return eval.EqualExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.NotEqual:
			return eval.NotEqualExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.LessThan:
			return eval.LessThanExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.LessThanOrEqual:
			return eval.LessThanOrEqualExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.GreaterThan:
			return eval.GreaterThanExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.GreaterThanOrEqual:
			return eval.GreaterThanOrEqualExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.Add:
			return eval.AddExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.Subtract:
			return eval.SubtractExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.Multiply:
			return eval.MultiplyExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.Divide:
			return eval.DivideExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.Modulo:
			return eval.ModuloExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.Exponent:
			return eval.ExponentExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.And:
			return eval.AndExpr(lhs, rhs, tn.SourceRange()), diags
		case ast.Or:
			return eval.OrExpr(lhs, rhs, tn.SourceRange()), diags
		default:
			panic(fmt.Errorf("compilation of binary %s is not implemented", tn.Op))
		}
	case *ast.ArithmeticUnary:
		operand, diags := CompileExpr(tn.Operand, scope)

		switch tn.Op {
		case ast.Negate:
			return eval.NegateExpr(operand, tn.SourceRange()), diags
		case ast.Not:
			return eval.NotExpr(operand, tn.SourceRange()), diags
		default:
			panic(fmt.Errorf("compilation of unary %s is not implemented", tn.Op))
		}
	case *ast.GetAttr:
		obj, diags := CompileExpr(tn.Source, scope)
		return eval.AttrExpr(obj, tn.Name, tn.SourceRange()), diags
	case *ast.GetIndex:
		var diags source.Diags
		coll, collDiags := CompileExpr(tn.Source, scope)
		diags = append(diags, collDiags...)
		index, indexDiags := CompileExpr(tn.Index, scope)
		diags = append(diags, indexDiags...)
		return eval.IndexExpr(coll, index, tn.SourceRange()), diags
	case *ast.Call:
		var diags source.Diags

		callee, calleeDiags := CompileExpr(tn.Callee, scope)
		diags = append(diags, calleeDiags...)

		posArgs := make([]eval.Expr, len(tn.Args.Positional))
		for i, cn := range tn.Args.Positional {
			var argDiags source.Diags
			posArgs[i], argDiags = CompileExpr(cn, scope)
			diags = append(diags, argDiags...)
		}

		namedArgs := make(map[string]eval.Expr, len(tn.Args.Named))
		for _, arg := range tn.Args.Named {
			expr, argDiags := CompileExpr(arg.Value, scope)
			diags = append(diags, argDiags...)
			if _, already := namedArgs[arg.Name]; already {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Duplicate named argument",
					Detail:  fmt.Sprintf("An argument named %q was already passed.", arg.Name),
					Ranges:  arg.Range.List(),
				})
				continue
			}
			namedArgs[arg.Name] = expr
		}

		return eval.CallExpr(callee, posArgs, namedArgs, tn.SourceRange()), diags
	case *ast.Invalid:
		// In normal situations ast.Invalid should not appear in the AST,
		// but if the parser encountered an error and we tried to compile
		// the resulting AST anyway (e.g. to try to implement autocomplete)
		// then we may encounter this, in which case we'll just treat it
		// as a placeholder literal to allow partial evaluation.
		return placeholderExpr(tn.SourceRange()), nil
	default:
		panic(fmt.Errorf("%#v cannot be compiled to an expression", node))
	}
}

// placeholderExpr returns a LiteralExpr representing cty.PlaceholderVal at
// the given source location. This is used as a return value from CompileExpr
// in erroneous situations where no reasonable real value can be built.
func placeholderExpr(rng source.Range) eval.Expr {
	return eval.LiteralExpr(cty.PlaceholderVal, rng)
}
