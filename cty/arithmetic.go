package cty

// typeWithArithmetic is an interface implemented by typeImpls for types that
// can do arithmetic.
type typeWithArithmetic interface {
	Add(a, b Value) Value
	Subtract(a, b Value) Value
	Multiply(a, b Value) Value
	Divide(a, b Value) Value
}
