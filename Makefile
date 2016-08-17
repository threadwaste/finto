SHELL = /bin/bash

FINTO_ROOT=github.com/threadwaste/finto
FINTO_PACKAGES=${FINTO_ROOT} ${FINTO_ROOT}/cmd/finto
FINTO_NOVENDOR=$(shell find . -type f -name \*.go -not -path ./vendor/\*)

all: build fmt vet test install

build:
	go build -a -v ./...

fmt:
	diff -u <(echo -n) <(gofmt -s -d ${FINTO_NOVENDOR})

install:
	go install -v .

test:
	go test -v ${FINTO_PACKAGES}

testall:
	go test -v -tags integration ${FINTO_PACKAGES}

vet:
	go vet -x ${FINTO_PACKAGES}
