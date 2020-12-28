.PHONY: all test unitTest build

VERSION ?= $(shell git rev-parse --short HEAD)

all: test build

test: unitTest 

unitTest:
	go test ${TEST_OPTS} .
	
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o gup ./cmd/goutprofile/main.go
