package eval

import (
	"bytes"
	"fmt"
)

type ExprSet map[Expr]struct{}

func NewExprSet(exprs ...Expr) ExprSet {
	ret := make(ExprSet, len(exprs))
	for _, expr := range exprs {
		ret.Add(expr)
	}
	return ret
}

func (s ExprSet) Has(expr Expr) bool {
	_, has := s[expr]
	return has
}

func (s ExprSet) Add(expr Expr) {
	s[expr] = struct{}{}
}

func (s ExprSet) Remove(expr Expr) {
	delete(s, expr)
}

func (s ExprSet) Equal(o ExprSet) bool {
	if len(s) != len(o) {
		return false
	}

	for expr := range s {
		if !o.Has(expr) {
			return false
		}
	}

	return true
}

func (s ExprSet) GoString() string {
	buf := bytes.Buffer{}
	buf.WriteString("eval.NewExprSet(")
	first := true
	for expr := range s {
		if !first {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%#v", expr)
		first = false
	}
	buf.WriteString(")")
	return buf.String()
}
