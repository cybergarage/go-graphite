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

MODULE_ROOT=github.com/cybergarage/go-graphite

PKG_PREFIX=net
PKG_NAME=graphite
PKG_COVER=${PKG_NAME}-cover
PKG_ID=${MODULE_ROOT}/${PKG_PREFIX}/${PKG_NAME}
PKG_SRC_DIR=${PKG_PREFIX}/${PKG_NAME}
PKG_SRCS=\
        ${PKG_SRC_DIR}
PKGS=\
        ${PKG_ID}

.PHONY: format vet lint cover clean

all: test

VERSION_GO="./net/graphite/version.go"

${VERSION_GO}: ./net/graphite/version.gen
	$< > $@

version: ${VERSION_GO}

format:
	gofmt -w ${PKG_PREFIX}/${PKG_NAME}

vet: format
	go vet ${PKG_ID}

lint: vet
	golangci-lint run ${PKG_SRCS}

test:
	 go test -v -timeout 60s ${PKGS} -cover -coverpkg=${PKG_ID} -coverprofile=${PKG_COVER}.out

clean:
	rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PKGS}
