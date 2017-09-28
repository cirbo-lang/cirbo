package erc

import (
	"fmt"
	"strings"

	"github.com/cirbo-lang/cirbo/cbo"
)

// Error represents an error encountered during an electrical rules check.
type Error interface {
	error
	errorSigil() isError
}

type isError struct {
}

func (e isError) errorSigil() isError {
	return e
}

// Errors represents a set of errors encountered during an electrical rules check.
type Errors []Error

func (es Errors) Error() string {
	if len(es) == 0 {
		// Should never actually be seen, since callers should not try to
		// treat a nil Errors as an error.
		return "successful electrical rules check"
	}

	if len(es) > 1 {
		return fmt.Sprintf("%s, and %d other errors", es[0].Error(), len(es)-1)
	}

	return es[0].Error()
}

// ErrorNoOutput is an error returned when a particular net has input endpoints
// but no output endpoints.
//
// This can be overridden with an ERC-only component that has a placeholder
// output endpoint.
type ErrorNoOutput struct {
	isError
	Inputs   cbo.EndpointSet
	Passives cbo.EndpointSet
}

// ErrorNoInput is an error returned when a particular net has outputs
// (possibly in conflict, signalled by a separate error) but no inputs.
//
// This can be overridden with an ERC-only component that has a placeholder
// input endpoint.
type ErrorNoInput struct {
	isError
	Outputs  cbo.EndpointSet
	Passives cbo.EndpointSet
}

// ErrorSignalAsPower is an error returned when a signal output is driving
// a power input.
//
// This can be overridden with an ERC-only component that has a signal input
// on one side and a power output on the other.
type ErrorSignalAsPower struct {
	isError
	Drivers cbo.EndpointSet
	Driving cbo.EndpointSet
}

// ErrorOutputConflict is an error returned when incompatible outputs are
// mixed on a net. For example, if both a PushPull and an Open collector/drain
// are present on the same net, or if two PushPull outputs are present.
//
// A special endpoint of direction MultiOutputSinkFlag can be connected to
// a net that otherwise contains only outputs, to override this flag. It is
// intended to be used with an ERC-only "device" that has a MultiOutputSinkFlag
// terminal on one side and a normal output on the other, thus indicating that
// several outputs should be considered as one output for the net on the
// "normal" side of the device.
type ErrorOutputConflict struct {
	isError
	Outputs cbo.EndpointSet
}

// ErrorUnconnected is an error returned when a net contains only one
// endpoint.
//
// A special endpoint of direction NoConnectFlag can be connected to a
// single-endpoint net to override this flag with no change to the final
// component network
type ErrorUnconnected struct {
	isError
	Endpoint *cbo.Endpoint
}

// ErrorNoConnectConnected is an error returned when a net has a no-connect
// flag but it also has two or more other endpoints that effectively conflict
// with the indication that the net is unconnected.
type ErrorNoConnectConnected struct {
	isError
	Endpoints cbo.EndpointSet
	Flags     cbo.EndpointSet
}

func (e ErrorNoOutput) Error() string {
	return fmt.Sprintf(
		"Input(s) %s are not driven by any output",
		strings.Join(e.Inputs.Names(), ", "),
	)
}

func (e ErrorNoInput) Error() string {
	return fmt.Sprintf(
		"Output(s) %s are not driving any input",
		strings.Join(e.Outputs.Names(), ", "),
	)
}

func (e ErrorSignalAsPower) Error() string {
	return fmt.Sprintf(
		"Signal output(s) %s driving power input(s) %s",
		strings.Join(e.Drivers.Names(), ", "),
		strings.Join(e.Driving.Names(), ", "),
	)
}

func (e ErrorOutputConflict) Error() string {
	return fmt.Sprintf(
		"Incompatible outputs %s are driving each other",
		strings.Join(e.Outputs.Names(), ", "),
	)
}

func (e ErrorUnconnected) Error() string {
	return fmt.Sprintf(
		"%s is not connected to anything",
		e.Endpoint.Name,
	)
}

func (e ErrorNoConnectConnected) Error() string {
	return fmt.Sprintf(
		"%s are flagged as no-connect but yet connected to each other",
		strings.Join(e.Endpoints.Names(), ", "),
	)
}
