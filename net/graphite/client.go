// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
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
	client := &Client{DefaultHost, DefaultCarbonPort, DefaultRenderPort}
	return client
}

// PostMetrics posts all metric datapoints to Carbon.
func (self *Client) PostMetrics(m *Metrics) error {
	for n, _ := range m.DataPoints {
		err := self.postMetricsDataPoint(m, n)
		if err != nil {
			return err
		}
	}
	return nil
}

// postMetricsDataPoint posts a specified metric to Carbon.
func (self *Client) postMetricsDataPoint(m *Metrics, n int) error {
	dpData, err := m.DataPointPlainTextString(n)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(self.Host, strconv.Itoa(self.CarbonPort))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	nWrote, err := fmt.Fprintf(conn, "%s", dpData)
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
func (self *Client) PostQuery(query *Query) ([]*Metrics, error) {
	// FIXME : Support other formats
	query.Format = QueryFormatTypeCSV

	url, err := query.RenderURLString(self.Host, self.RenderPort)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var metrics []*Metrics

	reader := bufio.NewReader(resp.Body)
	for {
		row, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return metrics, err
		}

		metric := NewMetrics()
		err = metric.ParseRenderCSV(string(row))
		if err != nil {
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
