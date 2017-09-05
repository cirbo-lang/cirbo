package parser

import (
	"github.com/cirbo-lang/cirbo/ast"
)

// This operation table maps from the operator's token type
// to the AST operation type. All expressions produced from
// binary operators are BinaryOp nodes.
//
// Binary operator groups are listed in order of precedence, with
// the *lowest* precedence first. Operators within the same group
// have left-to-right associativity.
var binaryOps = []map[TokenType]ast.ArithmeticOp{
	{
		TokenOr: ast.Or,
	},
	{
		TokenAnd: ast.And,
	},
	{
		TokenEqual:    ast.Equal,
		TokenNotEqual: ast.NotEqual,
	},
	{
		TokenGreaterThan:   ast.GreaterThan,
		TokenGreaterThanEq: ast.GreaterThanOrEqual,
		TokenLessThan:      ast.LessThan,
		TokenLessThanEq:    ast.LessThanOrEqual,
	},
	{
		TokenPlus:   ast.Add,
		TokenMinus:  ast.Subtract,
		TokenDotDot: ast.Concat,
	},
	{
		TokenStar:  ast.Multiply,
		TokenSlash: ast.Divide,
		// TODO: Modulo? It doesn't have a punctuation token associated with
		// it so it doesn't really fit into this model.
	},
}
