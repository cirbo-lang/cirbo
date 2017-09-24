package cbo

// An Endpoint is a low-level object representing a participant in a net.
type Endpoint struct {
	Name       string
	Net        *Net
	Type       TerminalType
	Dir        TerminalDir
	Role       TerminalRole
	OutputType TerminalOutputType
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
