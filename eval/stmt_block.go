package eval

import (
	"github.com/cirbo-lang/cirbo/source"
)

type StmtBlock struct {
	stmts []Stmt
}

// NewStmtBlock constructs and returns a new statement block containing the
// given statements.
//
// The caller must not read or write the given statements slice after it has
// been passed to NewStmtBlock. Ownership is transferred to the returned
// object and the slice's backing array may be modified in unspecified ways.
func MakeStmtBlock(stmts []Stmt) (StmtBlock, source.Diags) {
	// TODO: do an in-place topological sort on the provided statements, and
	// return a diagnostic describing a cycle if one is found.
	return StmtBlock{
		stmts: stmts,
	}, nil
}
