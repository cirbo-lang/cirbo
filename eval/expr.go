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

type binaryOpExpr struct {
	lhs Expr
	rhs Expr
	op  operator
	rng
}

func makeBinaryOpExpr(lhs, rhs Expr, op operator, rng source.Range) Expr {
	return &binaryOpExpr{
		lhs: lhs,
		rhs: rhs,
		op:  op,
		rng: srcRange(rng),
	}
}

func AddExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opAdd, rng)
}

func SubtractExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opSubtract, rng)
}

func MultiplyExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opMultiply, rng)
}

func DivideExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opDivide, rng)
}

func ModuloExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opModulo, rng)
}

func ExponentExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opExponent, rng)
}

func ConcatExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opConcat, rng)
}

func EqualExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opEqual, rng)
}

func NotEqualExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opNotEqual, rng)
}

func LessThanExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opLessThan, rng)
}

func LessThanOrEqualExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opLessThanOrEqual, rng)
}

func GreaterThanExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opGreaterThan, rng)
}

func GreaterThanOrEqualExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opGreaterThanOrEqual, rng)
}

func AndExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opAnd, rng)
}

func OrExpr(lhs, rhs Expr, rng source.Range) Expr {
	return makeBinaryOpExpr(lhs, rhs, opOr, rng)
}

func (e *binaryOpExpr) value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags) {
	return e.op.evalBinary(ctx, e.lhs, e.rhs, e.sourceRange())
}

func (e *binaryOpExpr) eachChild(cb walkCb) {
	cb(e.lhs)
	cb(e.rhs)
}

func (e *binaryOpExpr) GoString() string {
	name := e.op.String()[2:]
	return fmt.Sprintf("eval.%sExpr(%#v, %#v, %#v)", name, e.lhs, e.rhs, e.sourceRange())
}

type unaryOpExpr struct {
	val Expr
	op  operator
	rng
}

func makeUnaryOpExpr(val Expr, op operator, rng source.Range) Expr {
	return &unaryOpExpr{
		val: val,
		op:  op,
		rng: srcRange(rng),
	}
}

func NegateExpr(val Expr, rng source.Range) Expr {
	return makeUnaryOpExpr(val, opNegate, rng)
}

func NotExpr(val Expr, rng source.Range) Expr {
	return makeUnaryOpExpr(val, opNot, rng)
}

func (e *unaryOpExpr) value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags) {
	return e.op.evalUnary(ctx, e.val, e.sourceRange())
}

func (e *unaryOpExpr) eachChild(cb walkCb) {
	cb(e.val)
}

type callExpr struct {
	callee    Expr
	posArgs   []Expr
	namedArgs map[string]Expr
	rng
}

func CallExpr(callee Expr, posArgs []Expr, namedArgs map[string]Expr, rng source.Range) Expr {
	return &callExpr{
		callee:    callee,
		posArgs:   posArgs,
		namedArgs: namedArgs,
		rng:       srcRange(rng),
	}
}

func (e *callExpr) value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags) {
	panic("callExpr.value not yet implemented")
}

func (e *callExpr) eachChild(cb walkCb) {
	cb(e.callee)
	for _, expr := range e.posArgs {
		cb(expr)
	}
	for _, expr := range e.namedArgs {
		cb(expr)
	}
}

type attrExpr struct {
	obj  Expr
	name string
	rng
}

func AttrExpr(obj Expr, name string, rng source.Range) Expr {
	return &attrExpr{
		obj:  obj,
		name: name,
		rng:  srcRange(rng),
	}
}

func (e *attrExpr) value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags) {
	panic("attrExpr.value not yet implemented")
}

func (e *attrExpr) eachChild(cb walkCb) {
	cb(e.obj)
}

type indexExpr struct {
	coll  Expr
	index Expr
	rng
}

func IndexExpr(coll, index Expr, rng source.Range) Expr {
	return &indexExpr{
		coll:  coll,
		index: index,
		rng:   srcRange(rng),
	}
}

func (e *indexExpr) value(ctx *Context, targetSym *Symbol) (cty.Value, source.Diags) {
	panic("indexExpr.value not yet implemented")
}

func (e *indexExpr) eachChild(cb walkCb) {
	cb(e.coll)
	cb(e.index)
}
