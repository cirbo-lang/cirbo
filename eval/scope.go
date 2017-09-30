package eval

import (
	"fmt"
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
