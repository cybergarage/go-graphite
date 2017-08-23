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
	httpHeaderContentType = "Content-Type"
)

// RenderRequestListener represents a listener for Render protocol.
type RenderRequestListener interface {
	QueryRequestReceived(*graphite.Query, error) ([]*graphite.Metric, error)
}

// RenderListener represents a listener for all requests of Render.
type RenderListener interface {
	RenderRequestListener
}

// Render is an instance for Graphite render protocols.
type Render struct {
	Port           int
	RenderListener RenderRequestListener
	server         *http.Server
}

// NewRender returns a new Render.
func NewRender() *Render {
	server := &Render{
		Port:   graphite.RenderDefaultPort,
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
	case graphite.RenderDefaultPath:
		self.handleRenderRequest(httpWriter, httpReq)
		return
	}

	http.NotFound(httpWriter, httpReq)
}

// handleRenderRequest handles Render requests.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Render) handleRenderRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	query := graphite.NewQuery()
	err := query.Parse(httpReq.URL)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	if self.RenderListener == nil {
		self.responseInternalServerError(httpWriter, httpReq)
		return
	}

	metrics, err := self.RenderListener.QueryRequestReceived(query, nil)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	self.responseQueryMetrics(httpWriter, httpReq, query, metrics)
}

func (self *Render) responseBadRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	http.Error(httpWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (self *Render) responseInternalServerError(httpWriter http.ResponseWriter, httpReq *http.Request) {
	http.Error(httpWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (self *Render) responseQueryMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *graphite.Query, metrics []*graphite.Metric) {
	switch query.Format {
	case graphite.QueryFormatTypeCSV:
		self.responseQueryCSVMetrics(httpWriter, httpReq, query, metrics)
		return
	case graphite.QueryFormatTypeJSON:
		self.responseQueryJSONMetrics(httpWriter, httpReq, query, metrics)
		return
	}

	self.responseBadRequest(httpWriter, httpReq)
}

func (self *Render) responseQueryCSVMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *graphite.Query, metrics []*graphite.Metric) {
	httpWriter.Header().Set(httpHeaderContentType, graphite.QueryContentTypeCSV)
	httpWriter.WriteHeader(http.StatusOK)
	for _, m := range metrics {
		mRow := fmt.Sprintf("%s,%s,%f\n", m.Path, m.Value)
		httpWriter.Write([]byte(mRow))
	}
}

func (self *Render) responseQueryJSONMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *graphite.Query, metrics []*graphite.Metric) {
	httpWriter.Header().Set(httpHeaderContentType, graphite.QueryContentTypeJSON)
	httpWriter.WriteHeader(http.StatusOK)
	// FIXME : Not implemented yet
}
