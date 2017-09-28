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

	if len(endpoints) == 1 {
		errs = append(errs, ErrorUnconnected{
			Endpoint: endpoints.AnyOne(),
		})
		// This error overrides all the others, since saying "no inputs" or
		// "no outputs" is redundant with saying "nothing connected at all".
		return errs
	}

	classes := classify(endpoints)

	if len(classes.NoConnectFlags) > 0 {
		if len(endpoints) > (len(classes.NoConnectFlags) + 1) {
			errs = append(errs, ErrorNoConnectConnected{
				Endpoints: endpoints.Subtract(classes.NoConnectFlags),
				Flags:     classes.NoConnectFlags,
			})
		}
		// This error overrides all the others, since it is explicitly
		// saying that we don't want to connect two things together.
		return errs
	}

	if len(classes.Outputs) > 0 && len(classes.Inputs) == 0 && len(classes.Bidis) == 0 && len(classes.MultiOutputFlags) == 0 {
		errs = append(errs, ErrorNoInput{
			Outputs: classes.Outputs,
		})
	}

	if len(classes.Inputs) > 0 && len(classes.Outputs) == 0 && len(classes.Bidis) == 0 {
		errs = append(errs, ErrorNoOutput{
			Inputs: classes.Inputs,
		})
	}

	if len(classes.SignalOutputs) > 0 && len(classes.PowerInputs) > 0 {
		errs = append(errs, ErrorSignalAsPower{
			Drivers: classes.SignalOutputs,
			Driving: classes.PowerInputs,
		})
	}

	if len(classes.Outputs) > 0 && (len(classes.MultiOutputFlags) == 0 || len(classes.Inputs) > 0) {
		outClasses := classifyOutputs(classes.Outputs)

		if len(outClasses.PushPull) > 1 {
			errs = append(errs, ErrorOutputConflict{
				Outputs: outClasses.PushPull,
			})
		} else if len(outClasses.PushPull) > 0 && (len(outClasses.OpenCollector) > 0 || len(outClasses.OpenEmitter) > 0) {
			errs = append(errs, ErrorOutputConflict{
				Outputs: outClasses.PushPull.Union(outClasses.OpenCollector.Union(outClasses.OpenEmitter)),
			})
		} else if len(outClasses.OpenCollector) > 0 && len(outClasses.OpenEmitter) > 0 {
			errs = append(errs, ErrorOutputConflict{
				Outputs: outClasses.OpenCollector.Union(outClasses.OpenEmitter),
			})
		}
	}

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

type outputClassifications struct {
	PushPull      cbo.EndpointSet
	Tristate      cbo.EndpointSet
	OpenCollector cbo.EndpointSet
	OpenEmitter   cbo.EndpointSet
}

func classifyOutputs(endpoints cbo.EndpointSet) outputClassifications {
	var ret outputClassifications
	ret.PushPull = make(cbo.EndpointSet)
	ret.Tristate = make(cbo.EndpointSet)
	ret.OpenCollector = make(cbo.EndpointSet)
	ret.OpenEmitter = make(cbo.EndpointSet)

	for e := range endpoints {
		if e.ERC.Dir != cbo.Output && e.ERC.Dir != cbo.Bidirectional {
			// We only care about endpoints that can be outputs
			continue
		}

		switch e.ERC.OutputType {
		case cbo.PushPull:
			ret.PushPull.Add(e)
		case cbo.Tristate:
			ret.Tristate.Add(e)
		case cbo.OpenCollector:
			ret.OpenCollector.Add(e)
		case cbo.OpenEmitter:
			ret.OpenEmitter.Add(e)
		}
	}

	return ret
}
