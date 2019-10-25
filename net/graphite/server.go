// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"net"
	"time"
)

// Server is an instance for Graphite protocols.
type Server struct {
	boundInterface *net.Interface
	*Carbon
	*Render
}

// NewServer returns a new Server.
func NewServer() *Server {
	server := &Server{
		boundInterface: nil,
		Carbon:         NewCarbon(),
		Render:         NewRender(),
	}
	return server
}

// SetConfig sets a configuration to the server.
func (server *Server) SetConfig(conf *Config) {
	server.SetCarbonPort(conf.GetCarbonPort())
	server.SetRenderPort(conf.GetRenderPort())
	server.SetConnectionTimeout(conf.GetConnectionTimeout())
	server.SetConnectionWaitTimeout(conf.GetConnectionWaitTimeout())
}

// SetBoundInterface sets a bound interface to the server.
func (server *Server) SetBoundInterface(ifi *net.Interface) {
	server.boundInterface = ifi
}

// GetBoundInterface returns the bound interface.
func (server *Server) GetBoundInterface() *net.Interface {
	return server.boundInterface
}

// SetBoundAddress sets a bound address.
func (server *Server) SetBoundAddress(addr string) {
	server.Carbon.SetAddress(addr)
	server.Render.SetAddress(addr)
}

// GetBoundAddress returns the bound address.
func (server *Server) GetBoundAddress() string {
	return server.Render.GetAddress()
}

// SetCarbonPort sets a bind port for Carbon.
func (server *Server) SetCarbonPort(port int) {
	server.Carbon.SetPort(port)
}

// GetCarbonPort returns a bind port for Carbon.
func (server *Server) GetCarbonPort() int {
	return server.Carbon.GetPort()
}

// SetRenderPort sets a bind port for Render.
func (server *Server) SetRenderPort(port int) {
	server.Render.SetPort(port)
}

// GetRenderPort returns a bind port for Render.
func (server *Server) GetRenderPort() int {
	return server.Render.GetPort()
}

// SetConnectionTimeout sets the connection timeout for Render .
func (server *Server) SetConnectionTimeout(d time.Duration) {
	server.Render.SetConnectionTimeout(d)
}

// SetConnectionWaitTimeout sets the connection wait timeout for Carbon .
func (server *Server) SetConnectionWaitTimeout(d time.Duration) {
	server.Carbon.SetConnectionWaitTimeout(d)
}

// GetConnectionTimeout return the connection timeout.
func (server *Server) GetConnectionTimeout() time.Duration {
	return server.Render.GetConnectionTimeout()
}

// Start starts the server.
func (server *Server) Start() error {
	err := server.Carbon.Start()
	if err != nil {
		server.Stop()
		return err
	}

	err = server.Render.Start()
	if err != nil {
		server.Stop()
		return err
	}

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	err := server.Carbon.Stop()
	if err != nil {
		return err
	}

	err = server.Render.Stop()
	if err != nil {
		return err
	}

	return nil
}
