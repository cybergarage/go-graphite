// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import "net/http"

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
