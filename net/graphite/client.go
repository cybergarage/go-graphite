// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	// DefaultHost is the default host for Carbon and Render servers
	DefaultHost string = "localhost"
	// DefaultTimeoutSecond is the default request timeout for Carbon and Render servers
	DefaultTimeoutSecond = 60
)

const (
	errorPostMetric              = "Couldn't write metric [%d] : %v"
	errorFindMetricsStatusCode   = "Bad status code (%d) : %s"
	errorGetAllMetricsStatusCode = "Bad status code (%d)"
)

// Client is an instance for Graphite protocols.
type Client struct {
	Host       string
	CarbonPort int
	RenderPort int
	Timeout    time.Duration
	conn       net.Conn
}

// NewClient returns a new Client.
func NewClient() *Client {
	client := &Client{
		Host:       DefaultHost,
		CarbonPort: DefaultCarbonPort,
		RenderPort: DefaultRenderPort,
		Timeout:    (time.Second * DefaultTimeoutSecond),
	}

	return client
}

// SetHost sets a target host.
func (self *Client) SetHost(host string) {
	self.Host = host
}

// GetHost returns a target host.
func (self *Client) GetHost() string {
	return self.Host
}

// SetCarbonPort sets a target Carbon port.
func (self *Client) SetCarbonPort(port int) {
	self.CarbonPort = port
}

// GetCarbonPort returns a target Carbon port.
func (self *Client) GetCarbonPort() int {
	return self.CarbonPort
}

// SetRenderPort sets a target Carbon port.
func (self *Client) SetRenderPort(port int) {
	self.RenderPort = port
}

// GetRenderPort returns a target Carbon port.
func (self *Client) GetRenderPort() int {
	return self.RenderPort
}

// SetTimeout sets a timeout for the request.
func (self *Client) SetTimeout(d time.Duration) {
	self.Timeout = d
}

// GetTimeout return  the timeout for the request.
func (self *Client) GetTimeout() time.Duration {
	return self.Timeout
}

// Open connects to the specified host.
func (self *Client) Open() (net.Conn, error) {
	addr := net.JoinHostPort(self.Host, strconv.Itoa(self.CarbonPort))
	dialer := net.Dialer{Timeout: self.Timeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// FeedString posts a specified string to Carbon.
func (self *Client) FeedString(m string) error {
	conn, err := self.Open()
	if err != nil {
		return err
	}

	nWrote, err := fmt.Fprintf(conn, "%s", m)
	if err != nil {
		return err
	}
	if nWrote <= 0 {
		return fmt.Errorf(errorPostMetric, nWrote, m)
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	return nil
}

// FeedMetrics posts all metric datapoints to Carbon.
func (self *Client) FeedMetrics(m *Metrics) error {
	for n, _ := range m.DataPoints {
		err := self.feedMetricsDataPoint(m, n)
		if err != nil {
			return err
		}
	}
	return nil
}

// feedMetricsDataPoint posts a specified metric to Carbon.
func (self *Client) feedMetricsDataPoint(m *Metrics, n int) error {
	dpData, err := m.DataPointPlainTextString(n)
	if err != nil {
		return err
	}

	return self.FeedString(dpData)
}

// FindMetrics searches the specified metrics.
// Graphite - The Metrics API
// https://graphite-api.readthedocs.io/en/latest/api.html#the-metrics-api
func (self *Client) FindMetrics(q *Query) ([]*Metrics, error) {
	url, err := q.FindMetricsURL(self.Host, self.RenderPort)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: self.Timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorFindMetricsStatusCode, resp.StatusCode, q.Target)
	}

	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonMetrics renderFindMetricJSONResponse
	err = json.Unmarshal(jsonBytes, &jsonMetrics)
	if err != nil {
		return nil, err
	}

	ms := make([]*Metrics, 0)
	for _, jsonMetric := range jsonMetrics.Metrics {
		m := NewMetrics()
		m.SetName(jsonMetric.Path)
		ms = append(ms, m)
	}

	return ms, nil
}

// GetAllMetrics returns all metrics.
// Graphite - The Metrics API
// https://graphite-api.readthedocs.io/en/latest/api.html#the-metrics-api
func (self *Client) GetAllMetrics() ([]*Metrics, error) {
	url := fmt.Sprintf("http://%s:%d%s",
		self.Host,
		self.RenderPort,
		renderDefaultIndexRequestPath)

	client := http.Client{
		Timeout: self.Timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorGetAllMetricsStatusCode, resp.StatusCode)
	}

	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonMetrics renderMetricIndexJSONResponse
	err = json.Unmarshal(jsonBytes, &jsonMetrics)
	if err != nil {
		return nil, err
	}

	ms := make([]*Metrics, 0)
	for _, name := range jsonMetrics {
		m := NewMetrics()
		m.SetName(name)
		ms = append(ms, m)
	}

	return ms, nil
}

// QueryRender queries with the specified parameters to Render.
// Graphite - The Render API
// https://graphite-api.readthedocs.io/en/latest/api.html#the-render-api-render
func (self *Client) QueryRender(q *Query) ([]*Metrics, error) {
	// FIXME : Support other formats
	q.Format = QueryFormatTypeCSV

	url, err := q.RenderURLString(self.Host, self.RenderPort)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: self.Timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var lastErr error
	ms := make([]*Metrics, 0)

	reader := bufio.NewReader(resp.Body)
	for {
		row, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			lastErr = err
			break
		}

		m := NewMetrics()
		err = m.ParseRenderCSV(string(row))
		if err != nil {
			lastErr = err
			continue
		}
		ms = append(ms, m)
	}

	return ms, lastErr
}
