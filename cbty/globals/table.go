package globals

import (
	"github.com/cirbo-lang/cirbo/cbty"
)

// Table returns a table of symbols that should be included in the global
// scope.
//
// Many of these symbols are also defined as symbols within this Go package
// so that they can be conveniently accessed by integration code.
func Table() map[string]cbty.Value {
	return map[string]cbty.Value{
		"Angle":             Angle,
		"AngularSpeed":      AngularSpeed,
		"Area":              Area,
		"Bool":              Bool,
		"Capacitance":       Capacitance,
		"Charge":            Charge,
		"Conductance":       Conductance,
		"Conductivity":      Conductivity,
		"Current":           Current,
		"Force":             Force,
		"Frequency":         Frequency,
		"Illuminance":       Illuminance,
		"Inductance":        Inductance,
		"Length":            Length,
		"LuminousIntensity": LuminousIntensity,
		"Mass":              Mass,
		"Momentum":          Momentum,
		"Number":            Number,
		"Object":            Object,
		"Power":             Power,
		"Resistance":        Resistance,
		"Resistivity":       Resistivity,
		"Speed":             Speed,
		"String":            String,
		"Time":              Time,
		"Type":              Type,
		"Voltage":           Voltage,
	}
}
