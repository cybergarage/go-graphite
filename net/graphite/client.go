// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	// DefaultHost is the default host for Carbon and Render servers.
	DefaultHost string = "localhost"
	// DefaultTimeoutSecond is the default request timeout for Carbon and Render servers.
	DefaultTimeoutSecond = 60
)

const (
	errorPostMetric              = "couldn't write metric [%d] : %v"
	errorFindMetricsStatusCode   = "bad status code (%d) : %s"
	errorGetAllMetricsStatusCode = "bad status code (%d)"
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
func (client *Client) SetHost(host string) {
	client.Host = host
}

// GetHost returns a target host.
func (client *Client) GetHost() string {
	return client.Host
}

// SetCarbonPort sets a target Carbon port.
func (client *Client) SetCarbonPort(port int) {
	client.CarbonPort = port
}

// GetCarbonPort returns a target Carbon port.
func (client *Client) GetCarbonPort() int {
	return client.CarbonPort
}

// SetRenderPort sets a target Carbon port.
func (client *Client) SetRenderPort(port int) {
	client.RenderPort = port
}

// GetRenderPort returns a target Carbon port.
func (client *Client) GetRenderPort() int {
	return client.RenderPort
}

// SetTimeout sets a timeout for the request.
func (client *Client) SetTimeout(d time.Duration) {
	client.Timeout = d
}

// GetTimeout return  the timeout for the request.
func (client *Client) GetTimeout() time.Duration {
	return client.Timeout
}

// Open connects to the specified host.
func (client *Client) Open() (net.Conn, error) {
	addr := net.JoinHostPort(client.Host, strconv.Itoa(client.CarbonPort))
	dialer := net.Dialer{Timeout: client.Timeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// FeedString posts a specified string to Carbon.
func (client *Client) FeedString(m string) error {
	conn, err := client.Open()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = client.FeedStringWithConnection(conn, m)
	if err != nil {
		return err
	}

	return nil
}

// FeedStringWithConnection posts a specified string to the specified connection.
func (client *Client) FeedStringWithConnection(conn net.Conn, m string) error {
	nWrote, err := fmt.Fprintf(conn, "%s", m)
	if err != nil {
		return err
	}
	if nWrote <= 0 {
		return fmt.Errorf(errorPostMetric, nWrote, m)
	}

	return nil
}

// FeedMetrics posts all metric datapoints to Carbon.
func (client *Client) FeedMetrics(m *Metrics) error {
	for n := range m.DataPoints {
		err := client.feedMetricsDataPoint(m, n)
		if err != nil {
			return err
		}
	}
	return nil
}

// feedMetricsDataPoint posts a specified metric to Carbon.
func (client *Client) feedMetricsDataPoint(m *Metrics, n int) error {
	dpData, err := m.DataPointPlainTextString(n)
	if err != nil {
		return err
	}

	return client.FeedString(dpData)
}

// FindMetrics searches the specified metrics.
// Graphite - The Metrics API
// https://graphite-api.readthedocs.io/en/latest/api.html#the-metrics-api
func (client *Client) FindMetrics(q *Query) ([]*Metrics, error) {
	url, err := q.FindMetricsURL(client.Host, client.RenderPort)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Timeout: client.Timeout,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorFindMetricsStatusCode, resp.StatusCode, q.Target)
	}

	jsonBytes, err := io.ReadAll(resp.Body)
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
func (client *Client) GetAllMetrics() ([]*Metrics, error) {
	hostPort := net.JoinHostPort(client.Host, strconv.Itoa(client.RenderPort))
	url := fmt.Sprintf("http://%s%s", hostPort, renderDefaultIndexRequestPath)

	httpClient := http.Client{
		Timeout: client.Timeout,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorGetAllMetricsStatusCode, resp.StatusCode)
	}

	jsonBytes, err := io.ReadAll(resp.Body)
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
func (client *Client) QueryRender(q *Query) ([]*Metrics, error) {
	// FIXME : Support other formats
	q.Format = QueryFormatTypeCSV

	url, err := q.RenderURLString(client.Host, client.RenderPort)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Timeout: client.Timeout,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var lastErr error
	ms := make([]*Metrics, 0)

	reader := bufio.NewReader(resp.Body)
	for {
		row, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
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
