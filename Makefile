SHELL := /bin/bash

.DEFAULT: default

all: deps build fmt test vet install

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

vet:
	go vet ./...
