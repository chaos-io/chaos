package main

import (
	"fmt"
	"strings"
)

func main() {
	files := []string{"a", "x/b"}
	covs := []string{"/d/a", "/d/x/b"}

	stop := false
	for _, file := range files {
		for _, cov := range covs {
			if strings.HasSuffix(cov, file) {
				SourceCodeFilePrefix, _ := strings.CutSuffix(cov, file)
				fmt.Println("----", SourceCodeFilePrefix)
				stop = true
				break
			}
		}

		if stop {
			break
		}
	}
}
