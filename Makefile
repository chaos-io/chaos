.PHONY: env test

GOPROXY := https://goproxy.cn,direct
export GOPROXY

default: test

env:
	@go version
test: env
	go test -race -covermode=atomic -coverprofile=coverage.out.tmp -coverpkg ./... ./...
	@cat coverage.out.tmp | grep -v "mock_" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out -o coverage.txt
	@tail -n 1 coverage.txt
