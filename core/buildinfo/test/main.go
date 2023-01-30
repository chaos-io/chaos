package main

import (
	"fmt"

	"github.com/chaos-io/chaos/core/buildinfo"
)

func main() {
	if buildinfo.Info.ProgramVersion != "" {
		fmt.Print(buildinfo.Info.ProgramVersion)
	} else {
		fmt.Printf("ProgramVersion is not available\n")
	}
}
