PREFIX?=$(shell pwd)

.DEFAULT: default

all: deps build fmt test vet install

build:
	go build -a -v ./...

deps:
	go get -t -v ./...

fmt:
	gofmt -s -l .

install:
	go install -v .

test:
	go test -v ./...

vet:
	go vet ./...
