SHELL = /bin/bash

BUILD_BIN=_build_bin

GOVERSION:=$(shell go version)
GOOS=$(word 1,$(subst /, , $(lastword ${GOVERSION})))
GOARCH=$(word 2,$(subst /, , $(lastword ${GOVERSION})))

FINTO_ROOT=github.com/threadwaste/finto
FINTO_MAIN=${FINTO_ROOT}/cmd/finto
FINTO_PACKAGES=${FINTO_ROOT} ${FINTO_MAIN}
FINTO_NOVENDOR:=$(shell find . -type f -name \*.go -not -path ./vendor/\*)

HAVE_GLIDE:=$(shell which glide)

.PHONY: build fmt vet

${BUILD_BIN}/glide:
ifndef HAVE_GLIDE
	@mkdir -vp ${BUILD_BIN}
	@curl -vL https://github.com/Masterminds/glide/releases/download/v0.12.0/glide-v0.12.0-${GOOS}-${GOARCH}.tar.gz | tar zxv -C ${BUILD_BIN}
endif

build:
	go build -a -v ${FINTO_MAIN}

deps: glide ${FINTO_NOVENDOR}
	@PATH=${BUILD_BIN}/${GOOS}-${GOARCH}:${PATH} glide install

fmt:
	diff -u <(echo -n) <(gofmt -s -d ${FINTO_NOVENDOR})

glide: ${BUILD_BIN}/glide

test: deps
	go test -v ${FINTO_PACKAGES}

testall: deps
	go test -v -tags integration ${FINTO_PACKAGES}

vet:
	go vet -x ${FINTO_PACKAGES}
