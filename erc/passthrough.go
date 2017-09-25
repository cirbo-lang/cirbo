package erc

import (
	"github.com/cirbo-lang/cirbo/cbo"
)

// FlattenPassthrough analyses the graph for "passthrough" edges and produces
// a map from each of the given nets to the set of endpoints that interact
// with that net, either directly or indirectly.
func FlattenPassthrough(nets []*cbo.Net) map[*cbo.Net]cbo.EndpointSet {
	// What we're doing here is analogous to a data-flow analysis of a
	// control flow graph in an imperative programming language, with
	// the transfer function being the union of all of a net's endpoints
	// and the passthrough endpoints for each endpoint.

	queue := newNetQueue(len(nets))
	ret := make(map[*cbo.Net]cbo.EndpointSet, len(nets))

	// Fill the queue with the nets we were given. The order here doesn't
	// really matter, though we'll get better performance if the list happens
	// to be ordered such that most or all of the passthrough results have
	// already been evaluated before we reach a particular net.
	for _, net := range nets {
		queue.Append(net)
	}

	// Keep iterating until our queue is empty. If we make any changes
	// to a net's set then we'll add its passthroughs to the queue, so we
	// are likely to visit each net multiple times but we should eventually
	// converge on a fixpoint. We only ever grow the sets, so we know that
	// _at worst_ we will add every endpoint to every set and then exit.
	for current := queue.Take(); current != nil; current = queue.Take() {
		set := ret[current]
		if set == nil {
			set = make(cbo.EndpointSet, len(current.Endpoints))
			ret[current] = set
		}

		// We need to track if we add anything to our set, so we can re-queue
		// our dependents in that case.
		changed := false

		for ep := range current.Endpoints {
			if ep.Passthrough != nil {
				// For endpoints with passthrough, we traverse their
				// internal connection and import the endpoints from the
				// nets on "the other side".
				for bep := range ep.Passthrough {
					if bep.Net == nil {
						continue
					}

					bset := ret[bep.Net]
					for dep := range bset {
						if dep.Passthrough != nil {
							continue
						}

						if !set.Has(dep) {
							set.Add(dep)
							changed = true
						}
					}
				}

			} else {
				// For endpoints _without_ passthrough, we just add them
				// directly to our set.

				if !set.Has(ep) {
					set.Add(ep)
					changed = true
				}
			}
		}

		if changed {
			// Make sure that dependent nets get re-processed to see
			// our updated set.
			for ep := range current.Endpoints {
				if ep.Passthrough == nil {
					continue
				}

				for bep := range ep.Passthrough {
					if bep.Net != nil {
						// This might be a no-op, if this net was already
						// in the queue anyway.
						queue.Append(bep.Net)
					}
				}
			}
		}
	}

	return ret
}
