package source

import (
	"fmt"
)

// Diag represents a single diagnostic message (warning or error).
type Diag struct {
	Level   DiagLevel
	Summary string
	Detail  string
	Ranges  []Range
}

type DiagLevel rune

const (
	Error   DiagLevel = 'E'
	Warning DiagLevel = 'W'
)

type Diags []Diag

func (ds Diags) HasErrors() bool {
	for _, d := range ds {
		if d.Level == Error {
			return true
		}
	}
	return false
}

func (d Diag) String() string {
	if len(d.Ranges) > 0 {
		return fmt.Sprintf("%s: %s; %s", d.Ranges[0], d.Summary, d.Detail)
	} else {
		return fmt.Sprintf("%s; %s", d.Summary, d.Detail)
	}
}
