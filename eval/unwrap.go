package eval

import (
	"github.com/cirbo-lang/cirbo/cbty"
)

// Unwrapped is the result type of Unwrap, which extracts a value of a native
// Go type out of a cbty.Value for more convenient use in calling code.
//
// Unwrapped is really just a different name for the empty interface, but
// Unwrap guarantees that the result will have one of the following dynamic
// types, depending on the type of cbty.Value given.
//
//     any unknown value              nil
//     cbty.String                    string
//     cbty.Bool                      bool
//     any number or quantity type    units.Quantity
//     any generic object type        map[string]Unwrapped (recursive unwrap)
//     any device type                *cbo.Device
//     any device instance type       *cbo.DeviceInstance
//     any circuit type               *cbo.Circuit
//     any circuit instance type      *cbo.CircuitInstance
//     any land type                  *cbo.Land
//     pinout                         *cbo.Pinout
//     board                          *cbo.Board
//     cbty.TypeType                  cbty.Type
//     any function type              func (args cbty.CallArgs) (Unwrapped, source.Diags)
//
// Callers making use of this type should use a type switch or type assertion
// to select for their desired type(s) and handle gracefully any type that is
// not expected, including types not included in the above list to allow for
// future expansion of the above list.
type Unwrapped interface {
}

// Unwrap obtains a native Go value corresponding to the given value within
// Cirbo's type system. See the documentation of type Unwrapped for details
// on what can be returned by this function.
func Unwrap(val cbty.Value) Unwrapped {
	if val == cbty.NilValue || val.IsUnknown() {
		return nil
	}

	ty := val.Type()

	// Simple types first
	switch ty {
	case cbty.String:
		return val.AsString()
	case cbty.Bool:
		return val.True()
	case cbty.TypeType:
		return val.UnwrapType()
	}

	// FIXME: should have a mapping for every possible cbty type
	return nil
}
