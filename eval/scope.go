package eval

import (
	"fmt"
	"sort"
)

// Scope represents a symbol table for evaluation. Each symbol appearing in
// an expression belongs to exactly one scope.
type Scope struct {
	parent  *Scope
	symbols map[string]*Symbol
}

// Symbol represents a particular symbol used in expressions. Symbol represents
// the location where a value will be stored rather than the value itself.
//
// A value for each symbol is defined in an associated Context.
type Symbol struct {
	name  string
	scope *Scope
}

// NewChild creates a new, empty scope that is a child of the receiver.
func (s *Scope) NewChild() *Scope {
	return &Scope{
		parent:  s,
		symbols: map[string]*Symbol{},
	}
}

// Declare establishes a new symbol in the recieving scope and returns it.
//
// This method will panic if a symbol of the same name was already declared in
// the receiving scope. This does not apply if the symbol is declared in an
// ancestor scope. Use Declared to determine if a given name was already declared.
//
// This method will also panic if a caller attempts to declare a symbol in the
// global scope, since the global scope is an immutable singleton.
func (s *Scope) Declare(name string) *Symbol {
	if s == globalScope {
		panic(fmt.Errorf("attempt to declare %q in the immutable global scope", name))
	}
	if s.symbols[name] != nil {
		panic(fmt.Errorf("attempt to re-declare %q in %#v", name, s))
	}
	sym := &Symbol{
		name:  name,
		scope: s,
	}
	s.symbols[name] = sym
	return sym
}

// Declared returns true if a symbol of the given name was already declared in
// the receiving scope. Returns false if the name is not defined or if it is
// declared only in ancestor scopes.
func (s *Scope) Declared(name string) bool {
	_, declared := s.symbols[name]
	return declared
}

// Get returns the symbol with the given name in either the receiver or
// its closest ancestor that declares the given name.
//
// If no scope in the inheritance chain declares the given name, the result
// is nil.
func (s *Scope) Get(name string) *Symbol {
	current := s
	for current != nil {
		sym := current.symbols[name]
		if sym != nil {
			return sym
		}
		current = current.parent
	}
	return nil
}

// Parent returns the receiver's parent scope, or nil if the receiver is the
// global scope.
func (s *Scope) Parent() *Scope {
	return s.parent
}

// AllNames returns a slice of all of the names declared in the receiver and
// all of its ancestor scopes.
//
// This is primarily intended for giving feedback to the user in error messages
// about references to undeclared symbols. It is returned sorted first by
// scope in order of "closeness" (receiver first, global scope last) and then
// by name lexicographically within each scope.
//
// Due to the sorting strategy, the same name may appear multiple times in the
// list if it is declared in multiple different scopes.
func (s *Scope) AllNames() []string {
	if s == nil {
		return nil
	}

	inherited := s.parent.AllNames()
	ret := make([]string, 0, len(s.symbols)+len(inherited))
	for name := range s.symbols {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	ret = append(ret, inherited...)

	return ret
}

// DeclaredName returns the name that was used at the declaration of the
// receiving symbol.
//
// Symbol names are not globally unique since child scopes can shadow
// declarations in parent scopes; while this result can be useful to
// talk about scopes to the user, care must be taken to give enough context
// to avoid creating further confusion.
func (sym *Symbol) DeclaredName() string {
	return sym.name
}

// Scope returns the scope that the receiving symbol belongs to.
func (sym *Symbol) Scope() *Scope {
	return sym.scope
}
