// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package client provides interfaces for Graphite protocols.
package client

const (
	// CarbonDefaultPort is the default port number for Carbon Server
	CarbonDefaultPort int = 2003
	// CarbonDefaultHost is the default host for Carbon Server
	CarbonDefaultHost string = "localhost"
)

// Client is an instance for Graphite protocols.
type Client struct {
}

// NewClient returns a new Client.
func NewClient() *Client {
	client := &Client{}
	return client
}
