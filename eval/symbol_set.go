package eval

import (
	"bytes"
	"fmt"
	"sort"
)

type SymbolSet map[*Symbol]struct{}

func NewSymbolSet(syms ...*Symbol) SymbolSet {
	ret := make(SymbolSet, len(syms))
	for _, sym := range syms {
		ret.Add(sym)
	}
	return ret
}

func (s SymbolSet) Has(sym *Symbol) bool {
	_, has := s[sym]
	return has
}

func (s SymbolSet) Add(sym *Symbol) {
	s[sym] = struct{}{}
}

func (s SymbolSet) Remove(sym *Symbol) {
	delete(s, sym)
}

func (s SymbolSet) Equal(o SymbolSet) bool {
	if len(s) != len(o) {
		return false
	}

	for sym := range s {
		if !o.Has(sym) {
			return false
		}
	}

	return true
}

func (s SymbolSet) GoString() string {
	buf := bytes.Buffer{}
	buf.WriteString("eval.NewSymbolSet(")
	first := true
	for sym := range s {
		if !first {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%#v", sym)
		first = false
	}
	buf.WriteString(")")
	return buf.String()
}

func (s SymbolSet) Names() []string {
	ret := make([]string, 0, len(s))
	for sym := range s {
		ret = append(ret, sym.name)
	}
	sort.Strings(ret)
	return ret
}
