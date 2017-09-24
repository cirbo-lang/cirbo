package cbo

type Net struct {
	Endpoints EndpointSet

	onReplace []func(new *Net)
}

// Connect adds the given endpoint to the receiving net.
//
// If the endpoint is already a member of a net, that net merged with the
// receiver, causing the receiver to have the superset of both endpoint sets
// and the old net to have no endpoints at all. The old net should, at that
// point, be discarded entirely.
func (n *Net) Connect(e *Endpoint) {
	// If the given endpoint already has a net then we need to merge the two
	// nets together, collecting any other endpoints already associated with
	// the other net. The other net will be empty after we are done.
	if e.Net != nil {
		mn := e.Net
		otherEs := mn.Endpoints.List()
		for _, otherE := range otherEs {
			n.Endpoints.Add(otherE)
			otherE.Net = n
			mn.Endpoints.Remove(otherE)
		}

		for _, cb := range mn.onReplace {
			// Notify about the new net
			cb(n)

			// Incorporate this callback into the new net so that it will
			// be called again if the net is successively replaced by another.
			n.onReplace = append(n.onReplace, cb)
		}
		return
	}

	n.Endpoints.Add(e)
	e.Net = n
}

// SuggestedName attempts to suggest a name for the receiver based on the
// names of its member endpoints.
//
// It is not always possible to produce a good result, and the result is
// not guaranteed unique across a whole design. This is a "best effort" that
// hopefully gives good results in some common cases.
//
// If a name cannot be suggested at all, the result is an empty string.
func (n *Net) SuggestedName() string {
	if len(n.Endpoints) == 0 {
		return ""
	}

	nameOccurs := map[string]int{}
	for e := range n.Endpoints {
		nameOccurs[e.Name]++
	}

	for _, n := range priorityNetNames {
		if nameOccurs[n] > 0 {
			return n
		}
	}

	return ""
}

// OnReplace arranges for the given callback to be called if the receiving
// net is replaced due to being merged with another net.
//
// This allows a data structure that contains a net pointer to keep its pointer
// updated if its referent is merged with another net.
//
// The given callback may be called multiple times if the original net is
// progressively merged with further nets.
func (n *Net) OnReplace(cb func(new *Net)) {
	n.onReplace = append(n.onReplace, cb)
}

var priorityNetNames = []string{
	"GND",
	"AGND",
	"PGND",

	// TODO: ideally we would recognize _all_ voltage-level-looking names
	// here, but for now we'll just pick out some common ones.
	// We list these before the generic IC supply connector names below so
	// that in a multi-voltage circuit we'll prefer the more specific rail
	// name over the generic chip pin name.
	"+3V3",
	"-3V3",
	"+1V8",
	"-1V8",
	"+5V",
	"-5V",
	"+12V",
	"-12V",
	"+9V",
	"-9V",
	"+6V",
	"-6V",
	"+24V",
	"-24V",

	"V+",
	"VCC",
	"VDD",
	"V-",
	"VSS",
	"VEE",

	"MOSI",
	"MISO",
	"SCLK",
	"SCK",
	"SDA",
	"RX",
	"TX",
	"DTR",
	"DCD",
	"DSR",
	"RTS",
	"RTR",
	"CTS",
}
