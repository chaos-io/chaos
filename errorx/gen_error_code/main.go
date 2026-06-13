package main

import (
	"fmt"
	"os"

	"github.com/chaos-io/chaos/errorx/gen_error_code/internal/generator"
)

//go:generate go run . -out ./testdata/generated ./testdata

func main() {
	if err := generator.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
