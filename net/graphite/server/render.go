// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

import (
	"net/http"
)

const (
	// HTTPDefaultPort is the default port number for HTTP Server
	HTTPDefaultPort int    = 8080
	QueryTarget     string = "target"
	QueryFrom       string = "from"
	QueryUntil      string = "until"
	QueryFormat     string = "format"
	QueryFormatCSV  string = "csv"
	QueryFormatJSON string = "json"
)

// Render is an instance for Graphite render protocols.
type Render struct {
	Port   int
	server http.Server
}

// NewRender returns a new Render.
func NewRender() *Render {
	server := &Render{Port: HTTPDefaultPort}
	return server
}

// Start starts the HTTP server.
func (self *Render) Start() error {
	err := self.Stop()
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the HTTP server.
func (self *Render) Stop() error {
	return nil
}

// ServeHTTP handles HTTP requests.
// Support The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Render) ServeHTTP(httpWriter http.ResponseWriter, httpReq *http.Request) {

}
