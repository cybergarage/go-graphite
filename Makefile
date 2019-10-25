###################################################################
#
# go-grapite
#
# Copyright (C) The go-graphite Authors 2017
#
# This is licensed under BSD-style license, see file COPYING.
#
###################################################################

PREFIX?=$(shell pwd)
GOPATH=$(shell pwd)

MODULE_NAME_ROOT=github.com/cybergarage

PACKAGE_NAME=net/graphite
MODULE_NAME=go-graphite

PACKAGE_ID=${MODULE_NAME}/${PACKAGE_NAME}

PACKAGES=${PACKAGE_ID}

.PHONY:

all: test

VERSION_GO="./net/graphite/version.go"

${VERSION_GO}: ./net/graphite/version.gen
	$< > $@

version: ${VERSION_GO}

format:
	gofmt -w ${PACKAGE_NAME}

package: format $(shell find src/${PACKAGE_ID}  -type f -name '*.go')
	go build -v ${PACKAGES}

test: package
	go test -v -cover ${PACKAGES}

install: build
	go install ${PACKAGES}

clean:
	rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PACKAGES}
