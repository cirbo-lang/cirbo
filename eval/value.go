package eval

import (
	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/source"
)

// Value evaluates the given expression in the given context and returns the
// result, along with any diagnostics that are produced during the operation.
func Value(expr Expr, ctx *Context) (cty.Value, source.Diags) {
	return expr.value(ctx, nil)
}

// DefineSymbol evaluates the given expression in the given context and
// then defines the given symbol (in the same context) with the result.
//
// As a convenience, the result of the expression is returned along with
// any diagnostics produced during the evaluation.
//
// It is invalid to re-define a symbol, so this function will panic if the
// given symbol already has a definition in the given context.
func DefineSymbol(sym *Symbol, expr Expr, ctx *Context) (cty.Value, source.Diags) {
	v, diags := expr.value(ctx, sym)
	ctx.Define(sym, v)
	return v, diags
}
