package eval

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

func TestNewStmtBlock(t *testing.T) {
	type testCase struct {
		Stmts []Stmt
		Want  []Stmt
		Diags int
	}

	tests := map[string]func(scope *Scope) testCase{
		"empty": func(scope *Scope) testCase {
			return testCase{
				Stmts: []Stmt{},
				Want:  []Stmt{},
				Diags: 0,
			}
		},
		"single": func(scope *Scope) testCase {
			sym := scope.Declare("sym")
			stmt := makeMockStmt(sym, nil)

			return testCase{
				Stmts: []Stmt{stmt},
				Want:  []Stmt{stmt},
				Diags: 0,
			}
		},
		"simple": func(scope *Scope) testCase {
			sym := scope.Declare("sym")
			definer := makeMockStmt(sym, nil)
			user := makeMockStmt(nil, NewSymbolSet(sym))

			return testCase{
				Stmts: []Stmt{user, definer},
				Want:  []Stmt{definer, user},
				Diags: 0,
			}
		},
		"chain": func(scope *Scope) testCase {
			sym1 := scope.Declare("sym1")
			sym2 := scope.Declare("sym2")
			definer := makeMockStmt(sym1, nil)
			userDefiner := makeMockStmt(sym2, NewSymbolSet(sym1))
			user := makeMockStmt(nil, NewSymbolSet(sym2))

			return testCase{
				Stmts: []Stmt{user, definer, userDefiner},
				Want:  []Stmt{definer, userDefiner, user},
				Diags: 0,
			}
		},
		"fork": func(scope *Scope) testCase {
			sym1 := scope.Declare("sym1")
			sym2 := scope.Declare("sym2")
			sym3 := scope.Declare("sym3")
			definer := makeMockStmt(sym1, nil)
			userDefiner1 := makeMockStmt(sym2, NewSymbolSet(sym1))
			userDefiner2 := makeMockStmt(sym3, NewSymbolSet(sym1))
			user := makeMockStmt(nil, NewSymbolSet(sym2, sym3))

			return testCase{
				// userDefiner2 and userDefiner1 can be handled in any order,
				// so the sort should preserve the input ordering and place
				// userDefiner2 first.
				Stmts: []Stmt{user, definer, userDefiner2, userDefiner1},
				Want:  []Stmt{definer, userDefiner2, userDefiner1, user},
				Diags: 0,
			}
		},
		"mutually-dependent": func(scope *Scope) testCase {
			sym1 := scope.Declare("sym1")
			sym2 := scope.Declare("sym2")
			stmtA := makeMockStmt(sym1, NewSymbolSet(sym2))
			stmtB := makeMockStmt(sym2, NewSymbolSet(sym1))
			stmtC := makeMockStmt(nil, NewSymbolSet(sym1))

			return testCase{
				Stmts: []Stmt{stmtB, stmtA, stmtC},
				Want:  []Stmt{},
				Diags: 1, // dependency cycle
			}
		},
		"self-reference": func(scope *Scope) testCase {
			sym1 := scope.Declare("sym1")
			stmt := makeMockStmt(sym1, NewSymbolSet(sym1))

			return testCase{
				Stmts: []Stmt{stmt},
				Want:  []Stmt{},
				Diags: 1, // dependency cycle
			}
		},
		"long cycle": func(scope *Scope) testCase {
			sym1 := scope.Declare("sym1")
			sym2 := scope.Declare("sym2")
			sym3 := scope.Declare("sym3")
			stmtA := makeMockStmt(sym1, NewSymbolSet(sym3))
			stmtB := makeMockStmt(sym2, NewSymbolSet(sym1))
			stmtC := makeMockStmt(sym3, NewSymbolSet(sym2))

			return testCase{
				Stmts: []Stmt{stmtB, stmtA, stmtC},
				Want:  []Stmt{},
				Diags: 1, // dependency cycle
			}
		},
	}

	for name, cons := range tests {
		t.Run(name, func(t *testing.T) {
			scope := globalScope.NewChild()
			test := cons(scope)
			gotBlock, diags := MakeStmtBlock(scope, test.Stmts)
			assertDiagCount(t, diags, test.Diags)

			got := gotBlock.stmts
			bad := false
			if len(got) == len(test.Want) {
				for i := range test.Want {
					if got[i] != test.Want[i] {
						bad = true
					}
				}
			} else {
				bad = true
			}

			if bad {
				t.Fatalf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestStmtBlockAttributes(t *testing.T) {
	type testCase struct {
		Stmts []Stmt
		Want  StmtBlockAttrs
		Diags int
	}

	tests := map[string]func(scope *Scope) testCase{
		"empty": func(scope *Scope) testCase {
			return testCase{
				Stmts: []Stmt{},
				Want:  StmtBlockAttrs{},
				Diags: 0,
			}
		},
		"irrelevant statements": func(scope *Scope) testCase {
			return testCase{
				Stmts: []Stmt{
					makeMockStmt(nil, nil),
					makeMockStmt(nil, nil),
				},
				Want:  StmtBlockAttrs{},
				Diags: 0,
			}
		},
		"simple required": func(scope *Scope) testCase {
			sym := scope.Declare("foo")
			tyExpr := SymbolExpr(scope.Get("Length"), source.NilRange)
			return testCase{
				Stmts: []Stmt{
					AttrStmt(sym, tyExpr, source.NilRange),
				},
				Want: StmtBlockAttrs{
					"foo": {
						Symbol:   sym,
						Type:     cbty.Length,
						Required: true,
					},
				},
				Diags: 0,
			}
		},
		"simple optional": func(scope *Scope) testCase {
			sym := scope.Declare("foo")
			valExpr := SymbolExpr(scope.Get("Length"), source.NilRange)
			return testCase{
				Stmts: []Stmt{
					AttrStmtDefault(sym, valExpr, source.NilRange),
				},
				Want: StmtBlockAttrs{
					"foo": {
						Symbol:   sym,
						Type:     cbty.TypeType,
						Required: false,
					},
				},
				Diags: 0,
			}
		},
		"multiple": func(scope *Scope) testCase {
			sym1 := scope.Declare("sym1")
			sym2 := scope.Declare("sym2")
			expr := SymbolExpr(scope.Get("Length"), source.NilRange)
			return testCase{
				Stmts: []Stmt{
					AttrStmtDefault(sym1, expr, source.NilRange),
					AttrStmt(sym2, expr, source.NilRange),
				},
				Want: StmtBlockAttrs{
					"sym1": {
						Symbol:   sym1,
						Type:     cbty.TypeType,
						Required: false,
					},
					"sym2": {
						Symbol:   sym2,
						Type:     cbty.Length,
						Required: true,
					},
				},
				Diags: 0,
			}
		},
		"invalid type": func(scope *Scope) testCase {
			sym := scope.Declare("foo")
			tyExpr := LiteralExpr(cbty.StringVal("hello"), source.NilRange)
			return testCase{
				Stmts: []Stmt{
					AttrStmt(sym, tyExpr, source.NilRange),
				},
				Want: StmtBlockAttrs{
					"foo": {
						Symbol:   sym,
						Type:     cbty.PlaceholderVal.Type(),
						Required: true,
					},
				},
				Diags: 1, // Type expression evaluates to String, not Type
			}
		},
	}

	for name, cons := range tests {
		t.Run(name, func(t *testing.T) {
			scope := globalScope.NewChild()
			ctx := globalContext.NewChild()
			test := cons(scope)
			block, diags := MakeStmtBlock(scope, test.Stmts)

			got, attrDiags := block.Attributes(ctx)
			diags = append(diags, attrDiags...)

			assertDiagCount(t, diags, test.Diags)

			if !reflect.DeepEqual(got, test.Want) {
				t.Fatalf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

type mockStmt struct {
	defines  *Symbol
	requires SymbolSet
	rng
}

func makeMockStmt(defines *Symbol, requires SymbolSet) Stmt {
	return Stmt{&mockStmt{
		defines:  defines,
		requires: requires,
	}}
}

func (s *mockStmt) definedSymbol() *Symbol {
	return s.defines
}

func (s *mockStmt) requiredSymbols(*Scope) SymbolSet {
	return s.requires
}

func (s *mockStmt) execute(exec *StmtBlockExecute, result *StmtBlockResult) source.Diags {
	exec.Context.DefineLiteral(s.defines, cbty.PlaceholderVal)
	return nil
}
