package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cbty"
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

func (cb StmtBlock) RequiredSymbols(scope *Scope) SymbolSet {
	ret := NewSymbolSet()
	for _, stmt := range cb.stmts {
		reqd := stmt.RequiredSymbols(scope)
		for sym := range reqd {
			ret.Add(sym)
		}
	}
	return ret
}

func (sb StmtBlock) PackagesImported() []PackageRef {
	return sb.PackagesImportedAppend(nil)
}

func (sb StmtBlock) PackagesImportedAppend(ppaths []PackageRef) []PackageRef {
	for _, stmt := range sb.stmts {
		if imp, isImp := stmt.s.(*importStmt); isImp {
			ppaths = append(ppaths, PackageRef{
				Path:  imp.ppath,
				Range: imp.sourceRange(),
			})
		}
	}
	return ppaths
}

type PackageRef struct {
	Path  string
	Range source.Range
}

type StmtBlockAttr struct {
	Symbol   *Symbol
	Type     cbty.Type
	Default  cbty.Value
	DefRange source.Range
}

type StmtBlockAttrs map[string]StmtBlockAttr

// AttributeNames returns a map of the attribute names required by the block,
// which can be resolved without needing a context. This can be used for
// early validation, though Attributes should be called with a context to
// get the full description of the attributes.
func (sb StmtBlock) AttributeNames() map[string]*Symbol {
	ret := map[string]*Symbol{}
	for _, stmt := range sb.stmts {
		if attr, isAttr := stmt.s.(*attrStmt); isAttr {
			ret[attr.sym.DeclaredName()] = attr.sym
		}
	}
	return ret
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
func (sb StmtBlock) Attributes(ctx *Context) (StmtBlockAttrs, source.Diags) {
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

				if tyVal.Type() != cbty.TypeType {
					diags = append(diags, source.Diag{
						Level:   source.Error,
						Summary: "Invalid attribute type",
						Detail:  fmt.Sprintf("Expected a type, but given a value of type %s. To assign a default value, use the '=' (equals) symbol.", tyVal.Type().Name()),
						Ranges:  attr.valueType.sourceRange().List(),
					})
					def.Type = cbty.PlaceholderVal.Type()
					def.Default = cbty.NilValue
				} else {
					def.Type = tyVal.UnwrapType()
					def.Default = cbty.NilValue
				}
			case attr.defValue != NilExpr:
				val, exprDiags := attr.defValue.Value(ctx)
				diags = append(diags, exprDiags...)

				def.Type = val.Type()
				def.Default = val
			default:
				// should never happen
				panic("attrStmt with neither value type nor default value")
			}

			ret[name] = def
		}
	}
	return ret, diags
}

// ImplicitExports returns a SymbolSet of the symbols defined in the block's
// scope that are eligible to be included in an implicit export object.
//
// This includes most definitions, but specifically excludes imports as they are
// assumed to be for internal use and could be requested directly by any caller.
//
// This method is intended for creating an implicit export object for a module,
// and so it will likely not produce a useful or sensible result for blocks
// created in other contexts.
func (sb StmtBlock) ImplicitExports() SymbolSet {
	var ret SymbolSet
	for _, stmt := range sb.stmts {
		if _, isImport := stmt.s.(*importStmt); isImport {
			continue
		}
		sym := stmt.s.definedSymbol()
		if sym != nil {
			if ret == nil {
				ret = SymbolSet{}
			}
			ret.Add(sym)
		}
	}
	return ret
}

func (sb StmtBlock) Execute(exec StmtBlockExecute, initDefs map[*Symbol]cbty.Value) (*StmtBlockResult, source.Diags) {
	// Make a new child context to work in. (We get "exec" by value here, so
	// we can mutate it without upsetting the caller.)
	exec.Context = exec.Context.NewChild()

	for sym, val := range initDefs {
		exec.Context.DefineLiteral(sym, val)
	}

	result := StmtBlockResult{}
	var diags source.Diags

	result.Scope = sb.scope
	result.Context = exec.Context

	for _, stmt := range sb.stmts {
		stmtDiags := stmt.s.execute(&exec, &result)
		diags = append(diags, stmtDiags...)
	}

	result.Context.final = true // no more modifications allowed

	return &result, diags
}

type StmtBlockExecute struct {
	Context  *Context
	Packages map[string]cbty.Value
	Attrs    map[*Symbol]cbty.Value
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

	// Designator is the value passed to a "Designator" statement, if present.
	// The empty string if not present.
	//
	// If errors are encountered processing a given designator, this string may
	// have the special value "<invalid>". If the result depends on a value
	// that isn't yet known (due to an error or otherwise), this string will
	// have the special value "<unknown>". Neither special value should be
	// seen for a valid program where all input values are known.
	Designator string

	// ExportValue is the value exported by an "export" statement, if any.
	ExportValue cbty.Value
}

func (a StmtBlockAttrs) CallSignature(posParams PosParameters, result cbty.Type) (*cbty.CallSignature, source.Diags) {
	var diags source.Diags
	sig := &cbty.CallSignature{}
	sig.Parameters = map[string]cbty.CallParameter{}
	sig.Result = result
	for name, attr := range a {
		sig.Parameters[name] = cbty.CallParameter{
			Type:     attr.Type,
			Required: attr.Default == cbty.NilValue,
		}
	}
	for _, param := range posParams {
		if _, exists := a[param.Name]; !exists {
			// Should never actually get here, because we should catch this
			// during compilation of posParams. We'll ignore it so we can
			// still produce a result, albeit an incomplete one.
			continue
		}
		sig.Positional = append(sig.Positional, param.Name)
	}
	return sig, diags
}
