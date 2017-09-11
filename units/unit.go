package units

type Unit struct {
	dim  Dimensionality
	base baseUnits
}

var unitByName map[string]*Unit = map[string]*Unit{
	// Mass Units
	"kg": &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: kilogram}},
	"g":  &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: gram}},
	"lb": &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: pound}},
	"st": &Unit{Dimensionality{Mass: 1}, baseUnits{Mass: stone}},

	// Length Units
	"m":   &Unit{Dimensionality{Length: 1}, baseUnits{Length: meter}},
	"mm":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: millimeter}},
	"cm":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: centimeter}},
	"km":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: kilometer}},
	"mil": &Unit{Dimensionality{Length: 1}, baseUnits{Length: mil}},
	"in":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: inch}},
	"ft":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: foot}},
	"yd":  &Unit{Dimensionality{Length: 1}, baseUnits{Length: yard}},

	// Angle Units
	"deg":  &Unit{Dimensionality{Angle: 1}, baseUnits{Angle: degree}},
	"rad":  &Unit{Dimensionality{Angle: 1}, baseUnits{Angle: radian}},
	"turn": &Unit{Dimensionality{Angle: 1}, baseUnits{Angle: turn}},

	// Time Units
	"s":  &Unit{Dimensionality{Time: 1}, baseUnits{Time: second}}, // There is no secs in physics.
	"ms": &Unit{Dimensionality{Time: 1}, baseUnits{Time: millisecond}},
	"us": &Unit{Dimensionality{Time: 1}, baseUnits{Time: microsecond}},

	// Electric Current Units
	"A": &Unit{Dimensionality{ElectricCurrent: 1}, baseUnits{ElectricCurrent: ampere}},

	// Luminous Intensity Units
	"cd": &Unit{Dimensionality{LuminousIntensity: 1}, baseUnits{LuminousIntensity: candela}},

	// Electic Resistance Units
	"ohm": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
	},
	// TODO: represent kohm, Mohm?

	// Electric Voltage Units
	"V": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -1},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
	},
	// TODO: represent mV, kV

	// Frequency Units
	"Hz": &Unit{
		Dimensionality{Time: -1},
		baseUnits{Time: second},
	},
	// TODO: represent kHz, mHz

	// Force Units
	"N": &Unit{
		Dimensionality{Mass: 1, Length: 1, Time: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second},
	},

	// Power Units
	"W": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -3},
		baseUnits{Mass: kilogram, Length: meter, Time: second},
	},
	// TODO: represent mW, kW

	// Electrical Capacitance Units
	"F": &Unit{
		Dimensionality{Mass: -1, Length: -2, Time: 4, ElectricCurrent: 2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
	},
	// TODO: represent uF, mF

	// Electrical Inductance Units
	"H": &Unit{
		Dimensionality{Mass: 1, Length: 2, Time: -2, ElectricCurrent: -2},
		baseUnits{Mass: kilogram, Length: meter, Time: second, ElectricCurrent: ampere},
	},
	// TODO: represent uH

	// Illuminance units
	"lx": &Unit{
		Dimensionality{Length: -2, LuminousIntensity: 1},
		baseUnits{Length: meter, LuminousIntensity: candela},
	},
}

var unitName map[Unit]string

func init() {
	unitName = make(map[Unit]string, len(unitByName))
	for name, unit := range unitByName {
		unitName[*unit] = name
	}
}
