package cbo

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

func (s EndpointSet) Has(e *Endpoint) bool {
	_, has := s[e]
	return has
}

func (s EndpointSet) Add(e *Endpoint) {
	s[e] = struct{}{}
}

func (s EndpointSet) Remove(e *Endpoint) {
	delete(s, e)
}

func (s EndpointSet) List() []*Endpoint {
	ret := make([]*Endpoint, 0, len(s))
	for e := range s {
		ret = append(ret, e)
	}
	return ret
}
