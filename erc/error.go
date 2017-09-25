package erc

import (
	"fmt"

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

// ErrorNoDriver is an error returned when a particular net has no output
// endpoints connected.
type ErrorNoDriver struct {
	Driving []*cbo.Endpoint
}

// ErrorSignalAsPower is an error returned when a signal output is driving
// a power input.
type ErrorSignalAsPower struct {
	Driver  *cbo.Endpoint
	Driving []*cbo.Endpoint
}

// ErrorOutputConflict is an error returned when incompatible outputs are
// mixed on a net. For example, if both a PushPull and an Open collector/drain
// are present on the same net, or if two PushPull outputs are present.
type ErrorOutputConflict struct {
	Drivers []*cbo.Endpoint
}
