package cbty

import (
	"testing"
)

func TestTypeImplInterfaces(t *testing.T) {
	// All of the following are actually compile-time assertions, but
	// we wrap this in a test function so we'll see an item for this
	// in verbose test output.

	// Embeddable type helpers
	var _ typeWithAttributes = staticAttributes(nil)

	// Specific type implementations
	var _ typeImpl = numberImpl{}
	var _ typeWithArithmetic = numberImpl{}
}
