// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"net/http"
	"strings"
)

// handleFindRequest handles requests for Metrics API.
// The Render URL API
// http://readthedocs.io/en/latest/render_api.html
func (self *Render) handleFindRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	query := NewQuery()
	err := query.ParseHTTPRequest(httpReq)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	if self.RenderListener == nil {
		self.responseInternalServerError(httpWriter, httpReq)
		return
	}

	metrics, err := self.RenderListener.FindMetricsRequestReceived(query, nil)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	self.responseFindMetrics(httpWriter, httpReq, query, metrics)
}

func (self *Render) responseFindMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	switch query.Format {
	case QueryFormatTypeCompleter: // TODO : Not implemented yet
		self.responseBadRequest(httpWriter, httpReq)
		return
	case QueryFormatTypeTreeJSON:
		self.responseFindJSONMetrics(httpWriter, httpReq, query, metrics)
		return
	default:
		self.responseFindJSONMetrics(httpWriter, httpReq, query, metrics)
		return
	}

	self.responseBadRequest(httpWriter, httpReq)
}

func (self *Render) responseFindJSONMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeJSON)
	httpWriter.Header().Set(httpHeaderAccessControlAllowOrigin, httpHeaderAccessControlAllowOriginAll)
	httpWriter.WriteHeader(http.StatusOK)

	httpWriter.Write([]byte("{\"metrics\": [\n"))

	mCount := len(metrics)

	for i, m := range metrics {

		paths := strings.Split(m.Name, renderMetricsDelim)
		pathCount := len(paths)
		if pathCount <= 0 {
			continue
		}

		// Start bracket
		httpWriter.Write([]byte("{\n"))

		// Leaf infomation (static)
		httpWriter.Write([]byte("\"is_leaf\": 1,\n"))

		// Name and path

		name := paths[pathCount-1]
		prefix := ""
		if 1 < pathCount {
			prefixes := paths[:pathCount-1]
			prefix = strings.Join(prefixes, renderMetricsDelim)
		}

		httpWriter.Write([]byte(fmt.Sprintf("\"name\": \"%s\",\n", name)))
		httpWriter.Write([]byte(fmt.Sprintf("\"path\": \"%s\"\n", prefix)))

		// End bracket
		if i < (mCount - 1) {
			httpWriter.Write([]byte("},\n"))
		} else {
			httpWriter.Write([]byte("}\n"))
		}
	}

	httpWriter.Write([]byte("]}\n"))
}

// handleIndexRequest handles requests for Metrics API.
// The Render URL API
// http://readthedocs.io/en/latest/render_api.html
func (self *Render) handleIndexRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	if self.RenderListener == nil {
		self.responseInternalServerError(httpWriter, httpReq)
		return
	}

	// Find all metrics

	query := NewQuery()
	query.Target = renderMetricsAsterisk
	metrics, err := self.RenderListener.FindMetricsRequestReceived(query, nil)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	// Response by JSON array

	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeJSON)
	httpWriter.Header().Set(httpHeaderAccessControlAllowOrigin, httpHeaderAccessControlAllowOriginAll)
	httpWriter.WriteHeader(http.StatusOK)

	httpWriter.Write([]byte("[\n"))

	mCount := len(metrics)
	for i, m := range metrics {
		if i < (mCount - 1) {
			httpWriter.Write([]byte(fmt.Sprintf("\"%s\",\n", m.Name)))
		} else {
			httpWriter.Write([]byte(fmt.Sprintf("\"%s\"\n", m.Name)))
		}
	}

	httpWriter.Write([]byte("]\n"))
}
