// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"io/ioutil"
	"net/http"
	"testing"

	"fmt"
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

func NewTestServer() *TestServer {
	server := &TestServer{
		NewServer(),
		0,
		0,
	}

	server.CarbonListener = server

	return server
}

func (self *TestServer) InsertMetricsRequestReceived(m *Metrics, err error) {
	if err != nil {
		return
	}
	//fmt.Printf("InsertMetricsRequestReceived = %v\n", m)
	self.MetricsCount++
}

func (self *TestServer) HTTPRequestReceived(r *http.Request, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	self.HTTPRequestCount++
}

func TestNewServer(t *testing.T) {
	NewTestServer()
}

func TestServerQuery(t *testing.T) {
	server := NewTestServer()

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

		err = cli.PostMetrics(m)
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
	server := NewTestServer()
	server.SetHTTPRequestListener(testServerHTTPRequestPath, server)

	err := server.Start()
	if err != nil {
		t.Error(err)
	}

	loopCount := 0
	for n := 0; n < 10; n++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d%s", server.Render.Port, testServerHTTPRequestPath))
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
