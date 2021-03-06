package cbo

import (
	"sort"
)

// An Endpoint is a low-level object representing a participant in a net.
type Endpoint struct {
	Name string
	Net  *Net
	ERC  ERCMode

	// Passthrough, if non-nil, is a set of endpoints that "pass through"
	// ERC characteristics.
	//
	// This is used for devices like resistors where the ERC direction
	// and output mode "pass through" for rules-checking purposes. It is also
	// used with terminals, once connected, to ensure that the outer net
	// is compatible with the inner net.
	//
	// When set, the endpoint's own ERC is ignored and one is instead inferred
	// by combining the ERC modes arriving on "the other side".
	Passthrough EndpointSet
}

// An EndpointSet is a set of endpoints.
type EndpointSet map[*Endpoint]struct{}

func NewEndpointSet(es ...*Endpoint) EndpointSet {
	s := make(EndpointSet, len(es))
	for _, e := range es {
		s.Add(e)
	}
	return s
}

func (s EndpointSet) Has(e *Endpoint) bool {
	if s == nil {
		return false
	}
	_, has := s[e]
	return has
}

// AnyOne returns one of the endpoints in the set, or nil if the set is empty.
//
// For a set with more than one element the result is un-defined, so this
// method should generally be used only with one-element sets.
func (s EndpointSet) AnyOne() *Endpoint {
	for endpoint := range s {
		return endpoint
	}
	return nil
}

// Names returns the names for all of the endpoints in the set, sorted
// lexicographically.
//
// This is primarily a test-assertion utility.
func (s EndpointSet) Names() []string {
	var ret []string
	for endpoint := range s {
		ret = append(ret, endpoint.Name)
	}
	sort.Strings(ret)
	return ret
}

func (s EndpointSet) Add(e *Endpoint) {
	s[e] = struct{}{}
}

func (s EndpointSet) Remove(e *Endpoint) {
	if s == nil {
		return
	}
	delete(s, e)
}

func (s EndpointSet) List() []*Endpoint {
	if s == nil {
		return nil
	}

	ret := make([]*Endpoint, 0, len(s))
	for e := range s {
		ret = append(ret, e)
	}
	return ret
}

func (s EndpointSet) Union(o EndpointSet) EndpointSet {
	ret := make(EndpointSet, len(s)+len(o))
	for e := range s {
		ret[e] = struct{}{}
	}
	for e := range o {
		ret[e] = struct{}{}
	}
	return ret
}

func (s EndpointSet) Subtract(o EndpointSet) EndpointSet {
	ret := make(EndpointSet, len(s))
	for e := range s {
		if !o.Has(e) {
			ret[e] = struct{}{}
		}
	}
	return ret
}
