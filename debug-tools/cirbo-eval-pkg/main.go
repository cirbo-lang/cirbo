package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cirbo-lang/cirbo/compiler"
	"github.com/cirbo-lang/cirbo/projpath"
	"github.com/cirbo-lang/cirbo/source"

	"github.com/cirbo-lang/cirbo/ast"
	"github.com/cirbo-lang/cirbo/parser"
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
	modFilenames := proj.ListModuleFiles(pkgDir)

	files := make([]*ast.File, 0, len(modFilenames))
	for _, filename := range modFilenames {
		src, err := proj.ReadFile(filename)
		if err != nil {
			return err
		}

		file, fileDiags := p.ParseFile(filename, src)
		diags = append(diags, fileDiags...)
		files = append(files, file)
	}

	if diags.HasErrors() {
		return diagsError(diags)
	}

	pkgNode := &ast.Package{
		// FIXME: This is no longer the right place to retain this, since
		// the parser no longer wants to know about packages.
		DefaultName: "FIXME",
		Files:       files,
	}

	pkg, compileDiags := compiler.CompilePackage(pkgNode)
	diags = append(diags, compileDiags...)

	if diags.HasErrors() {
		return diagsError(diags)
	}

	// TODO: We don't have enough implemented yet to actually eval, so for
	// now we'll just print out some info.

	fmt.Printf("module imports %#v\n", pkg.PackagesImported())

	return nil
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
