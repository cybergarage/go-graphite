// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package client provides interfaces for Graphite protocols.
package client

import (
	"fmt"
	"net"

	"github.com/cybergarage/go-graphite/net/graphite"
)

const (
	// CarbonDefaultPort is the default port number for Carbon Server
	CarbonDefaultPort int = 2003
	// CarbonDefaultHost is the default host for Carbon Server
	CarbonDefaultHost string = "localhost"
)

// Client is an instance for Graphite protocols.
type Client struct {
	Host       string
	CarbonPort int
}

// NewClient returns a new Client.
func NewClient() *Client {
	client := &Client{CarbonDefaultHost, CarbonDefaultPort}
	return client
}

// PostMetric returns a metrics of the specified bytes.
func (self *Client) PostMetric(m *graphite.Metric) error {
	addr := fmt.Sprintf("%s:%d", self.Host, self.CarbonPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	nWrote, err := fmt.Fprintf(conn, "%s", m.GoString())
	if err != nil {
		return err
	}
	if nWrote <= 0 {
		return fmt.Errorf("Couldn't write metric [%d] : %v", nWrote, m)
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	return nil
}
