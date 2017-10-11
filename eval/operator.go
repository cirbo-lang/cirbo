package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

type operator rune

//go:generate stringer -type=operator

const (
	opNone operator = 0

	// Arithmetic
	opAdd      operator = '+'
	opSubtract operator = '-'
	opMultiply operator = '*'
	opDivide   operator = '/'
	opModulo   operator = '%'
	opExponent operator = '^'
	opNegate   operator = '±'

	// Concatenation
	opConcat operator = '…'

	// Equality
	opEqual    operator = '='
	opNotEqual operator = '≠'

	// Comparison
	opLessThan           operator = '<'
	opGreaterThan        operator = '>'
	opLessThanOrEqual    operator = '≤'
	opGreaterThanOrEqual operator = '≥'

	// Boolean Logic
	opAnd operator = '∧'
	opOr  operator = '∨'
	opNot operator = '¬'
)

func (o operator) evalBinary(ctx *Context, lhs, rhs Expr, rng source.Range) (cbty.Value, source.Diags) {
	var diags source.Diags
	lv, lhsDiags := lhs.value(ctx, nil)
	diags = append(diags, lhsDiags...)
	rv, rhsDiags := rhs.value(ctx, nil)
	diags = append(diags, rhsDiags...)

	invalidTypes := func() source.Diag {
		return source.Diag{
			Level:   source.Error,
			Summary: "Invalid operand types",
			Detail:  fmt.Sprintf("Cannot %s with %s values.", o.verb(), typePairStr(lv.Type(), rv.Type())),
		}
	}

	switch o {
	case opAdd, opSubtract:
		if !lv.Type().CanSum(rv.Type()) {
			diags = append(diags, invalidTypes())
			return cbty.PlaceholderVal, diags
		}
		switch o {
		case opAdd:
			return lv.Add(rv), diags
		case opSubtract:
			return lv.Subtract(rv), diags
		default:
			panic("invalid sum operator")
		}
	case opMultiply, opDivide, opModulo, opExponent:
		if !lv.Type().CanProduct(rv.Type()) {
			diags = append(diags, invalidTypes())
			return cbty.PlaceholderVal, diags
		}
		switch o {
		case opMultiply:
			return lv.Multiply(rv), diags
		case opDivide:
			return lv.Divide(rv), diags
		case opModulo:
			// TODO: implement
			panic("modulo not yet implemented")
		case opExponent:
			// TODO: implement
			panic("exponent not yet implemented")
		default:
			panic("invalid product operator")
		}
	case opConcat:
		if !lv.Type().CanConcat(rv.Type()) {
			diags = append(diags, invalidTypes())
			return cbty.PlaceholderVal, diags
		}
		return lv.Concat(rv), diags
	case opEqual:
		return lv.Equal(rv), diags
	case opNotEqual:
		return lv.Equal(rv).Not(), diags
	case opLessThan, opLessThanOrEqual, opGreaterThan, opGreaterThanOrEqual:
		panic("comparison not yet implemented")
	case opAnd, opOr:
		if !(lv.Type().Same(cbty.Bool) && rv.Type().Same(cbty.Bool)) {
			diags = append(diags, invalidTypes())
			return cbty.UnknownVal(cbty.Bool), diags
		}
		switch o {
		case opAnd:
			return lv.And(rv), diags
		case opOr:
			return lv.Or(rv), diags
		default:
			panic("invalid boolean operator")
		}
	default:
		panic(fmt.Errorf("%s is not a binary operator", o))
	}
}

func (o operator) evalUnary(ctx *Context, val Expr, rng source.Range) (cbty.Value, source.Diags) {
	vv, diags := val.value(ctx, nil)

	switch o {
	case opNegate:
		panic("negate not yet implemented")
	case opNot:
		if !vv.Type().Same(cbty.Bool) {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Invalid operand type",
				Detail:  fmt.Sprintf("Cannot %s with a %s value.", o.verb(), vv.Type().Name()),
			})
			return cbty.UnknownVal(cbty.Bool), diags
		}
		return vv.Not(), diags
	default:
		panic(fmt.Errorf("%s is not a unary operator", o))
	}
}

func (o operator) verb() string {
	switch o {
	case opAdd:
		return "add"
	case opSubtract:
		return "subtract"
	case opMultiply:
		return "multiply"
	case opDivide:
		return "divide"
	case opModulo:
		return "modulo-divide"
	case opExponent:
		return "exponentiate"
	case opNegate:
		return "negate"
	case opConcat:
		return "concatenate"
	case opEqual, opNotEqual:
		return "equate"
	case opLessThan, opGreaterThan, opLessThanOrEqual, opGreaterThanOrEqual:
		return "compare"
	case opAnd:
		return "apply AND"
	case opOr:
		return "apply OR"
	case opNot:
		return "apply NOT"
	default:
		return "evaluate" // should never hit this because above should be comprehensive
	}
}

func typePairStr(a, b cbty.Type) string {
	if a.Same(b) {
		return fmt.Sprintf("two %s", a.Name())
	} else {
		return fmt.Sprintf("%s and %s", a.Name(), b.Name())
	}
}
