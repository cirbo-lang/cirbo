package cbo

import (
	"github.com/cirbo-lang/cirbo/cbty"
)

// AttributesDef represents a set of attribute definitions.
type AttributesDef struct {
	All map[string]AttributeDef
	Pos []string
}

// AttributeDef represents a single attribute definition.
type AttributeDef struct {
	Type     cbty.Type
	Required bool
}

func (ad AttributesDef) PosDefs() []AttributeDef {
	ret := make([]AttributeDef, len(ad.Pos))
	for i, n := range ad.Pos {
		ret[i] = ad.All[n]
	}
	return ret
}
