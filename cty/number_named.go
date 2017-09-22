package cty

import (
	"github.com/cirbo-lang/cirbo/units"
)

// The Number type supports any dimensionality, but certain dimensionalities
// have special names due to their commonly-recognized meaning. These are
// tabulated here for easier reference.

//go:generate go run generate/generate_number_named_acc.go

var namedNumberTypes = map[string]Type{
	// Dimensionless numbers are just called "Number"
	"Number": Quantity(units.Dimensionality{}),

	// All the base dimensions get names
	"Mass": Quantity(units.Dimensionality{
		Mass: 1,
	}),
	"Length": Quantity(units.Dimensionality{
		Length: 1,
	}),
	"Angle": Quantity(units.Dimensionality{
		Angle: 1,
	}),
	"Time": Quantity(units.Dimensionality{
		Time: 1,
	}),
	"Current": Quantity(units.Dimensionality{
		ElectricCurrent: 1,
	}),
	"LuminousIntensity": Quantity(units.Dimensionality{
		LuminousIntensity: 1,
	}),

	// Some derived dimensions get names too
	"Force": Quantity(units.Dimensionality{
		Mass:   1,
		Length: 1,
		Time:   -2,
	}),
	"Momentum": Quantity(units.Dimensionality{
		Mass:   1,
		Length: 1,
		Time:   -1,
	}),
	"Area": Quantity(units.Dimensionality{
		Length: 2,
	}),
	"Speed": Quantity(units.Dimensionality{
		Length: 1,
		Time:   -1,
	}),
	"AngularSpeed": Quantity(units.Dimensionality{
		Angle: 1,
		Time:  -1,
	}),
	"Frequency": Quantity(units.Dimensionality{
		Time: -1,
	}),
	"Voltage": Quantity(units.Dimensionality{
		Mass:            1,
		Length:          2,
		Time:            -3,
		ElectricCurrent: -1,
	}),
	"Resistance": Quantity(units.Dimensionality{
		Mass:            1,
		Length:          2,
		Time:            -3,
		ElectricCurrent: -2,
	}),
	"Power": Quantity(units.Dimensionality{
		Mass:   1,
		Length: 2,
		Time:   -3,
	}),
	"Capacitance": Quantity(units.Dimensionality{
		Mass:            -1,
		Length:          -2,
		Time:            4,
		ElectricCurrent: 2,
	}),
	"Inductance": Quantity(units.Dimensionality{
		Mass:            1,
		Length:          2,
		Time:            -2,
		ElectricCurrent: -2,
	}),
	"Illuminance": Quantity(units.Dimensionality{
		Length:            -2,
		LuminousIntensity: 1,
	}),
}

var numberTypeNames map[units.Dimensionality]string

func init() {
	numberTypeNames = make(map[units.Dimensionality]string)
	for name, ty := range namedNumberTypes {
		numberTypeNames[ty.impl.(numberImpl).dim] = name
	}
}

// QuantityDimensionalities returns a map from each of the named quantity
// types to the dimensionality of that type.
func QuantityDimensionalities() map[string]units.Dimensionality {
	ret := make(map[string]units.Dimensionality)
	for name, ty := range namedNumberTypes {
		ret[name] = ty.impl.(numberImpl).dim
	}
	return ret
}

// QuantityByName returns the quantity type of the given name, or NilType
// if the name is not recognized.
func QuantityByName(name string) Type {
	return namedNumberTypes[name]
}
