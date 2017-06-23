// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

// HTTPServer is an instance for Graphite Web protocols.
type HTTPServer struct {
}

// NewHTTPServer returns a new HTTPServer.
func NewHTTPServer() *HTTPServer {
	server := &HTTPServer{}
	return server
}
