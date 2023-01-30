//go:build !arcadia
// +build !arcadia

package yatest

import (
	"github.com/chaos-io/chaos/yatool"
)

func doInit() {
	isRunningUnderGoTest = true

	arcadiaSourceRoot, err := yatool.ArcadiaRoot()
	if err != nil {
		panic(err)
	}
	context.Runtime.SourceRoot = arcadiaSourceRoot
}
