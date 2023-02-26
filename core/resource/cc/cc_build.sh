#!/bin/bash
package=main
o=test1.go
registryKey=a-server
blob=111

set -x
go run main.go -package=$package -o=$o - notafile$registryKey=$blob
