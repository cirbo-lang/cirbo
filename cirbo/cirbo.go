package cirbo

import (
	"github.com/cirbo-lang/cirbo/parser"
	"github.com/cirbo-lang/cirbo/projpath"
)

// Cirbo is the main entry-point object. It encapsulates state from other
// Cirbo subsystems to present a convenient interface to perform common
// actions.
type Cirbo struct {
	proj projpath.Project
	pars *parser.Parser
	pkgs pkgCache
}

type Config struct {
	// WorkingDir is an absolute path that will be used as the basis to
	// resolve relative paths passed to methods of Cirbo. A relative path
	// for WorkingDir is not permitted.
	//
	// WorkingDir may be empty when the caller is running in a context where
	// a working directory is not a meaningful concept, such as a GUI tool
	// or a web application. In that case, all paths passed to other methods
	// must themselves be absolute.
	//
	// If Cirbo is being instantiated in the context of an integration with
	// an IDE or text editor that has the notion of a "project" then
	// WorkingDir should usually be the root directory of the project.
	WorkingDir string

	// SystemPkgDir is the path where the system packages (those included in
	// the Cirbo distribution) can be found.
	//
	// If SystemPkgDir is a relative path then it is interpreted relative to
	// WorkingDir, if set. If WorkingDir is not set, SystemPkgDir must be
	// absolute.
	SystemPkgDir string
}

// New creates a new instance of Cirbo and returns it. If the given
// configuration is invalid (per the rules in the struct's own documentation)
// then this method will panic.
//
// It will not panic if the configuration is syntactically valid but target
// paths do not exist; that condition will instead result in failures during
// later compilation requests.
func New(config Config) *Cirbo {
	return &Cirbo{
		proj: projpath.NewProject(projpath.PathConfig{
			WorkingDir:   config.WorkingDir,
			SystemPkgDir: config.SystemPkgDir,
		}),
		pars: parser.NewParser(),
		pkgs: pkgCache{},
	}
}
