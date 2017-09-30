package eval

import (
	"github.com/cirbo-lang/cirbo/cty"
)

type Expr interface {
	value() cty.Value
}
