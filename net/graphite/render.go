// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"net/http"
	"time"
)

const (
	// RenderDefaultPort is the default port number for Render
	RenderDefaultPort int = 8080
	// RenderDefaultPath is the default path for Render
	RenderDefaultPath string = "/render"
)

const (
	httpHeaderContentType = "Content-Type"
)

// RenderRequestListener represents a listener for Render protocol.
type RenderRequestListener interface {
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
	case RenderDefaultPath:
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

// handleRenderRequest handles Render requests.
// The Render URL API
// http://readthedocs.io/en/latest/render_api.html
func (self *Render) handleRenderRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	query := NewQuery()
	err := query.Parse(httpReq.URL)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	if self.RenderListener == nil {
		self.responseInternalServerError(httpWriter, httpReq)
		return
	}

	metrics, err := self.RenderListener.QueryMetricsRequestReceived(query, nil)
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

func (self *Render) responseQueryMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	switch query.Format {
	case QueryFormatTypeRaw:
		self.responseQueryRawMetrics(httpWriter, httpReq, query, metrics)
		return
	case QueryFormatTypeCSV:
		self.responseQueryCSVMetrics(httpWriter, httpReq, query, metrics)
		return
	case QueryFormatTypeJSON:
		self.responseQueryJSONMetrics(httpWriter, httpReq, query, metrics)
		return
	}

	self.responseBadRequest(httpWriter, httpReq)
}

func (self *Render) responseQueryRawMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeRaw)
	httpWriter.WriteHeader(http.StatusOK)

	for _, m := range metrics {
		err := m.SortDataPoints()
		if err != nil {
			continue
		}
		dpCount := m.GetDataPointCount()
		if dpCount <= 0 {
			continue
		}

		var from, until time.Time
		var step int64

		switch dpCount {
		case 0:
			continue
		case 1:
			firstDp, err := m.GetDataPoint(0)
			if err != nil {
				continue
			}
			from = firstDp.Timestamp
			until = from
			step = 0
		default:
			firstDp, err := m.GetDataPoint(0)
			if err != nil {
				continue
			}
			from = firstDp.Timestamp

			lastDp, err := m.GetDataPoint((dpCount - 1))
			if err != nil {
				continue
			}
			until = lastDp.Timestamp

			secondDp, err := m.GetDataPoint(1)
			if err != nil {
				continue
			}
			// FIXME : Step is calculated only using the first few data points.
			step = secondDp.Timestamp.Unix() - firstDp.Timestamp.Unix()
		}

		msg := fmt.Sprintf("%s,%d,%d,%d|", m.Name, from.Unix(), until.Unix(), step)
		httpWriter.Write([]byte(msg))

		for n, dp := range m.DataPoints {
			var value string
			switch n {
			case (dpCount - 1):
				value = fmt.Sprintf("%f", dp.Value)
			default:
				value = fmt.Sprintf("%f,", dp.Value)
			}
			httpWriter.Write([]byte(value))
		}

		httpWriter.Write([]byte("\n"))
	}
}

func (self *Render) responseQueryCSVMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeCSV)
	httpWriter.WriteHeader(http.StatusOK)
	for _, m := range metrics {
		for _, dp := range m.DataPoints {
			mRow := fmt.Sprintf("%s,%s\n", m.Name, dp.RenderCSVString())
			httpWriter.Write([]byte(mRow))
		}
	}
}

func (self *Render) responseQueryJSONMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeJSON)
	httpWriter.WriteHeader(http.StatusOK)

	httpWriter.Write([]byte("[\n"))

	mCount := len(metrics)
	for i, m := range metrics {
		httpWriter.Write([]byte("{\n"))

		// Output the target name
		httpWriter.Write([]byte(fmt.Sprintf("\"target\": \"%s\",\n", m.Name)))

		// Output the datapoint array
		dpCount := m.GetDataPointCount()
		httpWriter.Write([]byte("\"datapoints\": [\n"))
		for j, dp := range m.DataPoints {
			httpWriter.Write([]byte(fmt.Sprintf("%f,%d\n", dp.Value, dp.UnixTimestamp())))
			if j < (dpCount - 1) {
				httpWriter.Write([]byte("],\n"))
			} else {
				httpWriter.Write([]byte("]\n"))
			}
		}
		httpWriter.Write([]byte("]\n"))

		if i < (mCount - 1) {
			httpWriter.Write([]byte("},\n"))
		} else {
			httpWriter.Write([]byte("}\n"))
		}
	}

	httpWriter.Write([]byte("]\n"))
}
