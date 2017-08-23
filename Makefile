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

GITHUB=github.com/cybergarage/go-graphite
PACKAGE_ID=${GITHUB}/net/graphite
PACKAGES=${PACKAGE_ID} ${PACKAGE_ID}/client ${PACKAGE_ID}/server

.PHONY: setup

VERSION_GO="./net/graphite/version.go"

${VERSION_GO}: ./net/graphite/version.gen
	$< > $@

version: ${VERSION_GO}

SETUP_CMD="./setup"

setup:
	@echo "#!/bin/bash" > ${SETUP_CMD}
	@echo "export GOPATH=\`pwd\`" >> ${SETUP_CMD}
	@echo "git pull" >> ${SETUP_CMD}
	@echo "go get -u ${PACKAGE_ID}" >> ${SETUP_CMD}
	@chmod a+x ${SETUP_CMD}
	@./${SETUP_CMD}

commit:
	pushd src/${GITHUB} && git commit -a && popd

push:
	pushd src/${GITHUB} && git push && popd

pull:
	pushd src/${GITHUB} && git pull && popd

diff:
	pushd src/${GITHUB} && git diff && popd

format:
	gofmt -w src/${GITHUB}

package: format $(shell find src/${GITHUB}  -type f -name '*.go')
	go build -v ${PACKAGES}

test: package
	go test -v -cover ${PACKAGES}

install: build
	go install ${PACKAGES}

clean:
	rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PACKAGES}
