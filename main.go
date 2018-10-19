package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	errUsage   = fmt.Errorf(`usage: %s {changelog|release}`, os.Args[0])
	errVersion = fmt.Errorf("version required")
)

func main() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	var err error
	switch cmd {
	case "release":
		flags := flag.NewFlagSet("release", flag.ExitOnError)
		prefix := flags.String("prefix", "", "release branch prefix")
		version := flags.String("version", "", "release version")
		verbose := flags.Bool("verbose", false, "verbose output")
		flags.Parse(os.Args[2:])
		if *version == "" {
			err = errVersion
		} else {
			err = release(*prefix, *version, *verbose)
		}
	case "changelog":
		flags := flag.NewFlagSet("changelog", flag.ExitOnError)
		prefix := flags.String("prefix", "", "release branch prefix")
		output := flags.String("output", "-", "output file to prepend changelog")
		flags.Parse(os.Args[2:])
		err = generateChangelog(*prefix, *output)
	default:
		err = errUsage
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
