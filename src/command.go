package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Actions struct {
	Help     bool
	Register bool
	Bind     bool
	UnBind   bool
	Cancle   bool
	License  string
	FileName string
	Name     string
	Generate string
}

func ParseCommandLine() Actions {
	var action Actions

	flag.BoolVar(&action.Help, "h", false, "")
	flag.BoolVar(&action.Help, "help", false, "")

	flag.BoolVar(&action.Register, "r", false, "")
	flag.BoolVar(&action.Register, "register", false, "")

	flag.BoolVar(&action.Bind, "b", false, "")
	flag.BoolVar(&action.Bind, "bind", false, "")

	flag.BoolVar(&action.UnBind, "u", false, "")
	flag.BoolVar(&action.UnBind, "unbind", false, "")

	flag.BoolVar(&action.Cancle, "c", false, "")
	flag.BoolVar(&action.Cancle, "cancle", false, "")

	flag.StringVar(&action.FileName, "f", "", "")
	flag.StringVar(&action.FileName, "file", "", "")

	flag.StringVar(&action.License, "l", "", "")
	flag.StringVar(&action.License, "license", "", "")

	flag.StringVar(&action.Name, "n", "", "")
	flag.StringVar(&action.Name, "name", "", "")

	flag.StringVar(&action.Generate, "g", "", "")
	flag.StringVar(&action.Generate, "generate", "", "")

	flag.Visit(
		func(f *flag.Flag) {
			if f.Name == "h" || f.Name == "help" {
				action.Help = true
			}
			if f.Name == "r" || f.Name == "register" {
				action.Register = true
			}
			if f.Name == "b" || f.Name == "bind" {
				action.Bind = true
			}
			if f.Name == "u" || f.Name == "unbind" {
				action.UnBind = true
			}
			if f.Name == "c" || f.Name == "cancle" {
				action.Cancle = true
			}
			if f.Name == "f" || f.Name == "file" {
				action.FileName = f.Value.String()
			}
			if f.Name == "l" || f.Name == "license" {
				action.License = f.Value.String()
			}
			if f.Name == "n" || f.Name == "name" {
				action.Name = f.Value.String()
			}
			if f.Name == "g" || f.Name == "generate" {
				action.Generate = f.Value.String()
			}
		},
	)

	flag.Usage = func() {
		help()
	}

	flag.Parse()

	if action.License != "" {
		expectedPattern := `^[0-9A-Za-z]{8}-[0-9A-Za-z]{8}-[0-9A-Za-z]{8}$`

		match, err := regexp.MatchString(expectedPattern, action.License)
		if err != nil {
			panic(err)
		}

		if !match {
			panic("License should be something matchs: ^[0-9A-Za-z]{8}-[0-9A-Za-z]{8}-[0-9A-Za-z]{8}$")
		}
	}

	return action
}

func help() {
	fmt.Fprintf(os.Stderr,
		`wg-cli Revision: __REVISION__
Usage:	%s [Options]
Options:		-h/--help			help
			-f/--file [string]		Configuration file (default "wgcf.json")
			-r/--register			Register an account
			-b/--bind			Get the account binding devices
			-n/--name [string]		Change the device name
			-l/--license [string]		Change the license
			-u/--unbind			Unbind a device from the account
			-c/--cancle			Cancle the account
			-g/--generate [wg/xray]		Generate a [wg/xray] configuration file
`, os.Args[0])
}
