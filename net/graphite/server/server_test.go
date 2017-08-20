// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"testing"

	"fmt"
	"time"

	"github.com/cybergarage/go-graphite/net/graphite"
	"github.com/cybergarage/go-graphite/net/graphite/client"
)

type TestServer struct {
	*Server
	MetricsCount int
}

func NewTestServer() *TestServer {
	server := &TestServer{NewServer(), 0}
	server.CarbonListener = server
	return server
}

func (self *TestServer) MetricRequestReceived(m *graphite.Metric, err error) {
	if err != nil {
		return
	}
	//fmt.Printf("MetricRequestReceived = %v\n", m)
	self.MetricsCount++
}

func TestNewServer(t *testing.T) {
	NewServer()
}

func TestServerThread(t *testing.T) {
	server := NewTestServer()

	err := server.Start()
	if err != nil {
		t.Error(err)
	}

	client := client.NewClient()

	loopCount := 0
	for n := 0; n < 10; n++ {
		m := graphite.NewMetric()
		m.Path = fmt.Sprintf("path%d", n)
		m.Value = float64(n)
		m.Timestamp = time.Now()

		err = client.PostMetric(m)
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
