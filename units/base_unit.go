package units

import (
	"fmt"
	"math/big"
)

var one = bf("1")

type baseUnits struct {
	Mass              *massUnit
	Length            *lengthUnit
	Angle             *angleUnit
	Time              *timeUnit
	ElectricCurrent   *electricCurrentUnit
	LuminousIntensity *luminousIntensityUnit
}

type massUnit struct {
	Scale big.Float
}

type lengthUnit struct {
	Scale big.Float
}

type angleUnit struct {
	Scale big.Float
}

type timeUnit struct {
	Scale big.Float
}

type electricCurrentUnit struct {
	Scale big.Float
}

type luminousIntensityUnit struct {
	Scale big.Float
}

var kilogram = &massUnit{bf("1")}
var gram = &massUnit{bf("1000")}
var pound = &massUnit{bf("2.2046226218487758072297380134502703385420702733602")}
var stone = &massUnit{bf("30.864716705882861301216332188303784739588983827043")}

var meter = &lengthUnit{bf("1")}
var centimeter = &lengthUnit{bf("100")}
var millimeter = &lengthUnit{bf("1000")}
var kilometer = &lengthUnit{bf("0.001")}
var yard = &lengthUnit{bf("1.0936132983377077865266841644794400699912510936133")}
var inch = &lengthUnit{bf("39.370078740157480314960629921259842519685039370079")}
var foot = &lengthUnit{bf("3.2808398950131233595800524934383202099737532808399")}
var mil = &lengthUnit{bf("39370.078740157480314960629921259842519685039370079")}

var degree = &angleUnit{bf("1")}
var radian = &angleUnit{bf("0.017453292519943295769236907684886127134428718885417")}
var turn = &angleUnit{bf("0.002777777777777777777777777777777777777777777778")}

var second = &timeUnit{bf("1")}
var millisecond = &timeUnit{bf("1000")}
var microsecond = &timeUnit{bf("1000000")}

var ampere = &electricCurrentUnit{bf("1")}
var milliampere = &electricCurrentUnit{bf("1000")}

var candela = &luminousIntensityUnit{bf("1")}

var massUnits = map[*massUnit]*Unit{
	kilogram: unitByName["kg"],
	gram:     unitByName["g"],
	pound:    unitByName["lb"],
	stone:    unitByName["st"],
}
var lengthUnits = map[*lengthUnit]*Unit{
	meter:      unitByName["m"],
	centimeter: unitByName["cm"],
	millimeter: unitByName["mm"],
	kilometer:  unitByName["km"],
	yard:       unitByName["yd"],
	inch:       unitByName["in"],
	foot:       unitByName["ft"],
	mil:        unitByName["mil"],
}
var angleUnits = map[*angleUnit]*Unit{
	degree: unitByName["deg"],
	radian: unitByName["rad"],
	turn:   unitByName["turn"],
}
var timeUnits = map[*timeUnit]*Unit{
	second:      unitByName["s"],
	millisecond: unitByName["ms"],
	microsecond: unitByName["us"],
}
var electricCurrentUnits = map[*electricCurrentUnit]*Unit{
	ampere:      unitByName["A"],
	milliampere: unitByName["mA"],
}
var luminousIntensityUnits = map[*luminousIntensityUnit]*Unit{
	candela: unitByName["cd"],
}

func bfp(s string) *big.Float {
	ret := &big.Float{}
	_, _, err := ret.Parse(s, 10)
	if err != nil {
		panic(fmt.Errorf("failed to parse float %q", s))
	}
	return ret
}

func bf(s string) big.Float {
	var ret big.Float
	_, _, err := ret.Parse(s, 10)
	if err != nil {
		panic(fmt.Errorf("failed to parse float %q", s))
	}
	return ret
}

func bff(v float64) big.Float {
	f := big.NewFloat(v)
	return *f
}

func powerScale(scale *big.Float, power int) *big.Float {
	if power == 1 {
		return scale
	}

	ret := (&big.Float{}).Copy(scale)
	var absPower int
	if power < 0 {
		absPower = -power
	} else {
		absPower = power
	}

	if absPower > 1 {
		p := (&big.Float{}).Copy(ret)
		for i := 1; i < absPower; i++ {
			ret.Mul(ret, p)
		}
	}

	if power < 0 {
		ret.Quo(&one, ret)
	}

	return ret
}
