package cbo

import (
	"github.com/cirbo-lang/cirbo/cbty"
)

// AttributesDef represents a set of attribute definitions.
type AttributesDef map[string]AttributeDef

// AttributeDef represents a single attribute definition.
type AttributeDef struct {
	Type     cbty.Type
	Required bool
}
