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

VERSION_GO="./net/graphite/server/version.go"

GITHUB=github.com/cybergarage/go-graphite-server

PACKAGES=${GITHUB}/net/graphite/server
	
${VERSION_GO}: ./net/graphite/server/version.gen
	$< > $@

version: ${VERSION_GO}

setup:
	export GOPATH=${GOPATH}
	go get -u ${GITHUB}/net/graphite/server

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
