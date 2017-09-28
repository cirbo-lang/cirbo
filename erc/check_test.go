package erc

import (
	"reflect"
	"testing"

	"github.com/cirbo-lang/cirbo/cbo"
	"github.com/davecgh/go-spew/spew"
)

func TestCheckNet(t *testing.T) {
	type testCase struct {
		Endpoints cbo.EndpointSet
		Want      Errors
	}

	tests := map[string]func() testCase{
		"empty net": func() testCase {
			return testCase{
				make(cbo.EndpointSet),
				Errors(nil),
			}
		},
		"one output many inputs": func() testCase {
			s := make(cbo.EndpointSet)
			s.Add(&cbo.Endpoint{
				Name: "OUT",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			})
			s.Add(&cbo.Endpoint{
				Name: "IN1",
				ERC: cbo.ERCMode{
					Type: cbo.Signal,
					Dir:  cbo.Input,
				},
			})
			s.Add(&cbo.Endpoint{
				Name: "IN2",
				ERC: cbo.ERCMode{
					Type: cbo.Signal,
					Dir:  cbo.Input,
				},
			})
			return testCase{
				Endpoints: s,
				Want:      Errors(nil),
			}
		},
		"two push-pull outputs one input": func() testCase {
			s := make(cbo.EndpointSet)
			out1 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			out2 := &cbo.Endpoint{
				Name: "OUT2",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			in := &cbo.Endpoint{
				Name: "IN",
				ERC: cbo.ERCMode{
					Type: cbo.Signal,
					Dir:  cbo.Input,
				},
			}
			s.Add(out1)
			s.Add(out2)
			s.Add(in)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorOutputConflict{
						Outputs: cbo.NewEndpointSet(out1, out2),
					},
				},
			}
		},
		"no input, conflicting outputs": func() testCase {
			s := make(cbo.EndpointSet)
			out1 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			out2 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			s.Add(out1)
			s.Add(out2)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorNoInput{
						Outputs: cbo.NewEndpointSet(out1, out2),
					},
					ErrorOutputConflict{
						Outputs: cbo.NewEndpointSet(out1, out2),
					},
				},
			}
		},
		"no input, compatible outputs": func() testCase {
			s := make(cbo.EndpointSet)
			out1 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.Tristate,
				},
			}
			out2 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.Tristate,
				},
			}
			s.Add(out1)
			s.Add(out2)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorNoInput{
						Outputs: cbo.NewEndpointSet(out1, out2),
					},
				},
			}
		},
		"no output": func() testCase {
			s := make(cbo.EndpointSet)
			in1 := &cbo.Endpoint{
				Name: "IN1",
				ERC: cbo.ERCMode{
					Type: cbo.Signal,
					Dir:  cbo.Input,
				},
			}
			in2 := &cbo.Endpoint{
				Name: "IN2",
				ERC: cbo.ERCMode{
					Type: cbo.Signal,
					Dir:  cbo.Input,
				},
			}
			s.Add(in1)
			s.Add(in2)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorNoOutput{
						Inputs: cbo.NewEndpointSet(in1, in2),
					},
				},
			}
		},
		"multi-output flag": func() testCase {
			s := make(cbo.EndpointSet)
			out1 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			out2 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			flag := &cbo.Endpoint{
				Name: "FLAG",
				ERC: cbo.ERCMode{
					Type: cbo.Passive,
					Dir:  cbo.MultiOutputSinkFlag,
				},
			}
			s.Add(out1)
			s.Add(out2)
			s.Add(flag)
			return testCase{
				Endpoints: s,
				Want:      nil,
			}
		},
		"multi-output flag with input": func() testCase {
			s := make(cbo.EndpointSet)
			out1 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			out2 := &cbo.Endpoint{
				Name: "OUT1",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			in := &cbo.Endpoint{
				Name: "IN",
				ERC: cbo.ERCMode{
					Type: cbo.Signal,
					Dir:  cbo.Input,
				},
			}
			flag := &cbo.Endpoint{
				Name: "FLAG",
				ERC: cbo.ERCMode{
					Type: cbo.Passive,
					Dir:  cbo.MultiOutputSinkFlag,
				},
			}
			s.Add(out1)
			s.Add(out2)
			s.Add(in)
			s.Add(flag)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorOutputConflict{
						Outputs: cbo.NewEndpointSet(out1, out2),
					},
				},
			}
		},
		"unconnected": func() testCase {
			s := make(cbo.EndpointSet)
			out := &cbo.Endpoint{
				Name: "OUT",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			s.Add(out)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorUnconnected{
						Endpoint: out,
					},
				},
			}
		},
		"unconnected with flag": func() testCase {
			s := make(cbo.EndpointSet)
			out := &cbo.Endpoint{
				Name: "OUT",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			flag := &cbo.Endpoint{
				Name: "FLAG",
				ERC: cbo.ERCMode{
					Type: cbo.Passive,
					Dir:  cbo.NoConnectFlag,
				},
			}
			s.Add(out)
			s.Add(flag)
			return testCase{
				Endpoints: s,
				Want:      nil,
			}
		},
		"no-connect flag but connected": func() testCase {
			s := make(cbo.EndpointSet)
			out := &cbo.Endpoint{
				Name: "OUT",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			in := &cbo.Endpoint{
				Name: "IN",
				ERC: cbo.ERCMode{
					Type: cbo.Signal,
					Dir:  cbo.Input,
				},
			}
			flag := &cbo.Endpoint{
				Name: "FLAG",
				ERC: cbo.ERCMode{
					Type: cbo.Passive,
					Dir:  cbo.NoConnectFlag,
				},
			}
			s.Add(out)
			s.Add(in)
			s.Add(flag)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorNoConnectConnected{
						Endpoints: cbo.NewEndpointSet(in, out),
						Flags:     cbo.NewEndpointSet(flag),
					},
				},
			}
		},
		"signal output driving power input": func() testCase {
			s := make(cbo.EndpointSet)
			out := &cbo.Endpoint{
				Name: "OUT",
				ERC: cbo.ERCMode{
					Type:       cbo.Signal,
					Dir:        cbo.Output,
					OutputType: cbo.PushPull,
				},
			}
			in := &cbo.Endpoint{
				Name: "IN",
				ERC: cbo.ERCMode{
					Type: cbo.Power,
					Dir:  cbo.Input,
				},
			}
			s.Add(out)
			s.Add(in)
			return testCase{
				Endpoints: s,
				Want: Errors{
					ErrorSignalAsPower{
						Drivers: cbo.NewEndpointSet(out),
						Driving: cbo.NewEndpointSet(in),
					},
				},
			}
		},
	}

	spewer := spew.NewDefaultConfig()
	spewer.DisableMethods = true

	for name, cons := range tests {
		t.Run(name, func(t *testing.T) {
			test := cons()
			got := CheckNet(test.Endpoints)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("wrong result\ngot: %swant: %s", spewer.Sdump(got), spewer.Sdump(test.Want))
			}
		})
	}
}

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
