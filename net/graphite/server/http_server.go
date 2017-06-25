// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

const (
	// HTTPDefaultPort is the default port number for HTTP Server
	HTTPDefaultPort int = 8080
)

// HTTPServer is an instance for Graphite Web protocols.
type HTTPServer struct {
	Port int
}

// NewHTTPServer returns a new HTTPServer.
func NewHTTPServer() *HTTPServer {
	server := &HTTPServer{Port: HTTPDefaultPort}
	return server
}

// Start starts the HTTP server.
func (self *HTTPServer) Start() error {
	return nil
}

// Stop stops the HTTP server.
func (self *HTTPServer) Stop() error {
	return nil
}
