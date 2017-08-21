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

// HTTPServer is an instance for Graphite Web protocols.
type HTTPServer struct {
	Port   int
	server http.Server
}

// NewHTTPServer returns a new HTTPServer.
func NewHTTPServer() *HTTPServer {
	server := &HTTPServer{Port: HTTPDefaultPort}
	return server
}

// Start starts the HTTP server.
func (self *HTTPServer) Start() error {
	err := self.Stop()
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the HTTP server.
func (self *HTTPServer) Stop() error {
	return nil
}

// ServeHTTP handles HTTP requests.
// Support The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *HTTPServer) ServeHTTP(httpWriter http.ResponseWriter, httpReq *http.Request) {

}
