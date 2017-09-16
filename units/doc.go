// Package units deals with the system of units used by Cirbo.
//
// It is not intended as a general-purpose units library; it's tailored for
// the needs of the Cirbo language. In particular, it considers angle to be
// a base quantity (contrary to SI standards) and treats decimal degrees as
// its primary unit (consistent with engineering practice). It also supports
// only integer powers, and thus cannot represent units such as
// "square root of seconds".
package units
