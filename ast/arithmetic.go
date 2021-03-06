package ast

type ArithmeticBinary struct {
	WithRange
	LHS Node
	RHS Node
	Op  ArithmeticOp
}

type ArithmeticUnary struct {
	WithRange
	Operand Node
	Op      ArithmeticOp
}

type ArithmeticOp rune

//go:generate stringer -type=ArithmeticOp

const (
	Equal              ArithmeticOp = '='
	NotEqual           ArithmeticOp = '≠'
	LessThan           ArithmeticOp = '<'
	GreaterThan        ArithmeticOp = '>'
	LessThanOrEqual    ArithmeticOp = '≤'
	GreaterThanOrEqual ArithmeticOp = '≥'

	Add      ArithmeticOp = '+'
	Subtract ArithmeticOp = '-'
	Multiply ArithmeticOp = '×'
	Divide   ArithmeticOp = '÷'
	Modulo   ArithmeticOp = 'm' // written as 'mod' because % is used for percentages
	Exponent ArithmeticOp = '^'
	Negate   ArithmeticOp = '±'

	Concat ArithmeticOp = '…'

	And ArithmeticOp = '∧'
	Or  ArithmeticOp = '∨'
	Not ArithmeticOp = '¬'

	ArithmeticOpNil ArithmeticOp = 0
)

func (n *ArithmeticBinary) walkChildNodes(cb internalWalkFunc) {
	cb(n.LHS)
	cb(n.RHS)
}

func (n *ArithmeticUnary) walkChildNodes(cb internalWalkFunc) {
	cb(n.Operand)
}
