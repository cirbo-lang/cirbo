package erc

import (
	"testing"

	"github.com/cirbo-lang/cirbo/cbo"
)

func TestFlattenPassthrough(t *testing.T) {
	nets := map[string]*cbo.Net{
		"empty":     testNet(), // stays empty
		"direct":    testNet(), // has only direct (non-passthrough) endpoints
		"passthruA": testNet(), // has passthrough to the one below
		"passthruB": testNet(), // has passthrough back to the one above
		"loopA":     testNet(), // loop A (connects to loop B and loop C)
		"loopB":     testNet(), // loop B (connects to loop A and loop C)
		"loopC":     testNet(), // loop C (connects to loop B and loop A)
	}
	netNames := map[*cbo.Net]string{}
	for k, net := range nets {
		netNames[net] = k
	}

	endpoints := map[string]*cbo.Endpoint{
		"direct1":        &cbo.Endpoint{},
		"direct2":        &cbo.Endpoint{},
		"passthru1":      &cbo.Endpoint{},
		"passthru2":      &cbo.Endpoint{},
		"passthrudirect": &cbo.Endpoint{},
		"loopAB":         &cbo.Endpoint{},
		"loopBA":         &cbo.Endpoint{},
		"loopBC":         &cbo.Endpoint{},
		"loopCB":         &cbo.Endpoint{},
		"loopCA":         &cbo.Endpoint{},
		"loopAC":         &cbo.Endpoint{},
		"loopdirect":     &cbo.Endpoint{},
	}
	endpointNames := map[*cbo.Endpoint]string{}
	for k, ep := range endpoints {
		endpointNames[ep] = k
	}

	// Set up the necessary passthrough connections
	endpoints["passthru1"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["passthru2"]})
	endpoints["passthru2"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["passthru1"]})
	endpoints["loopAB"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["loopBA"]})
	endpoints["loopBA"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["loopAB"]})
	endpoints["loopBC"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["loopCB"]})
	endpoints["loopCB"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["loopBC"]})
	endpoints["loopCA"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["loopAC"]})
	endpoints["loopAC"].Passthrough = testEndpointSet([]*cbo.Endpoint{endpoints["loopCA"]})

	// Put the endpoints in their relevant nets.
	nets["direct"].Connect(endpoints["direct1"])
	nets["direct"].Connect(endpoints["direct2"])
	nets["passthruA"].Connect(endpoints["passthru1"])
	nets["passthruB"].Connect(endpoints["passthru2"])
	nets["passthruB"].Connect(endpoints["passthrudirect"])
	nets["loopA"].Connect(endpoints["loopAB"])
	nets["loopA"].Connect(endpoints["loopAC"])
	nets["loopB"].Connect(endpoints["loopBA"])
	nets["loopB"].Connect(endpoints["loopBC"])
	nets["loopC"].Connect(endpoints["loopCA"])
	nets["loopC"].Connect(endpoints["loopCB"])
	nets["loopB"].Connect(endpoints["loopdirect"])

	// With all that setup done, we should now be able to do our flattening and
	// get back the merged sets of endpoints.
	allNets := make([]*cbo.Net, 0, len(nets))
	for _, net := range nets {
		allNets = append(allNets, net)
	}
	got := FlattenPassthrough(allNets)
	want := map[*cbo.Net]cbo.EndpointSet{
		nets["empty"]: cbo.EndpointSet{},
		nets["direct"]: cbo.EndpointSet{
			endpoints["direct1"]: struct{}{},
			endpoints["direct2"]: struct{}{},
		},
		nets["passthruA"]: cbo.EndpointSet{
			endpoints["passthrudirect"]: struct{}{},
		},
		nets["passthruB"]: cbo.EndpointSet{
			endpoints["passthrudirect"]: struct{}{},
		},
		nets["loopA"]: cbo.EndpointSet{
			endpoints["loopdirect"]: struct{}{},
		},
		nets["loopB"]: cbo.EndpointSet{
			endpoints["loopdirect"]: struct{}{},
		},
		nets["loopC"]: cbo.EndpointSet{
			endpoints["loopdirect"]: struct{}{},
		},
	}

	if gotLen, wantLen := len(got), len(allNets); gotLen != wantLen {
		t.Fatalf("wrong number of result nets %d; want %d", gotLen, wantLen)
	}

	for _, net := range allNets {
		gotSet := got[net]
		wantSet := want[net]
		if gotSet == nil {
			t.Errorf("result missing net %q", netNames[net])
			continue
		}

		for wantEndpoint := range wantSet {
			if !gotSet.Has(wantEndpoint) {
				t.Errorf("result for net %q missing endpoint %q", netNames[net], endpointNames[wantEndpoint])
			}
		}
		for gotEndpoint := range gotSet {
			if !wantSet.Has(gotEndpoint) {
				t.Errorf("result for net %q extraneous endpoint %q", netNames[net], endpointNames[gotEndpoint])
			}
		}
	}
}

func testNet() *cbo.Net {
	return &cbo.Net{
		Endpoints: make(cbo.EndpointSet),
	}
}

func testEndpointSet(items []*cbo.Endpoint) cbo.EndpointSet {
	ret := make(cbo.EndpointSet, len(items))
	for _, item := range items {
		ret.Add(item)
	}
	return ret
}
