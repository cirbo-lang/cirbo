package erc

import (
	"github.com/cirbo-lang/cirbo/cbo"
)

// CheckNet applies the ERC rules to the given set of endpoints, assuming that
// they all belong to the same net, and returns Error objects describing any
// deviations from the rules.
//
// The result will be nil if no inconsistencies are detected.
func CheckNet(endpoints cbo.EndpointSet) Errors {
	var errs Errors

	//classes := classify(endpoints)

	return errs
}

// CheckNets is a helper wrapper around CheckNet that applies checks to many
// nets in a single call.
//
// The map returned will be nil if no errors are encountered at all, and will
// have keys present only for nets that have errors.
func CheckNets(endpointSets map[*cbo.Net]cbo.EndpointSet) map[*cbo.Net]Errors {
	var ret map[*cbo.Net]Errors

	for net, endpoints := range endpointSets {
		errs := CheckNet(endpoints)
		if errs != nil {
			if ret == nil {
				ret = make(map[*cbo.Net]Errors)
			}
			ret[net] = errs
		}
	}

	return ret
}

type classifications struct {
	Inputs  cbo.EndpointSet
	Outputs cbo.EndpointSet
	Bidis   cbo.EndpointSet

	SignalOutputs cbo.EndpointSet
	PowerInputs   cbo.EndpointSet

	NoConnectFlags   cbo.EndpointSet
	MultiOutputFlags cbo.EndpointSet
}

func classify(endpoints cbo.EndpointSet) classifications {
	var ret classifications
	ret.Inputs = make(cbo.EndpointSet)
	ret.Outputs = make(cbo.EndpointSet)
	ret.Bidis = make(cbo.EndpointSet)
	ret.SignalOutputs = make(cbo.EndpointSet)
	ret.PowerInputs = make(cbo.EndpointSet)
	ret.NoConnectFlags = make(cbo.EndpointSet)
	ret.MultiOutputFlags = make(cbo.EndpointSet)

	for e := range endpoints {

		switch e.ERC.Dir {
		case cbo.Input:
			ret.Inputs.Add(e)
		case cbo.Output:
			ret.Outputs.Add(e)
		case cbo.Bidirectional:
			ret.Bidis.Add(e)
		case cbo.MultiOutputSinkFlag:
			ret.MultiOutputFlags.Add(e)
		case cbo.NoConnectFlag:
			ret.NoConnectFlags.Add(e)
		}

		switch e.ERC.Type {
		case cbo.Signal:
			if e.ERC.Dir == cbo.Output || e.ERC.Dir == cbo.Bidirectional {
				ret.SignalOutputs.Add(e)
			}
		case cbo.Power:
			if e.ERC.Dir == cbo.Input {
				ret.PowerInputs.Add(e)
			}
		}
	}

	return ret
}
