package cirbo

import (
	"fmt"
	"sort"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/compiler"
	"github.com/cirbo-lang/cirbo/eval"
	"github.com/cirbo-lang/cirbo/projpath"
	"github.com/cirbo-lang/cirbo/source"
)

// LoadPackage loads a Cirbo package from the given filesystem path, loads
// all of the packages it depends on (directly or indirectly) and then
// returns the value exported from the package.
//
// Note that the diagnostics for each package are only returned the first
// time that package is evaluated, even if evaluated indirectly. This is done
// under the assumption that chains of calls to LoadPackage on the same
// Cirbo will accumulate any diagnostics into a single diagnostic list, and
// thus re-returning the same diagnostic would cause it to be duplicated in the
// output.
func (cb *Cirbo) LoadPackage(dir string) (cbty.Value, source.Diags) {
	fp := cb.proj.FilePathFromUI(dir)
	if fp == projpath.NoPath {
		return cbty.PlaceholderVal, source.Diags{
			source.Diag{
				Level:   source.Error,
				Summary: "Invalid package directory",
				Detail:  fmt.Sprintf("The path %q could not be resolved as a package directory.", dir),
			},
		}
	}

	if entry, ok := cb.pkgs.GetOk(fp); ok {
		return entry.Value, nil
	}

	var diags source.Diags

	dependents := map[projpath.FilePath][]projpath.FilePath{}
	inDeg := map[projpath.FilePath]int{}
	var queue []projpath.FilePath
	queue = append(queue, fp)
	pkgs := map[projpath.FilePath]*eval.Package{}
	pkgRefMap := map[projpath.FilePath]map[string]projpath.FilePath{}

	// Compile the requested package and all of its transitive dependencies
	for i := 0; i < len(queue); i++ { // list will grow during the loop
		pkgDir := queue[i]
		if cb.pkgs.Has(pkgDir) {
			// Dependency is already cached, so we don't need to visit it
			// or any of its dependencies.
			continue
		}

		if _, compiled := pkgs[pkgDir]; compiled {
			// Dependency isn't already cached, but we already planned to
			// compile it so we don't need to re-compile it here.
			continue
		}

		pkg, pkgDiags := cb.compilePackage(pkgDir)
		diags = append(diags, pkgDiags...)
		if pkgDiags.HasErrors() {
			// If we can't compile a particular package at all, we'll
			// just stub it out with a placeholder and skip over it.
			cb.pkgs.Put(pkgDir, pkgCacheEntry{
				Value: cbty.PlaceholderVal,
			})
			continue
		}

		pkgs[pkgDir] = pkg
		pkgRefMap[pkgDir] = map[string]projpath.FilePath{}

		for _, dep := range pkg.PackagesImported() {
			depDir := cb.proj.FilePathForPackagePath(dep.Path, dep.Range.Filename)
			if depDir == projpath.NoPath {
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Unresolvable import",
					// TODO: Once we have a package installer, reference it in
					// this error message as a "next step" for the user.
					Detail: fmt.Sprintf("The path %q could not be resolved to an available package.", dep.Path),
					Ranges: dep.Range.List(),
				})
				continue
			}
			pkgRefMap[pkgDir][dep.Path] = depDir
			inDeg[pkgDir]++
			dependents[depDir] = append(dependents[depDir], pkgDir)
			queue = append(queue, depDir)
		}
	}

	// If any of the packages failed to compile then we can't proceed.
	if diags.HasErrors() {
		cb.pkgs.Put(fp, pkgCacheEntry{
			Value: cbty.PlaceholderVal,
		})
		return cbty.PlaceholderVal, diags
	}

	// If we _did_ manage to load all of the necessary packages, then our
	// next job is to evaluate them in topoligical order to expand
	// our package value cache.

	queue = queue[:0] // reusing the queue backing array for the traversal queue
	for depDir := range pkgRefMap {
		if inDeg[depDir] == 0 {
			queue = append(queue, depDir)
		}
	}
	for len(queue) > 0 {
		var pkgDir projpath.FilePath
		pkgDir, queue = queue[0], queue[1:] // dequeue next item
		pkg := pkgs[pkgDir]

		depVals := map[string]cbty.Value{}

		for depName, depPath := range pkgRefMap[pkgDir] {
			if depEntry, ok := cb.pkgs.GetOk(depPath); ok {
				depVals[depName] = depEntry.Value
			} else {
				// Should never happen unless there's a bug in our topological
				// traversal code here.
				diags = append(diags, source.Diag{
					Level:   source.Error,
					Summary: "Unresolved package",
					Detail:  fmt.Sprintf("Required package %q for %q was not resolved in time. This is a bug in Cirbo that should be reported!", depPath, pkgDir),
				})
				depVals[depName] = cbty.PlaceholderVal
			}
		}

		result, evalDiags := pkg.ExportedValue(depVals)
		diags = append(diags, evalDiags...)
		cb.pkgs.Put(pkgDir, pkgCacheEntry{
			Value: result,
		})

		for _, enabledDir := range dependents[pkgDir] {
			inDeg[enabledDir]--
			if inDeg[enabledDir] < 1 {
				queue = append(queue, enabledDir)
				delete(inDeg, enabledDir)
			}
		}
	}

	if len(inDeg) > 0 {
		switch len(inDeg) {
		case 1:
			var selfDep projpath.FilePath
			for dfp := range inDeg {
				selfDep = dfp
			}
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Package dependency cycle",
				Detail: fmt.Sprintf(
					"The package at %s depends on itself. Dependency cycles are not allowed.",
					cb.proj.FilePathForUI(selfDep),
				),
			})
		case 2:
			paths := make([]string, 0, 2)
			for dfp := range inDeg {
				paths = append(paths, cb.proj.FilePathForUI(dfp))
			}
			sort.Strings(paths)
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Package dependency cycle",
				Detail: fmt.Sprintf(
					"The packages in %s and %s both depend on each other. Dependency cycles are not allowed.",
					paths[0], paths[1],
				),
			})
		default:
			// If we have a large number of items leftover then we've got
			// a messy cycle situation that would require some analysis
			// to get a great error message. For now, we'll just settle for
			// a non-great error message. We will improve on this if this
			// error ends up appearing a lot in practice.
			paths := make([]string, 0, len(inDeg))
			for dfp := range inDeg {
				paths = append(paths, cb.proj.FilePathForUI(dfp))
			}
			sort.Strings(paths)
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Package dependency cycle",
				Detail: fmt.Sprintf(
					"The package in %s either depends on itself or on packages that form a dependency cycle. Dependency cycles are not allowed. If you recently added an \"import\" statement, investigate whether that statement created a cycle.",
					paths[0],
				),
			})
		}
	}

	// If we've encountered any errors along the way then we probably won't
	// have a reasonable value to return, so we'll stub it out.
	if diags.HasErrors() {
		cb.pkgs.Put(fp, pkgCacheEntry{
			Value: cbty.PlaceholderVal,
		})
		return cbty.PlaceholderVal, diags
	}

	result := cb.pkgs.Get(fp).Value
	return result, diags
}

func (cb *Cirbo) compilePackage(pkgDir projpath.FilePath) (*eval.Package, source.Diags) {
	var diags source.Diags
	modFilenames := cb.proj.ListModuleFiles(pkgDir)

	files := make([]*ast.File, 0, len(modFilenames))
	for _, filename := range modFilenames {
		src, err := cb.proj.ReadFile(filename)
		if err != nil {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Failed to read module file",
				Detail: fmt.Sprintf(
					"The module file %s (for package in %s) could not be read: %s",
					cb.proj.FilePathForUI(filename),
					cb.proj.FilePathForUI(pkgDir),
					err.Error(),
				),
			})
			continue
		}

		file, fileDiags := cb.pars.ParseFile(filename, src)
		diags = append(diags, fileDiags...)
		files = append(files, file)
	}

	if diags.HasErrors() {
		// We'll just make the package empty, so we can proceed with
		// compiling up to a point.
		files = nil
	}

	pkgNode := ast.Package(files)

	pkg, compileDiags := compiler.CompilePackage(pkgNode)
	diags = append(diags, compileDiags...)
	return pkg, diags
}
