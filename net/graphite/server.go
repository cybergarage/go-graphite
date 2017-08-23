// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package graphite

// Server is an instance for Graphite protocols.
type Server struct {
	*Carbon
	*Render
}

// NewServer returns a new Server.
func NewServer() *Server {
	server := &Server{}
	server.Carbon = NewCarbon()
	server.Render = NewRender()
	return server
}

// Start starts the server.
func (self *Server) Start() error {
	err := self.Carbon.Start()
	if err != nil {
		self.Stop()
		return err
	}

	err = self.Render.Start()
	if err != nil {
		self.Stop()
		return err
	}

	return nil
}

// Stop stops the server.
func (self *Server) Stop() error {
	err := self.Carbon.Stop()
	if err != nil {
		return err
	}

	err = self.Render.Stop()
	if err != nil {
		return err
	}

	return nil
}
