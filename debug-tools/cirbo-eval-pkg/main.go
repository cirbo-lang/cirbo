package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cirbo-lang/cirbo/cirbo"
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

	cb := cirbo.New(cirbo.Config{
		WorkingDir:   wd,
		SystemPkgDir: wd, // we don't have a SystemPkgDir for this debug tool
	})

	value, diags := cb.LoadPackage(args[0])
	if diags.HasErrors() {
		return diagsError(diags)
	}

	fmt.Printf("exported value is %#v\n", value)

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
		return errors.New("There were some errors, as shown above.")
	}
	return nil
}
