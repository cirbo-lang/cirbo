package eval

import (
	"bytes"
	"fmt"
)

type StmtSet map[Stmt]struct{}

func NewStmtSet(stmts ...Stmt) StmtSet {
	ret := make(StmtSet, len(stmts))
	for _, stmt := range stmts {
		ret.Add(stmt)
	}
	return ret
}

func (s StmtSet) Has(stmt Stmt) bool {
	_, has := s[stmt]
	return has
}

func (s StmtSet) Add(stmt Stmt) {
	s[stmt] = struct{}{}
}

func (s StmtSet) Remove(stmt Stmt) {
	delete(s, stmt)
}

func (s StmtSet) Equal(o StmtSet) bool {
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

func (s StmtSet) GoString() string {
	buf := bytes.Buffer{}
	buf.WriteString("eval.NewStmtSet(")
	first := true
	for stmt := range s {
		if !first {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%#v", stmt)
		first = false
	}
	buf.WriteString(")")
	return buf.String()
}
