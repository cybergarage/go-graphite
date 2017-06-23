// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

// Server is an instance for Graphite protocols.
type Server struct {
	carbon     *Carbon
	httpServer *HTTPServer
}

// NewServer returns a new Server.
func NewServer() *Server {
	server := &Server{}
	server.carbon = NewCarbon()
	server.httpServer = NewHTTPServer()
	return server
}
