package erc

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/cbo"
)

func TestClassify(t *testing.T) {
	s := make(cbo.EndpointSet)

	s.Add(&cbo.Endpoint{
		Name: "VCC",
		ERC: cbo.ERCMode{
			Type: cbo.Power,
			Dir:  cbo.Output,
		},
	})
	s.Add(&cbo.Endpoint{
		Name: "GND",
		ERC: cbo.ERCMode{
			Type: cbo.Power,
			Dir:  cbo.Output,
		},
	})
	s.Add(&cbo.Endpoint{
		Name: "SCLK",
		ERC: cbo.ERCMode{
			Type: cbo.Signal,
			Dir:  cbo.Input,
		},
	})
	s.Add(&cbo.Endpoint{
		Name: "MISO",
		ERC: cbo.ERCMode{
			Type:       cbo.Signal,
			Dir:        cbo.Output,
			OutputType: cbo.PushPull,
		},
	})
	s.Add(&cbo.Endpoint{
		Name: "SDA",
		ERC: cbo.ERCMode{
			Type:       cbo.Signal,
			Dir:        cbo.Bidirectional,
			OutputType: cbo.OpenCollector,
		},
	})
	s.Add(&cbo.Endpoint{
		Name: "AVCC",
		ERC: cbo.ERCMode{
			Type: cbo.Power,
			Dir:  cbo.Input,
		},
	})
	s.Add(&cbo.Endpoint{
		Name: "NC",
		ERC: cbo.ERCMode{
			Type: cbo.Passive,
			Dir:  cbo.NoConnectFlag,
		},
	})
	s.Add(&cbo.Endpoint{
		Name: "MOF",
		ERC: cbo.ERCMode{
			Type: cbo.Passive,
			Dir:  cbo.MultiOutputSinkFlag,
		},
	})

	classes := classify(s)

	if got, want := classes.Inputs.Names(), []string{"AVCC", "SCLK"}; !reflect.DeepEqual(got, want) {
		t.Errorf("wrong inputs\ngot:  %#v\nwant: %#v", got, want)
	}
	if got, want := classes.Outputs.Names(), []string{"GND", "MISO", "VCC"}; !reflect.DeepEqual(got, want) {
		t.Errorf("wrong outputs\ngot:  %#v\nwant: %#v", got, want)
	}
	if got, want := classes.Bidis.Names(), []string{"SDA"}; !reflect.DeepEqual(got, want) {
		t.Errorf("wrong bidis\ngot:  %#v\nwant: %#v", got, want)
	}
	if got, want := classes.SignalOutputs.Names(), []string{"MISO", "SDA"}; !reflect.DeepEqual(got, want) {
		t.Errorf("wrong signal outputs\ngot:  %#v\nwant: %#v", got, want)
	}
	if got, want := classes.PowerInputs.Names(), []string{"AVCC"}; !reflect.DeepEqual(got, want) {
		t.Errorf("wrong power inputs\ngot:  %#v\nwant: %#v", got, want)
	}
	if got, want := classes.NoConnectFlags.Names(), []string{"NC"}; !reflect.DeepEqual(got, want) {
		t.Errorf("wrong no-connect flags\ngot:  %#v\nwant: %#v", got, want)
	}
	if got, want := classes.MultiOutputFlags.Names(), []string{"MOF"}; !reflect.DeepEqual(got, want) {
		t.Errorf("wrong multi-output flags\ngot:  %#v\nwant: %#v", got, want)
	}
}
