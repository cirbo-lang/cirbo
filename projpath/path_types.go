package projpath

import (
	"strings"
	"unicode"
)

// FilePath is a specialization of string that represents an absolute path
// to a file or directory within the host filesystem.
//
// FilePaths may be constructed only by functions in this package, and must be
// considered to be opaque values by all calling code. They must not be exposed
// to end-users via the CLI; instead, they must be re-interpreted into
// working-directory-relative files using functions in this package.
type FilePath string

// NoPath is a special FilePath value representing the absense of a path.
const NoPath = FilePath("")

func (path FilePath) IsModule() bool {
	if !strings.HasSuffix(string(path), ".cbm") {
		return false
	}

	if strings.HasPrefix(string(path), ".") || strings.HasPrefix(string(path), "_") {
		return false
	}

	return true
}

// packagePath is an internal specialization of string that represents
// a package path. This is a utility that provides the OS-agnostic
// package path handling functionality that is common to all project
// implementations.
type packagePath string

// IsAbs returns true if the package path is an absolute one. "Absolute"
// here means that it is a name within the global package namespace and thus
// installable as an external dependency. "Relative" paths, on the other hand,
// have no global context and so cannot be installed an external dependency.
//
// The result is meaningful only if the path is valid, as defined by method
// Valid.
func (pp packagePath) IsAbs() bool {
	firstSlash := strings.IndexRune(string(pp), '/')
	var firstPart string
	switch {
	case firstSlash == -1:
		firstPart = string(pp)
	default:
		firstPart = string(pp)[:firstSlash]
	}

	return firstPart != "." && firstPart != ".."
}

// IsRel is the opposite of IsAbs
func (pp packagePath) IsRel() bool {
	return !pp.IsAbs()
}

func (pp packagePath) Valid() bool {
	// A valid package path is a forward-slash-separated path that contains
	// no "." or ".." components except at the start.
	// It's considered "bad form" to use .. to exit the repository containing
	// a particular file through relative traversal, but this is not enforced
	// by this function.
	first := true
	seenNonParent := false
	last := false

	s := string(pp)
	for !last {
		var cur string
		nextSlash := strings.IndexRune(s, '/')
		if nextSlash != -1 {
			cur, s = s[:nextSlash], s[nextSlash+1:]
		} else {
			cur = s
			last = true
		}

		// No empty components are allowed
		if len(cur) == 0 {
			return false
		}

		if cur == "." {
			// A "." segment is only allowed as the first segment, and only
			// if it's followed by at least one other segment.
			if last || !first {
				return false
			}
			seenNonParent = true
		} else if cur == ".." {
			// Any number of .. can appear only at the start of the path,
			// and only if followed by at least one non-.. segment.
			if seenNonParent || last {
				return false
			}
		} else {
			seenNonParent = true
		}

		for _, r := range cur {
			if last {
				// We are stricter about the final segment in a path, since
				// it should ideally be a valid identifier in the language.
				// (This is not actually exactly aligned with how the parser
				// thinks of identifiers, but that's the spirit and the
				// recommended usage.)
				switch {
				case unicode.IsLower(r): // okay
				case unicode.IsNumber(r): // okay
				case r == '_': // okay
				default: // everything else is not okay
					return false
				}
			} else {
				// The rest of the path is more liberal to allow for mapping
				// onto URLs for automatic resolution/installation.
				// However, we still require lowercase to ensure consistent
				// behavior on both case-sensitive and case-insensitive
				// filesystems.
				switch {
				case unicode.IsLower(r): // okay
				case unicode.IsNumber(r): // okay
				case r == '.' || r == '-' || r == '_': // okay
				default: // everything else is not okay
					return false
				}
			}
		}

		first = false
	}

	return true
}
