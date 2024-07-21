// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"net/http"
)

const (
	httpHeaderContentType                 = "Content-Type"
	httpHeaderAccessControlAllowOrigin    = "Access-Control-Allow-Origin"
	httpHeaderAccessControlAllowOriginAll = "*"
)

const (
	renderDefaultFindRequestPath   string = "/metrics/find"
	renderDefaultExpandRequestPath string = "/metrics/expand"
	renderDefaultIndexRequestPath  string = "/metrics/index.json"
	renderDefaultQueryRequestPath  string = "/render"
	renderMetricsDelim             string = "."
	renderMetricsAsterisk          string = "*"
)

// ServeHTTP handles HTTP requests.
func (render *Render) ServeHTTP(httpWriter http.ResponseWriter, httpReq *http.Request) {
	path := httpReq.URL.Path

	switch path {
	case renderDefaultFindRequestPath:
		render.handleFindRequest(httpWriter, httpReq)
		return
	case renderDefaultExpandRequestPath:
		// TODO : Not implemented yet
	case renderDefaultIndexRequestPath:
		render.handleIndexRequest(httpWriter, httpReq)
		return
	case renderDefaultQueryRequestPath:
		render.handleRenderRequest(httpWriter, httpReq)
		return
	}

	httpListener, ok := render.extraHTTPListeners[path]
	if ok {
		httpListener.HTTPRequestReceived(httpReq, httpWriter)
		return
	}

	http.NotFound(httpWriter, httpReq)
}

func (render *Render) responseBadRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	httpWriter.Header().Set(httpHeaderAccessControlAllowOrigin, httpHeaderAccessControlAllowOriginAll)
	http.Error(httpWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (render *Render) responseInternalServerError(httpWriter http.ResponseWriter, httpReq *http.Request) {
	httpWriter.Header().Set(httpHeaderAccessControlAllowOrigin, httpHeaderAccessControlAllowOriginAll)
	http.Error(httpWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
