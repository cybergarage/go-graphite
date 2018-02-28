// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"net/http"
)

const (
	// RenderDefaultPort is the default port number for Render and Metrics APIs
	RenderDefaultPort int = 8080
)

const (
	httpHeaderContentType                 = "Content-Type"
	renderDefaultFindRequestPath   string = "/metrics/find"
	renderDefaultExpandRequestPath string = "/metrics/expand"
	renderDefaultIndexRequestPath  string = "/metrics/index.json"
	renderDefaultQueryRequestPath  string = "/render"
	renderMetricsDelim             string = "."
	renderMetricsAsterisk          string = "*"
)

// RenderRequestListener represents a listener for Render protocol.
type RenderRequestListener interface {
	FindMetricsRequestReceived(*Query, error) ([]*Metrics, error)
	QueryMetricsRequestReceived(*Query, error) ([]*Metrics, error)
}

// RenderHTTPRequestListener represents a listener for HTTP requests.
type RenderHTTPRequestListener interface {
	HTTPRequestReceived(r *http.Request, w http.ResponseWriter)
}

// RenderListener represents a listener for all requests of Render.
type RenderListener interface {
	RenderRequestListener
}

// Render is an instance for Graphite render protocols.
type Render struct {
	Port               int
	RenderListener     RenderRequestListener
	server             *http.Server
	extraHTTPListeners map[string]RenderHTTPRequestListener
}

// NewRender returns a new Render.
func NewRender() *Render {
	server := &Render{
		Port:               RenderDefaultPort,
		RenderListener:     nil,
		server:             nil,
		extraHTTPListeners: make(map[string]RenderHTTPRequestListener),
	}

	return server
}

// SetHTTPRequestListener sets a extra HTTP request listner.
func (self *Render) SetHTTPRequestListener(path string, listener RenderHTTPRequestListener) error {
	if len(path) <= 0 || listener == nil {
		return fmt.Errorf(errorInvalidHTTPRequestListener, path, listener)
	}

	self.extraHTTPListeners[path] = listener

	return nil
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
	path := httpReq.URL.Path

	switch path {
	case renderDefaultFindRequestPath:
		self.handleFindRequest(httpWriter, httpReq)
		return
	case renderDefaultExpandRequestPath:
		// TODO : Not implemented yet
	case renderDefaultIndexRequestPath:
		self.handleIndexRequest(httpWriter, httpReq)
		return
	case renderDefaultQueryRequestPath:
		self.handleRenderRequest(httpWriter, httpReq)
		return
	}

	httpListener, ok := self.extraHTTPListeners[path]
	if ok {
		httpListener.HTTPRequestReceived(httpReq, httpWriter)
		return
	}

	http.NotFound(httpWriter, httpReq)
}

func (self *Render) responseBadRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	http.Error(httpWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (self *Render) responseInternalServerError(httpWriter http.ResponseWriter, httpReq *http.Request) {
	http.Error(httpWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
