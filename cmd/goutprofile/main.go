package main

import (
	"flag"
	"fmt"
	"os"

	gup "github.com/ahrtr/goutprofile"
)

var (
	dir         string
	file        string
	recursive   bool
	showVersion bool

	// Version is the version of GO UT Profile tool
	Version string
)

func init() {
	flag.StringVar(&dir, "d", "", "directory to be validated")
	flag.StringVar(&file, "f", "", "profile to be validated")
	flag.BoolVar(&recursive, "r", false, "whether do recursive validation")
	flag.BoolVar(&showVersion, "v", false, "show version")
}

func main() {
	flag.Parse()
	validateConfig()

	if len(file) > 0 {
		if err := gup.ValidateFile(file); err != nil {
			fmt.Printf("Error validating file, %s: %v\n", file, err)
		}
		return
	}

	if profiles, err := gup.ValidateDir(dir, recursive); err != nil {
		fmt.Printf("Error validating dir, %s: %v\n", dir, err)
	} else {
		fmt.Printf("%d profile(s) validated\n", len(profiles))
	}
}

func validateConfig() {
	if showVersion {
		fmt.Printf("GO UT Profile version: %s\n", Version)
		os.Exit(0)
	}

	if len(dir) == 0 && len(file) == 0 {
		showErrorMessage("'-d' or '-f' must be set")
	}

	if len(dir) > 0 && len(file) > 0 {
		showErrorMessage("'-d' and '-f' can't be set at the same time")
	}

	if len(file) > 0 && recursive {
		showErrorMessage("'-f' and '-r' can't be set at the same time")
	}
}

func showErrorMessage(errMsg string) {
	fmt.Println(errMsg)
	flag.Usage()
	os.Exit(1)
}
