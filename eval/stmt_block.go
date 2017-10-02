package eval

import (
	"github.com/cirbo-lang/cirbo/source"
)

type StmtBlock struct {
	stmts []Stmt
}

// NewStmtBlock constructs and returns a new statement block containing the
// given statements, which are assumed to be populating the given scope.
//
// The caller must not read or write the given statements slice after it has
// been passed to NewStmtBlock. Ownership is transferred to the returned
// object and the slice's backing array may be modified in unspecified ways.
func MakeStmtBlock(scope *Scope, stmts []Stmt) (StmtBlock, source.Diags) {
	var diags source.Diags

	providers := make(map[*Symbol]Stmt, len(stmts))
	enables := make(map[Stmt][]Stmt, len(stmts)) // slice so that we preserve input ordering when ordering is ambiguous
	inDeg := make(map[Stmt]int, len(stmts))
	for _, stmt := range stmts {
		if sym := stmt.s.definedSymbol(); sym != nil {
			providers[sym] = stmt
		}
	}
	for _, stmt := range stmts {
		syms := stmt.s.requiredSymbols(scope)
		for sym := range syms {
			if provider, provided := providers[sym]; provided {
				enables[provider] = append(enables[provider], stmt)
				inDeg[stmt]++
			}
		}
	}

	// We place both "result" and "queue" at the head of our input array.
	// We know that the length of the queue and the length of the result
	// must sum up to less than or equal to the original list, so we can
	// safely use the original underlying array as storage for both. The
	// start of the queue will gradually move through the array just as
	// the result slice grows to include the elements it has vacated.
	result := stmts[0:0]
	queueStart := 0 // index into stmts underlying array
	queue := stmts[queueStart:queueStart]

	// Seed the queue with statements that have no dependencies
	for _, stmt := range stmts {
		if inDeg[stmt] == 0 {
			queue = append(queue, stmt)
		}
	}

	for len(queue) > 0 {
		stmt := queue[0]

		// Adjust the head of the queue to one element later in our array.
		queueStart++
		queue = stmts[queueStart : queueStart+(len(queue)-1)]

		// Adjust the result list to include the element that we just
		// removed from the queue.
		result = stmts[:len(result)+1]

		for _, enabled := range enables[stmt] {
			inDeg[enabled]--
			if inDeg[enabled] == 0 {
				queue = append(queue, enabled)
				delete(inDeg, enabled)
			}
		}
	}

	// When we reach this point, if there were no cycles then result already
	// equals stmts, but the list may have shrunk if there _were_ cycles and
	// so we need to do some adjusting.
	stmts = result

	if len(inDeg) > 0 {
		// Indicates that we have at least one cycle.
		// TODO: This error message isn't great; ideally we would provide
		// more context to help the user understand the reason for the
		// cycle, since it might be via multiple levels of indirection.
		ranges := make([]source.Range, 0, len(inDeg))
		for stmt := range inDeg {
			ranges = append(ranges, stmt.s.sourceRange())
		}

		if len(ranges) == 1 {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Self-referential symbol definition",
				Detail:  "Definition statement depends (possibly indirectly) on its own result.",
				Ranges:  ranges,
			})
		} else {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Self-referential symbol definitions",
				Detail:  "Definition statements depend (possibly indirectly) on their own results.",
				Ranges:  ranges,
			})
		}
	}

	return StmtBlock{
		stmts: stmts,
	}, diags
}
