package units

import (
	"fmt"
	"math/big"
	"testing"
)

func TestQuantityConvert(t *testing.T) {
	tests := []struct {
		Q    Quantity
		U    *Unit
		Want string
	}{
		{
			q("1", unitByName["lb"]),
			unitByName["kg"],
			"0.45359237 kg",
		},
		{
			q("1", unitByName["lb"]),
			unitByName["st"],
			"14 st",
		},
		{
			q("1", unitByName["m"]),
			unitByName["cm"],
			"100 cm",
		},
		{
			q("1", unitByName["in"]),
			unitByName["cm"],
			"2.54 cm",
		},
		{
			q("1", unitByName["in"]),
			unitByName["mil"],
			"1000 mil",
		},
		{
			q("1", unitByName["ft"]),
			unitByName["in"],
			"12 in",
		},
		{
			q("1", unitByName["yd"]),
			unitByName["in"],
			"36 in",
		},
		{
			q("1", &Unit{Dimensionality{Length: 2}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: 2}, baseUnits{Length: centimeter}, 0},
			"10000 cm²",
		},
		{
			q("1", &Unit{Dimensionality{Length: 3}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: 3}, baseUnits{Length: centimeter}, 0},
			"1000000 cm³",
		},
		{
			q("1", &Unit{Dimensionality{Length: -2}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: -2}, baseUnits{Length: centimeter}, 0},
			"0.0001 cm⁻²",
		},
		{
			q("1", &Unit{Dimensionality{Length: -3}, baseUnits{Length: meter}, 0}),
			&Unit{Dimensionality{Length: -3}, baseUnits{Length: centimeter}, 0},
			"1e-06 cm⁻³",
		},
		{
			q("1", unitByName["MHz"]),
			unitByName["Hz"],
			"1000000 Hz",
		},
		{
			q("1", unitByName["ohm"]),
			unitByName["kohm"],
			"0.001 kohm",
		},
		{
			q("1", unitByName["kohm"]),
			unitByName["ohm"],
			"1000 ohm",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s to %s", test.Q, test.U), func(t *testing.T) {
			got := test.Q.Convert(test.U)
			gotStr := got.String()
			if gotStr != test.Want {
				t.Errorf("wrong result\nquant: %s\nunit:  %s\ngot:   %s\nwant:  %s", test.Q, test.U, gotStr, test.Want)
			}
		})
	}
}

func TestQuantityWithStandardUnits(t *testing.T) {
	tests := []struct {
		Q    Quantity
		Want string
	}{
		{
			q("20", dimless),
			"20",
		},
		{
			q("2", unitByName["lb"]),
			"0.90718474 kg",
		},
		{
			q("50", unitByName["cm"]),
			"0.5 m",
		},
		{
			q("5", unitByName["in"]),
			"0.127 m",
		},
		{
			q("2", unitByName["yd"]),
			"1.8288 m",
		},
		{
			q("1000", unitByName["mil"]),
			"0.0254 m",
		},
		{
			q("2", unitByName["turn"]),
			"720 deg",
		},
		{
			q("1", unitByName["turn"]),
			"360 deg",
		},
		{
			q("0.5", unitByName["turn"]),
			"180 deg",
		},
		{
			q("3.1415926535897932384626433832795028841971693993751", unitByName["rad"]),
			"180 deg",
		},
		{
			q("0.25", unitByName["turn"]),
			"90 deg",
		},
		{
			q("100", unitByName["ms"]),
			"0.1 s",
		},
		{
			q("1", unitByName["kohm"]),
			"1000 ohm",
		},
		{
			q("1000000", &Unit{
				Dimensionality{Length: 1, Time: -1},
				baseUnits{Length: inch, Time: microsecond},
				0,
			}),
			"0.0254 m s⁻¹",
		},
		{
			// This is a rather odd expression of voltage using inches,
			// which tests whether we end up normalizing the result to
			// be "V" after conversion. The un-normalized form is
			// kg m² s⁻³ A⁻¹, which is not the correct answer here.
			q("1000", &Unit{
				Dimensionality{Mass: 1, Length: 2, Time: -3, ElectricCurrent: -1},
				baseUnits{Mass: kilogram, Length: inch, Time: second, ElectricCurrent: ampere},
				0,
			}),
			"0.64516 V",
		},
	}

	for _, test := range tests {
		t.Run(test.Q.String(), func(t *testing.T) {
			got := test.Q.WithStandardUnits()
			gotStr := got.String()
			if gotStr != test.Want {
				t.Errorf("wrong result\ninput: %s\ngot:   %s\nwant:  %s", test.Q, gotStr, test.Want)
			}
		})
	}
}

func TestQuantityMultiply(t *testing.T) {
	tests := []struct {
		A    Quantity
		B    Quantity
		Want string
	}{
		{
			MakeDimensionless(bfp("2")),
			MakeDimensionless(bfp("2")),
			"4",
		},
		{
			MakeDimensionless(bfp("2")),
			MakeQuantity(bfp("2"), unitByName["m"]),
			"4 m",
		},
		{
			MakeQuantity(bfp("2"), unitByName["kg"]),
			MakeQuantity(bfp("2"), unitByName["m"]),
			"4 kg m",
		},
		{
			MakeQuantity(bfp("2"), unitByName["lb"]),
			MakeQuantity(bfp("2"), unitByName["in"]),
			"4 lb in",
		},
		{
			MakeQuantity(bfp("2"), unitByName["m"]),
			MakeQuantity(bfp("2"), unitByName["m"]),
			"4 m²",
		},
		{
			MakeQuantity(bfp("2"), unitByName["in"]),
			MakeQuantity(bfp("2"), unitByName["in"]),
			"4 in²",
		},
		{
			MakeQuantity(bfp("2"), unitByName["cm"]),
			MakeQuantity(bfp("2"), unitByName["in"]),
			"0.001016 m²",
		},
		{
			MakeQuantity(bfp("2"), unitByName["s"]).Multiply(MakeQuantity(bfp("2"), unitByName["s"])),
			MakeQuantity(bfp("2"), unitByName["s"]),
			"8 s³",
		},
		{
			MakeQuantity(bfp("2"), unitByName["kg"]),
			MakeQuantity(bfp("2"), unitByName["kg"]),
			"4 kg²",
		},
		{
			MakeQuantity(bfp("2"), unitByName["deg"]),
			MakeQuantity(bfp("2"), unitByName["deg"]),
			"4 deg²",
		},
		{
			MakeQuantity(bfp("2"), unitByName["A"]),
			MakeQuantity(bfp("2"), unitByName["A"]),
			"4 A²",
		},
		{
			MakeQuantity(bfp("2"), unitByName["cd"]),
			MakeQuantity(bfp("2"), unitByName["cd"]),
			"4 cd²",
		},
		{
			MakeQuantity(bfp("2"), unitByName["A"]),
			MakeQuantity(bfp("2"), unitByName["ohm"]),
			"4 V", // should normalize to V; kg m² s⁻³ A⁻¹ is equivalent but not what we want here
		},
		{
			MakeQuantity(bfp("2"), unitByName["A"]),
			MakeQuantity(bfp("2"), unitByName["kohm"]),
			"4000 V", // should normalize to V; kg m² s⁻³ A⁻¹ is equivalent but not what we want here
		},
		{
			MakeQuantity(bfp("2"), unitByName["mA"]),
			MakeQuantity(bfp("2"), unitByName["kohm"]),
			"4 V", // should normalize to V; kg m² s⁻³ A⁻¹ is equivalent but not what we want here
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s * %s", test.A, test.B), func(t *testing.T) {
			got := test.A.Multiply(test.B)
			gotStr := got.String()
			if gotStr != test.Want {
				t.Errorf("wrong result\ninput: %s * %s\ngot:   %s\nwant:  %s", test.A, test.B, gotStr, test.Want)
			}
		})
	}
}

func TestQuantityDivide(t *testing.T) {
	tests := []struct {
		A    Quantity
		B    Quantity
		Want string
	}{
		{
			MakeDimensionless(bfp("2")),
			MakeDimensionless(bfp("2")),
			"1",
		},
		{
			MakeDimensionless(bfp("2")),
			MakeQuantity(bfp("2"), unitByName["m"]),
			"1 m⁻¹",
		},
		{
			MakeQuantity(bfp("2"), unitByName["kg"]),
			MakeQuantity(bfp("2"), unitByName["m"]),
			"1 kg m⁻¹",
		},
		{
			MakeQuantity(bfp("2"), unitByName["lb"]),
			MakeQuantity(bfp("2"), unitByName["in"]),
			"1 lb in⁻¹",
		},
		{
			MakeQuantity(bfp("2"), unitByName["m"]),
			MakeQuantity(bfp("2"), unitByName["m"]),
			"1",
		},
		{
			MakeQuantity(bfp("2"), unitByName["in"]),
			MakeQuantity(bfp("2"), unitByName["in"]),
			"1",
		},
		{
			MakeQuantity(bfp("2"), unitByName["cm"]),
			MakeQuantity(bfp("2"), unitByName["in"]),
			"0.3937007874",
		},
		{
			MakeQuantity(bfp("2"), unitByName["s"]).Multiply(MakeQuantity(bfp("2"), unitByName["s"])),
			MakeQuantity(bfp("2"), unitByName["s"]),
			"2 s",
		},
		{
			MakeQuantity(bfp("2"), unitByName["kg"]),
			MakeQuantity(bfp("2"), unitByName["kg"]),
			"1",
		},
		{
			MakeQuantity(bfp("2"), unitByName["deg"]),
			MakeQuantity(bfp("2"), unitByName["deg"]),
			"1",
		},
		{
			MakeQuantity(bfp("2"), unitByName["A"]),
			MakeQuantity(bfp("2"), unitByName["A"]),
			"1",
		},
		{
			MakeQuantity(bfp("2"), unitByName["cd"]),
			MakeQuantity(bfp("2"), unitByName["cd"]),
			"1",
		},
		{
			MakeQuantity(bfp("2"), unitByName["A"]),
			MakeQuantity(bfp("2"), unitByName["ohm"]),
			"1 kg⁻¹ m⁻² s³ A³",
		},
		{
			MakeQuantity(bfp("2"), unitByName["A"]),
			MakeQuantity(bfp("2"), unitByName["kohm"]),
			"0.001 kg⁻¹ m⁻² s³ A³",
		},
		{
			MakeQuantity(bfp("2"), unitByName["mA"]),
			MakeQuantity(bfp("2"), unitByName["kohm"]),
			"1e-06 kg⁻¹ m⁻² s³ A³",
		},
		{
			MakeQuantity(bfp("4"), unitByName["m"]),
			MakeQuantity(bfp("2"), unitByName["s"]),
			"2 m s⁻¹",
		},
		{
			MakeQuantity(bfp("4"), unitByName["V"]),
			MakeQuantity(bfp("2"), unitByName["A"]),
			"2 ohm", // unit must normalize to ohm
		},
		{
			MakeQuantity(bfp("4"), unitByName["W"]),
			MakeQuantity(bfp("2"), unitByName["A"]),
			"2 V", // unit must normalize to volt
		},
		{
			MakeQuantity(bfp("4"), unitByName["W"]),
			MakeQuantity(bfp("2"), unitByName["V"]),
			"2 A",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s ÷ %s", test.A, test.B), func(t *testing.T) {
			got := test.A.Divide(test.B)
			gotStr := got.String()
			if gotStr != test.Want {
				t.Errorf("wrong result\ninput: %s ÷ %s\ngot:   %s\nwant:  %s", test.A, test.B, gotStr, test.Want)
			}
		})
	}
}

func TestQuantityString(t *testing.T) {
	tests := []struct {
		Input Quantity
		Want  string
	}{
		{
			q("1", unitByName["cm"]),
			"1 cm",
		},
		{
			q("1", nil),
			"1",
		},
		{
			q("2.89", &Unit{
				Dimensionality{Length: 1, Time: -2},
				baseUnits{Length: meter, Time: second},
				0,
			}),
			"2.89 m s⁻²",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Input), func(t *testing.T) {
			got := test.Input.String()
			if got != test.Want {
				t.Errorf("wrong result\ninput: %#v\ngot:   %s\nwant:  %s", test.Input, got, test.Want)
			}
		})
	}
}

func q(v string, u *Unit) Quantity {
	f, _, err := (&big.Float{}).Parse(v, 10)
	if err != nil {
		panic(err)
	}

	return MakeQuantity(f, u)
}
