###################################################################
#
# go-grapite-server
#
# Copyright (C) The go-graphite Authors 2017
#
# This is licensed under BSD-style license, see file COPYING.
#
###################################################################

PREFIX?=$(shell pwd)
GOPATH=$(shell pwd)

GITHUB=github.com/cybergarage/go-graphite-server
PACKAGES=${GITHUB}/net/graphite/server

.PHONY: setup

VERSION_GO="./net/graphite/server/version.go"

${VERSION_GO}: ./net/graphite/server/version.gen
	$< > $@

version: ${VERSION_GO}

SETUP_CMD="./setup"

setup:
	@echo "export GOPATH=${GOPATH}" > ${SETUP_CMD}
	@echo "go get -u ${GITHUB}/net/graphite/server" >> ${SETUP_CMD}
	@chmod a+x ${SETUP_CMD}
	@./${SETUP_CMD}

commit:
	cd src/${GITHUB} && git commit

format:
	gofmt -w src/${GITHUB} net

package: format $(shell find . -type f -name '*.go')
	go build -v ${PACKAGES}

test: package
	go test -v -cover ${PACKAGES}

install: build
	go install ${PACKAGES}

clean:
	rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PACKAGES}
