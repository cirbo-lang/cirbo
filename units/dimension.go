package units

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

type BaseDimension int

// The order of these (really, their value) is used to create a canonical
// ordering for stringification, etc.
const (
	Invalid BaseDimension = iota
	Mass
	Length
	Angle
	ElectricCurrent
	LuminousIntensity
	Time
)

var powerReplacer = strings.NewReplacer(
	"0", "⁰",
	"1", "¹",
	"2", "²",
	"3", "³",
	"4", "⁴",
	"5", "⁵",
	"6", "⁶",
	"7", "⁷",
	"8", "⁸",
	"9", "⁹",
	"-", "⁻",
)

func (d BaseDimension) Symbol() string {
	switch d {
	case Mass:
		return "M"
	case Length:
		return "L"
	case Angle:
		// this is non-standard, so there is no symbol
		return "angle"
	case Time:
		return "T"
	case ElectricCurrent:
		return "I"
	case LuminousIntensity:
		return "J"
	default:
		panic("can't call Symbol() on invalid BaseDimension")
	}
}

func (d BaseDimension) String() string {
	return "[" + d.Symbol() + "]"
}

// Dimensionality represents powers of each supported dimension.
//
// Each field of this struct represents a specific dimension, and its value
// represents the power that it is raised to.
//
// Dimensionality is represented as a struct so that it can be compared with
// the == operator and used as a map key.
type Dimensionality struct {
	Mass              int
	Length            int
	Angle             int
	Time              int
	ElectricCurrent   int
	LuminousIntensity int
}

// Multiply returns the product of the receiver and the other given
// Dimensionality.
func (d Dimensionality) Multiply(o Dimensionality) Dimensionality {
	var ret Dimensionality
	ret.Mass = d.Mass + o.Mass
	ret.Length = d.Length + o.Length
	ret.Angle = d.Angle + o.Angle
	ret.Time = d.Time + o.Time
	ret.ElectricCurrent = d.ElectricCurrent + o.ElectricCurrent
	ret.LuminousIntensity = d.LuminousIntensity + o.LuminousIntensity
	return ret
}

// Reciprocal returns the multiplicative inverse of the receiver.
func (d Dimensionality) Reciprocal() Dimensionality {
	var ret Dimensionality
	ret.Mass = -d.Mass
	ret.Length = -d.Length
	ret.Angle = -d.Angle
	ret.Time = -d.Time
	ret.ElectricCurrent = -d.ElectricCurrent
	ret.LuminousIntensity = -d.LuminousIntensity
	return ret
}

// ToPower returns a new dimensionality that is the receiver rased to
// the given integer power.
//
// In other words, it multiplies all of the powers already present in the
// receiver by the given power, and returns the result.
func (d Dimensionality) ToPower(power int) Dimensionality {
	return Dimensionality{
		Mass:              d.Mass * power,
		Length:            d.Length * power,
		Angle:             d.Angle * power,
		Time:              d.Time * power,
		ElectricCurrent:   d.ElectricCurrent * power,
		LuminousIntensity: d.LuminousIntensity * power,
	}
}

// String returns a compact string representation of a dimensionality,
// which is stable for a given dimensionality value.
//
// This method is not optimized, since it's primarily intended for debugging
// and error messages.
func (d Dimensionality) String() string {
	b := &bytes.Buffer{}

	e := d.dimEntries()

	for _, ei := range e {
		b.WriteString(ei.Dimension.String())
		if ei.Power != 1 {
			b.WriteString(powerReplacer.Replace(strconv.Itoa(ei.Power)))
		}
	}

	return strings.TrimSpace(b.String())
}

func (d Dimensionality) dimEntries() dimEntries {
	e := make(dimEntries, 0, 6)
	if d.Mass != 0 {
		e = append(e, dimEntry{Mass, d.Mass})
	}
	if d.Length != 0 {
		e = append(e, dimEntry{Length, d.Length})
	}
	if d.Angle != 0 {
		e = append(e, dimEntry{Angle, d.Angle})
	}
	if d.ElectricCurrent != 0 {
		e = append(e, dimEntry{ElectricCurrent, d.ElectricCurrent})
	}
	if d.LuminousIntensity != 0 {
		e = append(e, dimEntry{LuminousIntensity, d.LuminousIntensity})
	}
	if d.Time != 0 {
		e = append(e, dimEntry{Time, d.Time})
	}
	sort.Stable(e)
	return e
}

type dimEntry struct {
	Dimension BaseDimension
	Power     int
}

type dimEntries []dimEntry

func (e dimEntries) Len() int {
	return len(e)
}

func (e dimEntries) Less(i, j int) bool {
	if e[i].Power != e[j].Power {
		// Sort unit powers before all others
		if e[i].Power == 1 {
			return true
		}
		if e[j].Power == 1 {
			return false
		}

		// Sort negative powers after positive powers
		if e[i].Power < 0 && e[j].Power >= 0 {
			return false
		}
		if e[j].Power < 0 && e[i].Power >= 0 {
			return true
		}

		api := e[i].Power
		apj := e[j].Power
		if api < 0 {
			api = -api
		}
		if apj < 0 {
			apj = -apj
		}
		return api < apj
	}

	return int(e[i].Dimension) < int(e[j].Dimension)
}

func (e dimEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
