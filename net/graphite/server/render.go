// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

import (
	"fmt"
	"net/http"

	"github.com/cybergarage/go-graphite/net/graphite"
)

const (
	// DefaultPort is the default port number for Render
	DefaultPort int = 8080
	// DefaultPath is the default path for Render
	DefaultPath string = "/render"
)

// RenderRequestListener represents a listener for Render protocol.
type RenderRequestListener interface {
	QueryRequestReceived(*graphite.Query, error)
}

// RenderListener represents a listener for all requests of Render.
type RenderListener interface {
	RenderRequestListener
}

// Render is an instance for Graphite render protocols.
type Render struct {
	Port     int
	Listener RenderRequestListener
	server   *http.Server
}

// NewRender returns a new Render.
func NewRender() *Render {
	server := &Render{
		Port:   DefaultPort,
		server: nil,
	}
	return server
}

// Start starts the HTTP server.
func (self *Render) Start() error {
	err := self.Stop()
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%d", self.Port)

	self.server = &http.Server{
		Addr:    addr,
		Handler: self,
	}

	// FIXE : Handle error
	go self.server.ListenAndServe()
	/*
		err = go self.server.ListenAndServe()
		if err != nil {
			return err
		}
	*/

	return nil
}

// Stop stops the HTTP server.
func (self *Render) Stop() error {
	if self.server == nil {
		return nil
	}

	err := self.server.Close()
	if err != nil {
		return err
	}

	return nil
}

// ServeHTTP handles HTTP requests.
func (self *Render) ServeHTTP(httpWriter http.ResponseWriter, httpReq *http.Request) {

	switch httpReq.URL.Path {
	case DefaultPath:
		self.handleRenderRequest(httpWriter, httpReq)

	}

	http.NotFound(httpWriter, httpReq)
}

// handleRenderRequest handles Render requests.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Render) handleRenderRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	q := graphite.NewQuery()
	err := q.Parse(httpReq.URL)
	if err != nil {
		http.NotFound(httpWriter, httpReq)
		return
	}

	http.NotFound(httpWriter, httpReq)
}
