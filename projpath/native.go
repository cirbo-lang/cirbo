package projpath

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// nativeFS is the primary projectImpl, which operates against the real
// host filesystem.
type nativeFS struct {
	isProjectImpl

	cfg PathConfig
}

// PathConfig specifies the root locations that will be used to resolve
// user-provided and source-provided paths within this project.
type PathConfig struct {
	// WorkingDir is the directory that any relative filesystem paths will
	// be resolved relative to. Should be provided as an absolute filesystem
	// path.
	//
	// WorkingDir is primarily relevant for CLI use. If using Cirbo from
	// a context where "working directory" isn't a relevant concept, leave
	// this blank. If working in an IDE or text editor that has the concept
	// of a "project root", set WorkingDir to the project root so that
	// file paths can be given relative to the project root.
	WorkingDir string

	// SystemPkgDir is the directory where "system packages" (the built-in
	// packages that ship with Cirbo) are installed. Should be provided either
	// as an absolute filesystem path or a path relative to WorkingDir.
	SystemPkgDir string
}

// NewProject constructs and returns a new Project object using the
// given path configuration.
//
// The given working directory must either be blank or an absolute path,
// or else this function will panic. It will also panic if WorkingDir
// is blank but SystemPkgDir is a relative path.
func NewProject(config PathConfig) Project {
	// Validate and normalize the given paths
	if config.WorkingDir != "" {
		if !filepath.IsAbs(config.WorkingDir) {
			panic(fmt.Errorf("attempt to create project with relative WorkingDir %q", config.WorkingDir))
		}

		config.WorkingDir = filepath.Clean(config.WorkingDir)
		if !filepath.IsAbs(config.SystemPkgDir) {
			config.SystemPkgDir = filepath.Join(config.WorkingDir, config.SystemPkgDir)
		}
	} else {
		if !filepath.IsAbs(config.SystemPkgDir) {
			panic(fmt.Errorf("attempt to create project with relative SystemPkgDir %q and no WorkingDir", config.SystemPkgDir))
		}
	}

	return Project{nativeFS{
		cfg: config,
	}}
}

func (fs nativeFS) ReadFile(path FilePath) ([]byte, error) {
	return ioutil.ReadFile(string(path))
}

func (fs nativeFS) ListFiles(path FilePath) []FilePath {
	entries, err := ioutil.ReadDir(string(path))
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
		ret = append(ret, FilePath(filepath.Join(string(path), info.Name())))
	}
	return ret
}

func (fs nativeFS) Stat(path FilePath) (os.FileInfo, error) {
	return os.Stat(string(path))
}

func (fs nativeFS) FilePathFromUI(path string) FilePath {
	if fs.cfg.WorkingDir != "" {
		path = absoluteNativePath(fs.cfg.WorkingDir, path)
	} else {
		if !filepath.IsAbs(path) {
			return NoPath
		}
		path = filepath.Clean(path)
	}

	return FilePath(path)
}

func (fs nativeFS) FilePathForUI(path FilePath) string {
	if fs.cfg.WorkingDir == "" {
		// If we have no working directory then we just always show
		// absolute paths.
		return string(path)
	}
	relPath, err := filepath.Rel(fs.cfg.WorkingDir, string(path))
	if err != nil {
		// Fall back on absolute path, though in practice this should
		// never happen because both of our given paths are absolute.
		return string(path)
	}
	return relPath
}

func (fs nativeFS) FilePathForPackagePath(pp packagePath, from FilePath) FilePath {
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
		fromDir, _ = filepath.Split(string(from))
		fromDir = filepath.Clean(fromDir) // trim off trailing slash, unless we're at a root
	}

	fspp := filepath.FromSlash(string(pp))

	// A relative path is directly relative to the requesting file.
	if pp.IsRel() {
		return FilePath(filepath.Join(fromDir, fspp))
	}

	// For an absolute path, we expect to find a suitable directory either
	// under a "cirbo-pkg" dir within the "from" directory (or some parent
	// directorty) or within the system package directory, in that order of
	// preference. The former is where the auto-installer ("cirbo get") would
	// put a third-party package retrieved from the Internet.

	path := fromDir
	volume := filepath.VolumeName(path)
	traversal := 0
	for {
		pkgPath := filepath.Join(fromDir, "cirbo-pkg", fspp)
		info, err = os.Stat(pkgPath)
		if err == nil && info.IsDir() {
			return FilePath(pkgPath)
		}

		if path[len(path)-1] == filepath.Separator {
			// If our path ends with the separator then we've reached the
			// root directory, so we'll stop trying to traverse further.
			break
		}

		if volume != "" && path == volume {
			// For Windows, we should also recognize when we've reached
			// the top of the "volume" so that we don't try to traverse up
			// past the "share name" of a UNC path.
			// For example, given the path //foo/bar/baz we must stop when
			// we reach //foo/bar, because //foo is not really a directory.
			//
			// volume might also be a drive prefix like "c:", but we don't
			// end up here in that case because we'll reach "c:\" first and
			// thus break out in the previous conditional branch above.
			//
			// (We never get in here on Unix because volume is always ""
			// on Unix systems.)
			break
		}

		path = filepath.Join(path, "..")
		traversal++

		// As a safety measure against behaviors we don't expect on
		// new Go target platforms, we'll always bail out after 15 steps
		// up because we're probably looping infinitely if we end up doing that.
		// This effectively means that local package directories within a
		// repository can't go more than 15 levels of nesting deep, which should
		// be a fine assumption in practice since nesting should be shallow
		// in all reasonable repositories.
		if traversal == 15 {
			break
		}
	}

	{
		sysPath := filepath.Join(fs.cfg.SystemPkgDir, fspp)
		info, err = os.Stat(sysPath)
		if err == nil && info.IsDir() {
			return FilePath(sysPath)
		}
	}

	return NoPath
}

func absoluteNativePath(base, given string) string {
	if filepath.IsAbs(given) {
		return given
	} else {
		return filepath.Join(base, given)
	}
}
