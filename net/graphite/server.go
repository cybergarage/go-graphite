// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import "net"

// Server is an instance for Graphite protocols.
type Server struct {
	*net.Interface
	*Carbon
	*Render
}

// NewServer returns a new Server.
func NewServer() *Server {
	server := &Server{
		Interface: nil,
		Carbon:    NewCarbon(),
		Render:    NewRender(),
	}
	return server
}

// SetConfig sets a configuration to the server.
func (server *Server) SetConfig(conf *Config) {
	server.SetCarbonPort(conf.GetCarbonPort())
	server.SetRenderPort(conf.GetRenderPort())
}

// SetInterface sets a bound interface to the server.
func (server *Server) SetInterface(ifi *net.Interface) error {
	ifaddr, err := GetInterfaceAddress(ifi)
	if err != nil {
		return err
	}

	server.Interface = ifi
	server.SetAddress(ifaddr)

	return nil
}

// GetInterface returns the bound interface.
func (server *Server) GetInterface() *net.Interface {
	return server.Interface
}

// SetAddress sets a bind address.
func (server *Server) SetAddress(addr string) {
	server.Carbon.SetAddress(addr)
	server.Render.SetAddress(addr)
}

// GetAddress returns the bound address.
func (server *Server) GetAddress() string {
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
