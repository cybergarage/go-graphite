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
PKGS=\
        ${PKG_ID}

TEST_PKG_NAME=test
TEST_PKG_ID=${MODULE_ROOT}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG=${MODULE_ROOT}/${TEST_PKG_DIR}

BIN_DIR=examples
BIN_ID=${MODULE_ROOT}/${BIN_DIR}
BIN_SERVER=go-graphited
BIN_SERVER_ID=${BIN_ID}/${BIN_SERVER}
BINS=\
	${BIN_SERVER_ID}

.PHONY: format vet lint cover clean test install

all: test

VERSION_GO="./net/graphite/version.go"

${VERSION_GO}: ./net/graphite/version.gen
	$< > $@

version: ${VERSION_GO}

format:
	gofmt -w ${PKG_SRC_DIR} ${TEST_PKG_DIR} ${BIN_DIR}

vet: format
	go vet ${PKG_ID}

lint: vet
	golangci-lint run ${PKG_SRC_DIR}/... ${BIN_DIR}/... ${TEST_PKG_DIR}/...

test:
	go test -v -timeout 60s ${PKGS} ${TEST_PKG} -cover -coverpkg=${PKG_ID} -coverprofile=${PKG_COVER}.out
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

install: test
	go install ${BINS}

clean:
	rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PKGS}
