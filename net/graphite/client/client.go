// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package client provides interfaces for Graphite protocols.
package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/cybergarage/go-graphite/net/graphite"
)

const (
	// DefaultHost is the default host for Carbon and Render servers
	DefaultHost string = "localhost"
)

// Client is an instance for Graphite protocols.
type Client struct {
	Host       string
	CarbonPort int
	RenderPort int
}

// NewClient returns a new Client.
func NewClient() *Client {
	client := &Client{DefaultHost, graphite.CarbonDefaultPort, graphite.RenderDefaultPort}
	return client
}

// PostMetric posts the specified metric to Carbon.
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

// PostQuery queries with the specified parameters to Render.
func (self *Client) PostQuery(query *graphite.Query) ([]*graphite.Metric, error) {
	// FIXME : Support other formats
	query.Format = graphite.QueryFormatTypeCSV

	url, err := query.URLString(self.Host, self.RenderPort)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var metrics []*graphite.Metric

	reader := bufio.NewReader(resp.Body)
	for {
		row, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return metrics, err
		}

		metric := graphite.NewMetric()
		err = metric.ParseRenderCSV(string(row))
		if err != nil {
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
