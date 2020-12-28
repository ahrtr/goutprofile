// +build ignore

package main

import (
	"fmt"
	gup "github.com/ahrtr/goutprofile"
	"os"
)

func main() {
	if _, err := gup.ValidateDir(".", false); err != nil {
		fmt.Printf("Error validating goutprofile, error: %v\n", err)
		os.Exit(1)
	}
}
