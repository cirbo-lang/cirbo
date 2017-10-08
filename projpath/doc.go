// Package projpath contains utilities for working with paths within the
// Cirbo project filesystem structure.
//
// Internally the system uses absolute filesystem paths in the canonical
// form for the host operating system. When dealing with CLI input from the
// user or output to the user we re-interpret paths as relative to the given
// working directory. When dealing with paths appearing in source files, we
// use a simpler alternative path representation that is portable across
// host operating systems and relative to the referring file.
package projpath
