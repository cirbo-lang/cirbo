package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/source"
)

type StmtBlock struct {
	scope *Scope
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

	// Don't permit any future modifications to the scope, since we're now
	// depending on its contents.
	scope.final = true

	return StmtBlock{
		scope: scope,
		stmts: stmts,
	}, diags
}

func (sb StmtBlock) ModulesImported() []string {
	return sb.ModulesImportedAppend(nil)
}

func (sb StmtBlock) ModulesImportedAppend(ppaths []string) []string {
	for _, stmt := range sb.stmts {
		if imp, isImp := stmt.s.(*importStmt); isImp {
			ppaths = append(ppaths, imp.ppath)
		}
	}
	return ppaths
}

type StmtBlockAttr struct {
	Symbol   *Symbol
	Type     cty.Type
	Required bool
	DefRange source.Range
}

// Attributes returns a description of the attributes defined by the block.
//
// When executing the block, values for some or all of these (depending on
// their Required status) should be provided in the StmtBlockExecute
// instance.
//
// The given context is used to resolve the type or default value expressions
// in the attribute statements. The given context must therefore be the same
// context that would ultimately be provided to Execute in the StmtBlockExecute
// object or else the result may be incorrect.
func (sb StmtBlock) Attributes(ctx *Context) (map[string]StmtBlockAttr, source.Diags) {
	var diags source.Diags
	ret := map[string]StmtBlockAttr{}
	for _, stmt := range sb.stmts {
		if attr, isAttr := stmt.s.(*attrStmt); isAttr {
			name := attr.sym.DeclaredName()
			def := StmtBlockAttr{
				Symbol: attr.sym,
			}

			switch {
			case attr.valueType != NilExpr:
				tyVal, exprDiags := attr.valueType.Value(ctx)
				diags = append(diags, exprDiags...)

				if tyVal.Type() != cty.TypeType {
					diags = append(diags, source.Diag{
						Level:   source.Error,
						Summary: "Invalid attribute type",
						Detail:  fmt.Sprintf("Expected a type, but given a value of type %s. To assign a default value, use the '=' (equals) symbol.", tyVal.Type().Name()),
						Ranges:  attr.valueType.sourceRange().List(),
					})
					def.Type = cty.PlaceholderVal.Type()
					def.Required = true
				} else {
					def.Type = tyVal.UnwrapType()
					def.Required = true
				}
			case attr.defValue != NilExpr:
				val, exprDiags := attr.defValue.Value(ctx)
				diags = append(diags, exprDiags...)

				def.Type = val.Type()
				def.Required = false
			default:
				// should never happen
				panic("attrStmt with neither value type nor default value")
			}

			ret[name] = def
		}
	}
	return ret, diags
}

func (sb StmtBlock) Execute(exec StmtBlockExecute) (*StmtBlockResult, source.Diags) {
	// Make a new child context to work in. (We get "exec" by value here, so
	// we can mutate it without upsetting the caller.)
	exec.Context = exec.Context.NewChild()

	result := StmtBlockResult{}
	var diags source.Diags

	result.Scope = sb.scope
	result.Context = exec.Context

	for _, stmt := range sb.stmts {
		stmtDiags := stmt.s.execute(&exec, &result)
		diags = append(diags, stmtDiags...)
	}

	result.Values = result.Context.AllValues(sb.scope)
	result.Context.final = true // no more modifications allowed

	return &result, diags
}

type StmtBlockExecute struct {
	Context *Context
	Modules map[string]cty.Value
	Attrs   map[*Symbol]cty.Value
}

type StmtBlockResult struct {
	// Context is the context that was created to evaluate the block.
	// It is provided only so that NewChild may be called on it for child
	// blocks; it should not be modified once returned.
	Context *Context

	// Scope is the scope that was created for the block during its
	// compilation. This object is shared between all executions of the same
	// block, and so should not be modified.
	Scope *Scope

	// Values gives the local values that were defined during execution.
	// These values are created from the perspective of the "inside" of the
	// block and so should not generally be exposed to expressions outside of
	// the block.
	Values map[string]cty.Value
}
