package projpath

import (
	"os"
	"sort"
)

type Project struct {
	i projectImpl
}

type projectImpl interface {
	projectImplSigil() isProjectImpl

	FilePathFromUI(string) FilePath
	FilePathForUI(FilePath) string
	FilePathForPackagePath(packagePath, FilePath) FilePath
	ReadFile(FilePath) ([]byte, error)
	ListFiles(FilePath) []FilePath
	Stat(FilePath) (os.FileInfo, error)
}

// Embed isProjectImpl into a struct to mark it as being a project
// implementation, along with implementing the other projectImpl methods.
type isProjectImpl struct {
}

func (i isProjectImpl) projectImplSigil() isProjectImpl {
	return i
}

// FilePathFromUI interprets a file path string given in the CLI (or
// equivalent) into a canonical FilePath, or returns NoPath if the given
// path is invalid.
func (p Project) FilePathFromUI(path string) FilePath {
	return p.i.FilePathFromUI(path)
}

// FilePathForDisplay converts an canonical, opaque FilePath into a string
// suitable to show to an end-user to describe the given path. The result
// takes into account any working directory context to show a relative
// path where possible.
func (p Project) FilePathForUI(path FilePath) string {
	return p.i.FilePathForUI(path)
}

// ReadFile reads the contents of a file identified by an opaque FilePath.
// Use FilePathFromUI to convert a user-specified path to a FilePath.
func (p Project) ReadFile(path FilePath) ([]byte, error) {
	return p.i.ReadFile(path)
}

// ListFiles returns a list of the files in the directory identified by an
// opaque FilePath.
//
/// Use FilePathFromUI to convert a user-specified path to a FilePath.
func (p Project) ListFiles(path FilePath) []FilePath {
	ret := p.i.ListFiles(path)

	// Sort the files lexicographically just so the result is in some
	// reasonable predictable order. (This is also guaranteed by the
	// ListModuleFiles wrapper so that module files are parsed and
	// evaluated in a predictable order.)
	sort.Slice(ret, func(i, j int) bool {
		return string(ret[i]) < string(ret[j])
	})

	return ret
}

// ListModuleFiles is a convenience wrapper around ListFiles that filters
// the resulting list to include only Cirbo module files.
//
// The result is sorted into lexicographical order so that module files used as
// package contents can be processed in a predictable order.
func (p Project) ListModuleFiles(path FilePath) []FilePath {
	ret := p.ListFiles(path)

	// We're doing this filter in-place here, by copying entries earlier
	// in the list in order to close gaps left by items we don't want.
	// This causes the whole underlying slice to remain in memory, but
	// that's okay because we assume that package dirs will only contain
	// a small number of non-module files.
	l := 0
	for _, fpath := range ret {
		if fpath.IsModule() {
			// this is a no-op as long as l == i, but they will diverge if we
			// have non-module files in the list.
			ret[l] = fpath
			l++
		}
	}

	// Empty out any trailing slice elements so that we don't leak the strings
	// that they refer to.
	for i := l; i < len(ret); i++ {
		ret[i] = ""
	}

	// Truncate off any trailing items.
	return ret[:l]
}

// FilePathForPackagePath finds the FilePath of the directory containing
// the source files for the given package path when requested from the
// given FilePath.
//
// A package paths is always relative to some other path, which is usually
// the source file that imports it. Package paths are a simplified, OS-agnostic
// path representation which is mapped on to the real filesystem by this
// method.
//
// The result is NoPath if the given package path is invalid or if no suitable
// mapping to the underlying filesystem can be found.
func (p Project) FilePathForPackagePath(ppath string, from FilePath) FilePath {
	pp := packagePath(ppath)
	if !pp.Valid() {
		return NoPath
	}

	return p.i.FilePathForPackagePath(pp, from)
}
