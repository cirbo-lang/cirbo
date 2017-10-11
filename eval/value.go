package eval

import (
	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

// Value evaluates the given expression in the given context and returns the
// result, along with any diagnostics that are produced during the operation.
func (expr Expr) Value(ctx *Context) (cbty.Value, source.Diags) {
	return expr.value(ctx, nil)
}
