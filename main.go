package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var errNoProgramName = errors.New("no program name (os.Args empty)")

func parseArgs(osArgs []string) error {
	if len(osArgs) == 0 {
		return errNoProgramName
	}

	fs := flag.NewFlagSet(osArgs[0], flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), "usage: ", fs.Name(), "\n")
		fmt.Fprint(fs.Output(), `
Fetch a webpage regularly and send an email when it changes.

See the README for how to set parameters.
`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(osArgs[1:]); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unrecognized arguments: %#v", fs.Args())
	}
	return nil
}

func errExit(err error) {
	if len(os.Args) > 0 && os.Args[0] != "" {
		fmt.Fprint(os.Stderr, os.Args[0], ": ")
	}
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func main_err() error {
	err := parseArgs(os.Args)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := main_err(); err != nil {
		errExit(err)
	}
}
