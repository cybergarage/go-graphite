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

PACKAGE_NAME=net/graphite
GITHUB=github.com/cybergarage/go-graphite

GITHUB_ID=${GITHUB}.git/${PACKAGE_NAME}
PACKAGE_ID=${GITHUB}/${PACKAGE_NAME}

PACKAGES=${PACKAGE_ID}

.PHONY: setup

all: setup test

VERSION_GO="./net/graphite/version.go"

${VERSION_GO}: ./net/graphite/version.gen
	$< > $@

version: ${VERSION_GO}

SETUP_CMD="./setup"

setup:
	@echo "#!/bin/bash" > ${SETUP_CMD}
	@echo "export GOPATH=\`pwd\`" >> ${SETUP_CMD}
	@echo "git pull" >> ${SETUP_CMD}
	@echo "pushd src && rm -rf ${GITHUB}.git ${GITHUB} && popd" >> ${SETUP_CMD}
	@echo "go get -u ${GITHUB_ID}" >> ${SETUP_CMD}
	@echo "pushd src && mv ${GITHUB}.git ${GITHUB} && popd" >> ${SETUP_CMD}
	@chmod a+x ${SETUP_CMD}
	@./${SETUP_CMD}

format:
	gofmt -w src/${GITHUB}

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
