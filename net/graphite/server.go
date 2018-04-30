// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	server.SetAddress("")
	return server
}

// SetAddress sets a bind address.
func (self *Server) SetAddress(addr string) error {
	self.Carbon.Addr = addr
	self.Render.Addr = addr
	return nil
}

// GetAddress returns the bind address.
func (self *Server) GetAddress() string {
	return self.Render.Addr
}

// SetCarbonPort sets a bind port for Carbon.
func (self *Server) SetCarbonPort(port int) error {
	self.Carbon.Port = port
	return nil
}

// GetCarbonPort returns a bind port for Carbon.
func (self *Server) GetCarbonPort() int {
	return self.Carbon.Port
}

// SetRenderPort sets a bind port for Render.
func (self *Server) SetRenderPort(port int) error {
	self.Render.Port = port
	return nil
}

// GetRenderPort returns a bind port for Render.
func (self *Server) GetRenderPort() int {
	return self.Render.Port
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
