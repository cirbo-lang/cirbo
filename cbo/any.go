package cbo

// Any is a special value that represents a value that is one of a number of
// different types depending on what was specified in the source Cirbo
// program.
//
// "Any" values are produced by unwrapping cbty.Value values used in the
// language interpreter, producing a native Go type that is more convenient
// to use in calling code once interpretation is complete.
//
// "Any" is really just a different name for the empty interface, but
// eval.Unwrapper guarantees that the result will have one of the following
// dynamic types, depending on the type of cbty.Value given.
//
//     any unknown value              nil
//     cbty.String                    string
//     cbty.Bool                      bool
//     any number or quantity type    units.Quantity
//     any generic object type        map[string]Any
//     any device type                *cbo.Device
//     any device instance type       *cbo.DeviceInstance
//     any circuit type               *cbo.Circuit
//     any circuit instance type      *cbo.CircuitInstance
//     any land type                  *cbo.Land
//     pinout                         *cbo.Pinout
//     board                          *cbo.Board
//     cbty.TypeType                  cbty.Type
//     any function type              func (args cbty.CallArgs) (Any, source.Diags)
//
// Callers making use of this type should use a type switch or type assertion
// to select for their desired type(s) and handle gracefully any type that is
// not expected, including types not included in the above list to allow for
// future expansion of the above list.
type Any interface {
}
