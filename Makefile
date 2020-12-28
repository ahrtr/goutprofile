.PHONY: all test unitTest build validate

VERSION ?= $(shell git rev-parse --short HEAD)

all: test validate build

test: unitTest 

unitTest:
	go test ${TEST_OPTS} .

validate:
	go run cmd/goutprofile/main.go -d . 
	
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o goutprofile ./cmd/goutprofile/main.go
