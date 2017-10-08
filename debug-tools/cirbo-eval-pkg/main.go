package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/cbo"
	"github.com/cirbo-lang/cirbo/compiler"
	"github.com/cirbo-lang/cirbo/cty"
	"github.com/cirbo-lang/cirbo/parser"
	"github.com/cirbo-lang/cirbo/projpath"
	"github.com/cirbo-lang/cirbo/source"
)

func main() {
	err := realMain(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", err)
		os.Exit(1)
	}
}

func realMain(args []string) error {
	fl := flag.NewFlagSet("cirbo-eval-pkg", flag.ExitOnError)
	err := fl.Parse(args)
	if err != nil {
		return err
	}
	args = fl.Args()

	if len(args) != 1 {
		fl.Usage()
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}
	proj := projpath.NewProject(projpath.PathConfig{
		WorkingDir:   wd,
		SystemPkgDir: wd, // not actually used here because we don't resolve imports
	})

	p := parser.NewParser()

	var diags source.Diags

	pkgDir := proj.FilePathFromUI(args[0])
	pkg, compileDiags := compilePackage(pkgDir, proj, p)
	diags = append(diags, compileDiags...)

	if diags.HasErrors() {
		return diagsError(diags)
	}

	packages := map[projpath.FilePath]*cbo.Package{}
	enables := map[projpath.FilePath][]projpath.FilePath{}
	inDeg := map[projpath.FilePath]int{}
	var queue []projpath.FilePath
	queue = append(queue, pkgDir)
	packages[pkgDir] = pkg

	// We're doing this "manually" instead of using range because we're going
	// to modify the slice as we go, potentially copying it.
	for i := 0; i < len(queue); i++ {
		pkgDir := queue[i]
		pkg := packages[pkgDir]
		for _, dep := range pkg.PackagesImported() {
			depDir := proj.FilePathForPackagePath(dep.Path, dep.Range.Filename)
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
			inDeg[pkgDir]++
			enables[depDir] = append(enables[depDir], pkgDir)
			if _, compiled := packages[depDir]; !compiled {
				depPkg, compileDiags := compilePackage(depDir, proj, p)
				diags = append(diags, compileDiags...)
				packages[depDir] = depPkg
				queue = append(queue, depDir)
			}
		}
	}

	if diags.HasErrors() {
		return diagsError(diags)
	}

	// We're now going to reuse the queue backing array for our topological
	// traversal queue.
	queue = queue[:0]
	pkgVal := map[projpath.FilePath]cty.Value{}
	for depDir := range packages {
		if inDeg[depDir] == 0 {
			queue = append(queue, depDir)
		}
	}
	for len(queue) > 0 {
		thisDir := queue[0]
		queue = queue[1:]

		thisPkg := packages[thisDir]
		depVals := map[string]cty.Value{}
		for _, dep := range thisPkg.PackagesImported() {
			depDir := proj.FilePathForPackagePath(dep.Path, dep.Range.Filename)
			if depDir == projpath.NoPath {
				// should never happen unless the filesystem is drifting
				// while we're working, so we'll ignore it.
				continue
			}

			if depVal, exists := pkgVal[depDir]; exists {
				depVals[dep.Path] = depVal
			} else {
				// This is indicative of an evaluation error, but we'll
				// allow it to pass through here so we can potentially
				// collect more diagnostics on subsequent loops.
				depVals[dep.Path] = cty.PlaceholderVal
			}
		}

		fmt.Printf("evaluating %#v\n", thisDir)
		result, evalDiags := thisPkg.ExportedValue(depVals)
		diags = append(diags, evalDiags...)
		pkgVal[thisDir] = result

		for _, enabledDir := range enables[thisDir] {
			inDeg[enabledDir]--
			if inDeg[enabledDir] < 1 {
				queue = append(queue, enabledDir)
				delete(inDeg, enabledDir)
			}
		}
	}

	// TODO: check if anything is left in inDeg, which would indicate
	// a dependency cycle.

	if diags.HasErrors() {
		return diagsError(diags)
	}

	finalResult := pkgVal[pkgDir]
	fmt.Printf("exported value is %#v\n", finalResult)

	return nil
}

func compilePackage(pkgDir projpath.FilePath, proj projpath.Project, p *parser.Parser) (*cbo.Package, source.Diags) {
	var diags source.Diags
	modFilenames := proj.ListModuleFiles(pkgDir)

	files := make([]*ast.File, 0, len(modFilenames))
	for _, filename := range modFilenames {
		src, err := proj.ReadFile(filename)
		if err != nil {
			diags = append(diags, source.Diag{
				Level:   source.Error,
				Summary: "Failed to read module file",
				Detail: fmt.Sprintf(
					"The module file %s (for package in %s) could not be read: %s",
					proj.FilePathForUI(filename),
					proj.FilePathForUI(pkgDir),
					err.Error(),
				),
			})
			continue
		}

		file, fileDiags := p.ParseFile(filename, src)
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

func diagsError(diags source.Diags) error {
	if len(diags) > 0 {
		os.Stderr.WriteString("\n")
		for _, diag := range diags {
			fmt.Fprintf(os.Stderr, "- %s\n", diag.String())
		}
		os.Stderr.WriteString("\n")
	}
	if diags.HasErrors() {
		return errors.New("There were some errors during parsing, as shown above.")
	}
	return nil
}
