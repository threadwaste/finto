SHELL := /bin/bash

.DEFAULT: default

all: deps build fmt vet test install

build:
	go build -a -v ./...

deps:
	go get -t -v ./...

fmt:
	diff -u <(echo -n) <(gofmt -s -d .)

install:
	go install -v .

test:
	go test -v ./...

testall:
	go test -v -tags integration ./...

vet:
	go vet ./...
