package cbty

import (
	"fmt"
	"testing"

	"github.com/cirbo-lang/cirbo/units"
)

func TestNumberTypeName(t *testing.T) {
	tests := []struct {
		Type Type
		Want string
	}{
		{Angle, "Angle"},
		{AngularSpeed, "AngularSpeed"},
		{Area, "Area"},
		{Capacitance, "Capacitance"},
		{Conductance, "Conductance"},
		{Conductivity, "Conductivity"},
		{Current, "Current"},
		{Force, "Force"},
		{Frequency, "Frequency"},
		{Illuminance, "Illuminance"},
		{Inductance, "Inductance"},
		{Length, "Length"},
		{LuminousIntensity, "LuminousIntensity"},
		{Mass, "Mass"},
		{Momentum, "Momentum"},
		{Number, "Number"},
		{Power, "Power"},
		{Resistance, "Resistance"},
		{Resistivity, "Resistivity"},
		{Speed, "Speed"},
		{Time, "Time"},
		{Voltage, "Voltage"},
		{Quantity(units.Dimensionality{Length: -1}), "Quantity([L]⁻¹)"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Type), func(t *testing.T) {
			got := test.Type.Name()
			if got != test.Want {
				t.Errorf("wrong result\ngot:  %s\nwant: %s", got, test.Want)
			}
		})
	}
}
