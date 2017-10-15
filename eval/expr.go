package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

type Expr struct {
	e exprImpl
}

func (e Expr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
	return e.e.value(ctx, targetSym)
}

func (e Expr) eachChild(cb walkCb) {
	e.e.eachChild(cb)
}

func (e Expr) sourceRange() source.Range {
	return e.e.sourceRange()
}

// NilExpr is an invalid expression that serves as the zero value of Expr.
//
// NilExpr indicates the absense of an expression and is not itself a valid
// expression. Any methods called on it will panic.
var NilExpr Expr

type exprImpl interface {
	// value evaluates the expression in the given context.
	//
	// If the result is being used directly as the definition of a symbol
	// then the that symbol is provided in "targetSym"; otherwise,
	// targetSym is nil. Some expression types are valid only when being
	// assigned directly to a symbol.
	value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags)

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
	return Expr{&symbolExpr{
		sym: sym,
		rng: srcRange(rng),
	}}
}

func (e *symbolExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
	val := ctx.Value(e.sym)
	if val == cbty.NilValue {
		// This is actually an implementation error in Cirbo rather than a
		// user error, but we'll return it as a diagnostic anyway since that's
		// more graceful.
		return cbty.PlaceholderVal, source.Diags{
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
	val cbty.Value
	rng
	leafExpr
}

func LiteralExpr(val cbty.Value, rng source.Range) Expr {
	return Expr{&literalExpr{
		val: val,
		rng: srcRange(rng),
	}}
}

func (e *literalExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
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
	return Expr{&binaryOpExpr{
		lhs: lhs,
		rhs: rhs,
		op:  op,
		rng: srcRange(rng),
	}}
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

func (e *binaryOpExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
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
	return Expr{&unaryOpExpr{
		val: val,
		op:  op,
		rng: srcRange(rng),
	}}
}

func NegateExpr(val Expr, rng source.Range) Expr {
	return makeUnaryOpExpr(val, opNegate, rng)
}

func NotExpr(val Expr, rng source.Range) Expr {
	return makeUnaryOpExpr(val, opNot, rng)
}

func (e *unaryOpExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
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
	return Expr{&callExpr{
		callee:    callee,
		posArgs:   posArgs,
		namedArgs: namedArgs,
		rng:       srcRange(rng),
	}}
}

func (e *callExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
	callee, diags := e.callee.Value(ctx)
	if diags.HasErrors() {
		return cbty.PlaceholderVal, diags
	}

	sig := callee.Type().CallSignature()
	if sig == nil {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Value is not callable",
			Detail:  fmt.Sprintf("A value of type %s cannot be called.", callee.Type().Name()),
			Ranges:  e.sourceRange().List(),
		})
		return cbty.PlaceholderVal, diags
	}

	call := cbty.CallArgs{
		Explicit:  map[string]cbty.Value{},
		Context:   ctx,
		CallRange: e.sourceRange(),
	}

	if targetSym != nil {
		call.TargetName = targetSym.DeclaredName()
	}

	if len(e.posArgs) < len(sig.Positional) {
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Insufficient positional arguments",
			Detail:  fmt.Sprintf("This function requires %d positional arguments.", len(sig.Positional)),
			Ranges:  e.sourceRange().List(),
		})
		return cbty.UnknownVal(sig.Result), diags
	}
	if len(e.posArgs) > len(sig.Positional) && !sig.AcceptsVariadicPositional {
		extras := e.posArgs[len(sig.Positional):]
		extrasRange := source.RangeBetween(extras[0].sourceRange(), extras[len(extras)-1].sourceRange())
		diags = append(diags, source.Diag{
			Level:   source.Error,
			Summary: "Extraneous positional arguments",
			Detail:  fmt.Sprintf("This function requires %d positional arguments.", len(sig.Positional)),
			Ranges:  extrasRange.List(),
		})
		return cbty.UnknownVal(sig.Result), diags
	}

	for i, name := range sig.Positional {
		argExpr := e.posArgs[i]
		val, valDiags := argExpr.Value(ctx)
		diags = append(diags, valDiags...)
		call.Explicit[name] = val
	}

	if len(e.posArgs) > len(sig.Positional) {
		extras := e.posArgs[len(sig.Positional):]
		call.PosVariadic = make([]cbty.Value, len(extras))
		for i, argExpr := range extras {
			var valDiags source.Diags
			call.PosVariadic[i], valDiags = argExpr.Value(ctx)
			diags = append(diags, valDiags...)
		}
	}

	if sig.AcceptsVariadicNamed {
		call.NamedVariadic = make(map[string]cbty.Value)
	}

	for name, argExpr := range e.namedArgs {
		_, isExplicit := sig.Parameters[name]
		if isExplicit {
			if _, defined := call.Explicit[name]; defined {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Duplicate argument definition",
					Detail:  fmt.Sprintf("Argument %q has already been assigned a value within this call.", name),
					Ranges:  argExpr.sourceRange().List(),
				})
				continue
			}
			val, valDiags := argExpr.Value(ctx)
			diags = append(diags, valDiags...)
			call.Explicit[name] = val
		} else {
			if call.NamedVariadic == nil {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Extraneous argument",
					Detail:  fmt.Sprintf("This function does not expect an argument named %q.", name),
					Ranges:  argExpr.sourceRange().List(),
				})
				continue
			}
			if _, defined := call.NamedVariadic[name]; defined {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Duplicate argument definition",
					Detail:  fmt.Sprintf("Argument %q has already been assigned a value within this call.", name),
					Ranges:  argExpr.sourceRange().List(),
				})
				continue
			}
			val, valDiags := argExpr.Value(ctx)
			diags = append(diags, valDiags...)
			call.NamedVariadic[name] = val
		}
	}

	// If we encountered any errors during argument processing then we won't
	// actually try the call, since we'll probably end up just passing the
	// callee garbage that it won't be able to deal with.
	if diags.HasErrors() {
		return cbty.UnknownVal(sig.Result), diags
	}

	// Make sure we got all the arguments we needed and that they are of
	// the required types.
	for name, def := range sig.Parameters {
		val, defined := call.Explicit[name]
		if !defined {
			if def.Required {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Missing required argument",
					Detail:  fmt.Sprintf("This function requires an argument named %q.", name),
					Ranges:  e.sourceRange().List(),
				})
			}
			continue
		}

		if !val.Type().Same(def.Type) {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Incorrect argument type",
				Detail:  fmt.Sprintf("Argument %q must be of type %s, not %s.", name, def.Type.Name(), val.Type().Name()),
				Ranges:  e.sourceRange().List(),
			})
			continue
		}
	}

	if diags.HasErrors() {
		return cbty.UnknownVal(sig.Result), diags
	}

	result, callDiags := callee.Call(call)
	diags = append(diags, callDiags...)
	return result, diags
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
	return Expr{&attrExpr{
		obj:  obj,
		name: name,
		rng:  srcRange(rng),
	}}
}

func (e *attrExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
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
	return Expr{&indexExpr{
		coll:  coll,
		index: index,
		rng:   srcRange(rng),
	}}
}

func (e *indexExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
	panic("indexExpr.value not yet implemented")
}

func (e *indexExpr) eachChild(cb walkCb) {
	cb(e.coll)
	cb(e.index)
}

type passthroughExpr struct {
	expr Expr
	rng
}

func PassthroughExpr(expr Expr, rng source.Range) Expr {
	return Expr{&passthroughExpr{
		expr: expr,
		rng:  srcRange(rng),
	}}
}

func (e *passthroughExpr) value(ctx *Context, targetSym *Symbol) (cbty.Value, source.Diags) {
	// This is the one situation where we _do_ pass through a targetSym
	// value, since we want PassthroughExpr to act like it isn't there at all.
	// (It's only there to represent parethesized exprs so we can draw ranges
	// properly around them in the event of errors.)
	return e.expr.value(ctx, targetSym)
}

func (e *passthroughExpr) eachChild(cb walkCb) {
	cb(e.expr)
}

func (e *passthroughExpr) GoString() string {
	return fmt.Sprintf("eval.PassthroughExpr(%#v)", e.expr)
}
