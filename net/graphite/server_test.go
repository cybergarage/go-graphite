// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

const (
	testServerHTTPRequestPath = "/hello"
)

type TestServer struct {
	*Server
	MetricsCount     int
	HTTPRequestCount int
}

func newTestServer() *TestServer {
	server := &TestServer{
		NewServer(),
		0,
		0,
	}

	server.SetCarbonListener(server)

	return server
}

func (server *TestServer) InsertMetricsRequestReceived(ms []*Metrics, err error) {
	if err != nil {
		return
	}

	for range ms {
		// fmt.Printf("InsertMetricsRequestReceived = %v\n", m)
		server.MetricsCount++
	}
}

func (server *TestServer) HTTPRequestReceived(r *http.Request, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	server.HTTPRequestCount++
}

func TestNewServer(t *testing.T) {
	newTestServer()
}

func TestServerQuery(t *testing.T) {
	server := newTestServer()

	err := server.Start()
	if err != nil {
		t.Error(err)
	}

	cli := NewClient()

	loopCount := 0
	for n := 0; n < 10; n++ {
		m := NewMetrics()
		m.Name = fmt.Sprintf("path%d", n)

		dp := NewDataPoint()
		dp.Value = float64(n)
		dp.Timestamp = time.Now()
		m.AddDataPoint(dp)

		err = cli.FeedMetrics(m)
		if err != nil {
			t.Error(err)
		}

		loopCount++
	}

	time.Sleep(1 * time.Second)

	if server.MetricsCount != loopCount {
		t.Error(fmt.Errorf("%d != %d", server.MetricsCount, loopCount))
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestServerHTTPRequest(t *testing.T) {
	server := newTestServer()
	server.SetHTTPRequestListener(testServerHTTPRequestPath, server)

	err := server.Start()
	if err != nil {
		t.Error(err)
	}

	loopCount := 0
	for n := 0; n < 10; n++ {
		url := fmt.Sprintf("http://%s%s",
			net.JoinHostPort(server.GetBoundAddress(), strconv.Itoa(server.Render.GetPort())), testServerHTTPRequestPath)
		resp, err := http.Get(url)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		loopCount++
	}

	if server.HTTPRequestCount != loopCount {
		t.Error(fmt.Errorf("%d != %d", server.HTTPRequestCount, loopCount))
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}
