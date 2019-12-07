TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build:
	go install

test:
	go test ./...

testacc:
	TF_LOG=debug TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

.PHONY: build test testacc
