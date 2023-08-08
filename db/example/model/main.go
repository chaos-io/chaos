package main

import (
	"fmt"

	"github.com/chaos-io/chaos/db/example/model/templates"
)

func main() {
	m := templates.InitModel()
	generateFile, err := m.Generate()
	if err != nil {
		fmt.Printf("generate error: %v", err)
		return
	}
	fmt.Printf("generate file into %s", generateFile)
}
