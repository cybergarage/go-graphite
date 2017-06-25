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

// Start starts the server.
func (self *Server) Start() error {
	err := self.carbon.Start()
	if err != nil {
		self.Stop()
		return err
	}

	err = self.httpServer.Start()
	if err != nil {
		self.Stop()
		return err
	}

	return nil
}

// Stop stops the server.
func (self *Server) Stop() error {
	err := self.carbon.Stop()
	if err != nil {
		return err
	}

	err = self.httpServer.Stop()
	if err != nil {
		return err
	}

	return nil
}
