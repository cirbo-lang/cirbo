package units

// The following are some units exposed as variables for convenience of
// code in other packages that needs to work directly with specific units.
// This is intentionally not a comprehensive catalog of all of the named
// units (they can all be accessed via units.ByName) but these ones are
// commonly-used-enough to warrant a shorthand representation.
//
// These should generally be units that would be expected to be used when
// generating application-specific file formats, for convenient use of
// Quantity.FormatValue, though other uses may be warranted too.

// Kilogram is equivalent to ByName("kg")
var Kilogram *Unit

// Kilogram is equivalent to ByName("m")
var Meter *Unit

// Kilogram is equivalent to ByName("mm")
var Millimeter *Unit

// Kilogram is equivalent to ByName("in")
var Inch *Unit

// Kilogram is equivalent to ByName("mil")
var Mil *Unit

// Kilogram is equivalent to ByName("deg")
var Degree *Unit

// Kilogram is equivalent to ByName("s")
var Second *Unit

// Kilogram is equivalent to ByName("A")
var Ampere *Unit

// Kilogram is equivalent to ByName("cd")
var Candela *Unit

// DegreeTenths is a strange unit provided only for use with Quantity.FormatValue
// when generating angles for Kicad, which represents angles in tenths of a degree.
//
// This unit should not be used for general quantity computation, since it does
// not obey the usual rule that only named units can be scaled versions of other
// units and it may thus produce incorrect results under operations other than
// conversion.
var DegreeTenths *Unit

func init() {
	Kilogram = unitByName["kg"]
	Meter = unitByName["m"]
	Millimeter = unitByName["mm"]
	Inch = unitByName["in"]
	Mil = unitByName["mil"]
	Degree = unitByName["deg"]
	DegreeTenths = &Unit{Dimensionality{Angle: 1}, baseUnits{Angle: degree}, -10}
	Second = unitByName["s"]
	Ampere = unitByName["A"]
	Candela = unitByName["cd"]
}
