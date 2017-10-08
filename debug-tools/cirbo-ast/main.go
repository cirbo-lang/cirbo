package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

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
	fl := flag.NewFlagSet("cirbo-ast", flag.ExitOnError)
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

	fpath := proj.FilePathFromUI(args[0])
	src, err := proj.ReadFile(fpath)
	if err != nil {
		return err
	}

	p := parser.NewParser()
	f, diags := p.ParseFile(fpath, src)

	ast.Walk(f, &walker{})

	if len(diags) > 0 {
		os.Stderr.WriteString("\n")
		for _, diag := range diags {
			fmt.Fprintf(os.Stderr, "- %s\n", diag.String())
		}
		os.Stderr.WriteString("\n")
		return errors.New("There were some errors during parsing, as shown above.")
	}
	return nil
}

type walker struct {
	Depth int
}

func (w *walker) EnterNode(node ast.Node) bool {
	buf := &bytes.Buffer{}
	buf.WriteString(strings.Repeat("  ", w.Depth))
	buf.WriteString(fmt.Sprintf("%T", node)[5:])
	nodeType := reflect.ValueOf((*ast.Node)(nil)).Type().Elem()
	nv := reflect.ValueOf(node).Elem()
	if nv.Kind() == reflect.Struct {
		fieldCt := nv.NumField()
		ty := nv.Type()
		for i := 0; i < fieldCt; i++ {
			f := ty.Field(i)
			fv := nv.Field(i)
			if f.Type.ConvertibleTo(nodeType) {
				// Don't show nodes since we'll walk into them
				continue
			}
			if f.Type.Kind() == reflect.Slice && f.Type.Elem().ConvertibleTo(nodeType) {
				// Don't show slices of nodes either
				continue
			}
			if f.Type.AssignableTo(reflect.TypeOf([]byte(nil))) {
				// Don't show source slices
				continue
			}
			if f.Type.AssignableTo(reflect.TypeOf(ast.WithRange{})) || f.Type.AssignableTo(reflect.TypeOf(source.Range{})) {
				// Don't show ranges
				continue
			}
			fmt.Fprintf(buf, " %s=%#v", f.Name, fv.Interface())
		}
	}
	buf.WriteString("\n")
	fmt.Print(buf.String())
	w.Depth++
	return true
}

func (w *walker) ExitNode(node ast.Node) {
	w.Depth--
}
