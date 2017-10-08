package projpath

import (
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

type inmemFS struct {
	isProjectImpl

	vfs vfs.FileSystem
}

// MockProject creates and returns a project that uses the given map as
// its source of files and directories.
//
// A mock project uses a virtual filesystem based on forward-slash-separated
// paths, regardless of the host OS, so that tests don't have to make any
// special effort to be portable unless they are innately OS-specific.
//
// This function is provided as a convenience for unit testing.
func MockProject(files map[string]string) Project {
	return Project{inmemFS{
		vfs: mapfs.New(files),
	}}
}

func (fs inmemFS) ReadFile(p FilePath) ([]byte, error) {
	r, err := fs.vfs.Open(string(p))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func (fs inmemFS) ListFiles(p FilePath) []FilePath {
	entries, err := fs.vfs.ReadDir(string(p))
	if err != nil {
		return nil
	}
	if len(entries) == 0 {
		return nil
	}

	ret := make([]FilePath, 0, len(entries))
	for _, info := range entries {
		if info.IsDir() {
			continue
		}
		ret = append(ret, FilePath(path.Join(string(p), info.Name())))
	}
	return ret
}

func (fs inmemFS) Stat(path FilePath) (os.FileInfo, error) {
	return fs.vfs.Stat(string(path))
}

func (fs inmemFS) FilePathFromUI(p string) FilePath {
	p = path.Join("/", p) // resolve relative to root and clean
	p = p[1:]             // trim off leading slash
	return FilePath(p)
}

func (fs inmemFS) FilePathForUI(p FilePath) string {
	// re-add the leading slash to indicate we're relative to the vfs root
	return path.Join("/", string(p))
}

func (fs inmemFS) FilePathForPackagePath(pp packagePath, from FilePath) FilePath {
	var fromDir string
	info, err := os.Stat(string(from))
	if err != nil {
		// If the "from" object doesn't exist at all, then no relative
		// package path is possible.
		return NoPath
	}
	if info.IsDir() {
		fromDir = string(from)
	} else {
		fromDir, _ = path.Split(string(from))
		fromDir = path.Clean(fromDir) // trim off trailing slash, unless we're at a root
	}

	// A relative path is directly relative to the requesting file.
	if pp.IsRel() {
		return FilePath(path.Join(fromDir, string(pp)))
	}

	// An absolute path is expected to be resolved under "/cirbo-pkg".
	// Notice that this is intentionally simpler than the native filesystem
	// implementation where we try to find the "closest" cirbo-pkg dir
	// in parent directories, since we expect users of the mock filesystem
	// to be testing simple, contrived situations.
	sysPath := path.Join("cirbo-pkg", string(pp))
	info, err = fs.vfs.Stat(sysPath)
	if err == nil && info.IsDir() {
		return FilePath(sysPath)
	}

	return NoPath
}
