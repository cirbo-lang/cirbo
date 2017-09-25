// Package erc contains Cirbo's Electrical Rules Checker.
//
// Electrical rules check is a special sort of semantic check which verifies
// that a particular circuit net has a suitable combination of endpoints that
// comply with the electrical rules.
//
// By default Cirbo considers it a fatal error to make connections that violate
// the electrical rules. However, certain rules can be overridden where the
// circuit author is intentionally doing something unusual.
package erc
