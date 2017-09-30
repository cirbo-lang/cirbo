package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/source"
)

type Expr interface {
	// value evaluates the expression in the given context.
	//
	// If the result is being used directly as the definition of a symbol
	// then the that symbol is provided in "targetSym"; otherwise,
	// targetSym is nil. Some expression types are valid only when being
	// assigned directly to a symbol.
	value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags)

	// eachChild should pass each of its child expressions (e.g. operands)
	// to the given callback in some reasonable order.
	eachChild(cb walkCb)

	// sourceRange returns the source code range that represents the receiving
	// expression.
	sourceRange() source.Range
}

type symbolExpr struct {
	sym *Symbol
	rng
	leafExpr
}

func SymbolExpr(sym *Symbol, rng source.Range) Expr {
	return &symbolExpr{
		sym: sym,
		rng: srcRange(rng),
	}
}

func (e *symbolExpr) value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags) {
	val := ctx.Value(e.sym)
	if val == cty.NilValue {
		// This is actually an implementation error in Cirbo rather than a
		// user error, but we'll return it as a diagnostic anyway since that's
		// more graceful.
		return cty.PlaceholderVal, source.Diags{
			{
				Level:   source.Error,
				Summary: "Symbol not yet defined",
				Detail:  fmt.Sprintf("The symbol %q has not yet been defined. This is a bug in Cirbo that should be reported!", e.sym.name),
				Ranges:  e.sourceRange().List(),
			},
		}
	}
	return val, nil
}

func (e *symbolExpr) GoString() string {
	return fmt.Sprintf("eval.SymbolExpr(%#v, %#v)", e.sym, e.rng.sourceRange())
}

type literalExpr struct {
	val cty.Value
	rng
	leafExpr
}

func LiteralExpr(val cty.Value, rng source.Range) Expr {
	return &literalExpr{
		val: val,
		rng: srcRange(rng),
	}
}

func (e *literalExpr) value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags) {
	return e.val, nil
}

func (e *literalExpr) GoString() string {
	return fmt.Sprintf("eval.LiteralExpr(%#v, %#v)", e.val, e.rng.sourceRange())
}
