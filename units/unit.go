package units

import (
	"bytes"
	"strconv"
	"strings"
)

type Unit struct {
	dim  Dimensionality
	base baseUnits

	// For scaled units (millis, kilos, megas, etc) this stores the scale
	// factor. If positive then a value must be divided by it to recover
	// the unscaled unit, while if it's negative a value must be multiplied
	// by its absolute value. 0 means "unscaled"
	//
	// Scaling is only used for derived units. Units of base dimensions are
	// just represented directly.
	scale int64
}

// ByName returns the unit with the given name, or nil if the name is not
// recognized as a unit.
//
// This function only works with units that have specifically been named. It
// cannot be used to create unnamed derived units, such as "meters squared";
// these must be constructed, such as:
//
//     meters := units.ByName("m")
//     metersSquared := meters.Multiply(meters)
func ByName(n string) *Unit {
	return unitByName[n]
}

var dimless = &Unit{Dimensionality{}, baseUnits{}, 0}

var unitByName map[string]*Unit = map[string]*Unit{
	// Dimensionless
	"": dimless,

	// Mass Units
	"kg": &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: kilogram}, 0},
	"g":  &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: gram}, 0},
	"lb": &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: pound}, 0},
	"st": &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: stone}, 0},

	// Length Units
	"m":   &Unit{Dimensionality{Length: 1}, baseUnits{Length: meter}, 0},
	"mm":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: millimeter}, 0},
	"cm":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: centimeter}, 0},
	"km":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: kilometer}, 0},
	"mil": &Unit{Dimensionality{Length: 1}, baseUnits{Length: mil}, 0},
	"in":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: inch}, 0},
	"ft":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: foot}, 0},
	"yd":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: yard}, 0},

	// Angle Units
	"deg":  &Unit{Dimensionality{Angle: 1}, baseUnits{Angle: degree}, 0},
	"rad":  &Unit{Dimensionality{Angle: 1}, baseUnits{Angle: radian}, 0},
	"turn": &Unit{Dimensionality{Angle: 1}, baseUnits{Angle: turn}, 0},

	// Time Units
	"s":  &Unit{Dimensionality{Time: 1}, baseUnits{Time: second}, 0}, // There is no secs in physics.
	"ms": &Unit{Dimensionality{Time: 1}, baseUnits{Time: millisecond}, 0},
	"us": &Unit{Dimensionality{Time: 1}, baseUnits{Time: microsecond}, 0},

	// Electric Current Units
	"A":  &Unit{Dimensionality{ElectricCurrent: 1}, baseUnits{ElectricCurrent: ampere}, 0},
	"mA": &Unit{Dimensionality{ElectricCurrent: 1}, baseUnits{ElectricCurrent: milliampere}, 0},

	// Luminous Intensity Units
	"cd": &Unit{Dimensionality{LuminousIntensity: 1}, baseUnits{LuminousIntensity: candela}, 0},

	// Electic Resistance Units
	"ohm": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		0,
	},
	"kohm": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		1000,
	},
	"Mohm": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		1000000,
	},

	// Electric Voltage Units
	"V": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -1},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		0,
	},
	"mV": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -1},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		-1000,
	},
	"kV": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -1},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		1000,
	},

	// Frequency Units
	"Hz": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
		0,
	},
	"kHz": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
		1000,
	},
	"MHz": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
		1000000,
	},

	// Force Units
	"N": &Unit{
		Dimensionality{Mass: 1, Length: 1, Time: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second},
		0,
	},

	// Power Units
	"W": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3},
		baseUnits{Mass: kilogram, Length: meter, Time: second},
		0,
	},
	"mW": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
		-1000,
	},
	"kW": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
		1000,
	},
	"MW": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
		1000000,
	},
	"GW": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
		1000000000,
	},

	// Electrical Capacitance Units
	"F": &Unit{
		Dimensionality{Mass: -1, Length: -2, Time: 4, ElectricCurrent: 2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		0,
	},
	"mF": &Unit{
		Dimensionality{Mass: -1, Length: -2, Time: 4, ElectricCurrent: 2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		-1000,
	},
	"uF": &Unit{
		Dimensionality{Mass: -1, Length: -2, Time: 4, ElectricCurrent: 2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		-1000000,
	},

	// Electrical Inductance Units
	"H": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -2, ElectricCurrent: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		0,
	},
	"uH": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -2, ElectricCurrent: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
		-1000000,
	},

	// Illuminance units
	"lx": &Unit{
		Dimensionality{Length: -2, LuminousIntensity: 1},
		baseUnits{Length: meter, LuminousIntensity: candela},
		0,
	},
}

var units map[Unit]*Unit
var unitName map[*Unit]string

func init() {
	units = make(map[Unit]*Unit, len(unitByName))
	unitName = make(map[*Unit]string, len(unitByName))
	for name, unit := range unitByName {
		units[*unit] = unit
		unitName[unit] = name
	}
}

// Dimensionality returns the dimensionality of the receiver.
func (u *Unit) Dimensionality() Dimensionality {
	return u.dim
}

// CommensurableWith returns true if the receiver and the given unit
// have the same dimensionality.
//
// For example, two length units are commensurable but a mass unit is not
// commensurable with a length unit.
func (u *Unit) CommensurableWith(other *Unit) bool {
	if u == nil && other == nil {
		return true
	}
	if u == nil || other == nil {
		return false
	}

	return u.dim == other.dim
}

func (u *Unit) ToStandardUnits() *Unit {
	nu := *u

	if u.base.Mass != nil {
		nu.base.Mass = kilogram
	}
	if u.base.Length != nil {
		nu.base.Length = meter
	}
	if u.base.Angle != nil {
		nu.base.Angle = degree
	}
	if u.base.Time != nil {
		nu.base.Time = second
	}
	if u.base.ElectricCurrent != nil {
		nu.base.ElectricCurrent = ampere
	}
	if u.base.LuminousIntensity != nil {
		nu.base.LuminousIntensity = candela
	}
	nu.scale = 0

	return (&nu).normalize()
}

// Multiply returns a new unit that is the result of multiplying the receiver
// with the given unit.
//
// Since we can't convert an associated value when multiplying units alone,
// both units must be using the same units for each base dimension.
// For example, it's not valid to multiply inches by centimeters because
// the result would be inch-centimeters, and that's a nonsense unit.
// This method will panic if given differing base units.
//
// Use SameBaseUnits to determine if two units are safe to multiply.
// Quantity.Multiply is an easier method to use, since it can handle value
// conversions automatically.
func (u *Unit) Multiply(o *Unit) *Unit {
	n := &Unit{}
	n.dim.Mass = u.dim.Mass + o.dim.Mass
	n.dim.Length = u.dim.Length + o.dim.Length
	n.dim.Angle = u.dim.Angle + o.dim.Angle
	n.dim.Time = u.dim.Time + o.dim.Time
	n.dim.ElectricCurrent = u.dim.ElectricCurrent + o.dim.ElectricCurrent
	n.dim.LuminousIntensity = u.dim.LuminousIntensity + o.dim.LuminousIntensity

	if !u.SameBaseUnits(o) {
		panic("can't multiply units with differing base units")
	}

	if u.base.Mass != nil {
		n.base.Mass = u.base.Mass
	}
	if u.base.Length != nil {
		n.base.Length = u.base.Length
	}
	if u.base.Angle != nil {
		n.base.Angle = u.base.Angle
	}
	if u.base.Time != nil {
		n.base.Time = u.base.Time
	}
	if u.base.ElectricCurrent != nil {
		n.base.ElectricCurrent = u.base.ElectricCurrent
	}
	if u.base.LuminousIntensity != nil {
		n.base.LuminousIntensity = u.base.LuminousIntensity
	}

	if o.base.Mass != nil {
		n.base.Mass = o.base.Mass
	}
	if o.base.Length != nil {
		n.base.Length = o.base.Length
	}
	if o.base.Angle != nil {
		n.base.Angle = o.base.Angle
	}
	if o.base.Time != nil {
		n.base.Time = o.base.Time
	}
	if o.base.ElectricCurrent != nil {
		n.base.ElectricCurrent = o.base.ElectricCurrent
	}
	if o.base.LuminousIntensity != nil {
		n.base.LuminousIntensity = o.base.LuminousIntensity
	}

	n.scale = u.scale

	return n.normalize()
}

// SameBaseUnits returns true if the two units have either the same base
// units for each dimension or have non-overlapping dimensionality.
//
// For derived units that have an associated SI scale factor, this too must
// match for the result to be true.
func (u *Unit) SameBaseUnits(o *Unit) bool {
	if u.base.Mass != o.base.Mass && u.base.Mass != nil && o.base.Mass != nil {
		return false
	}
	if u.base.Length != o.base.Length && u.base.Length != nil && o.base.Length != nil {
		return false
	}
	if u.base.Angle != o.base.Angle && u.base.Angle != nil && o.base.Angle != nil {
		return false
	}
	if u.base.Time != o.base.Time && u.base.Time != nil && o.base.Time != nil {
		return false
	}
	if u.base.ElectricCurrent != o.base.ElectricCurrent && u.base.ElectricCurrent != nil && o.base.ElectricCurrent != nil {
		return false
	}
	if u.base.LuminousIntensity != o.base.LuminousIntensity && u.base.LuminousIntensity != nil && o.base.LuminousIntensity != nil {
		return false
	}
	if u.scale != o.scale {
		return false
	}
	return true
}

// String returns a compact, human-readable string representation of a given
// unit.
//
// If the unit has its own name then that will be returned. Otherwise,
// a name will be constructed from the set of base units that the given
// derived unit is built from.
func (u *Unit) String() string {
	if name, hasName := unitName[u]; hasName {
		return name
	}

	// If the unit doesn't have a name then we'll construct one from its
	// base units, using its dimensionality entries.
	buf := &bytes.Buffer{}
	e := u.dim.dimEntries()

	if u.scale != 0 && u.scale != 1 {
		// should never happen, since scale should be used only for named units.
		panic("constructed unit may not have scale")
	}

	for _, ei := range e {
		var unit *Unit
		switch ei.Dimension {
		case Mass:
			unit = massUnits[u.base.Mass]
		case Length:
			unit = lengthUnits[u.base.Length]
		case Angle:
			unit = angleUnits[u.base.Angle]
		case Time:
			unit = timeUnits[u.base.Time]
		case ElectricCurrent:
			unit = electricCurrentUnits[u.base.ElectricCurrent]
		case LuminousIntensity:
			unit = luminousIntensityUnits[u.base.LuminousIntensity]
		default:
			// should never happen if dimEntries is working correctly
			panic("String called on Unit with invalid base dimension entry")
		}

		// We assume here that all units in the base tables will have names;
		// if any don't, that's a bug to be fixed.
		unitName := unitName[unit]
		buf.WriteString(unitName)
		if ei.Power != 1 {
			buf.WriteString(powerReplacer.Replace(strconv.Itoa(ei.Power)))
		}
		buf.WriteString(" ")
	}

	return strings.TrimSpace(buf.String())
}

// normalize checks if the receiver is one of the named units, and if so
// returns the pointer to the canonical instance of that unit, which can
// then in turn by looked up in the unitName table to recover the name.
func (u *Unit) normalize() *Unit {
	if nu := units[*u]; nu != nil {
		return nu
	}
	return u
}
