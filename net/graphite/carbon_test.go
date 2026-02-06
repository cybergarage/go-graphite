// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

const (
	carbonTestFeedDataFilename = "carbon_feed_test.dat"
)

type TestCarbon struct {
	*Carbon
	MetricsCount int
}

func NewTestCarbon() *TestCarbon {
	carbon := &TestCarbon{NewCarbon(), 0}
	carbon.SetCarbonListener(carbon)
	return carbon
}

func (carbon *TestCarbon) InsertMetricsRequestReceived(ms []*Metrics, err error) {
	if err != nil {
		return
	}
	for range ms {
		carbon.MetricsCount++
	}
}

func TestNewCarbon(t *testing.T) {
	NewCarbon()
}

func TestCarbonParseMetrics(t *testing.T) {
	carbon := NewTestCarbon()

	loopCount := 0
	for n := range 10 {
		path := fmt.Sprintf("path%d", n)
		value := float64(n)
		ts := time.Now().Unix() + int64(n)

		line := fmt.Sprintf("%s %f %d", path, value, ts)

		ms, err := carbon.FeedPlainTextString(line)
		if err != nil {
			t.Error(err)
		}

		m := ms[0]

		if m.Name != path {
			t.Error(fmt.Errorf("%s != %s", m.Name, path))
		}

		if len(m.DataPoints) != 1 {
			t.Error(fmt.Errorf("%d", len(m.DataPoints)))
		}

		dp := m.DataPoints[0]

		if int64(dp.Value) != int64(value) {
			t.Error(fmt.Errorf("%f != %f", dp.Value, value))
		}

		if dp.Timestamp.Unix() != ts {
			t.Error(fmt.Errorf("%d != %d", dp.Timestamp.Unix(), ts))
		}

		loopCount++
	}

	if carbon.MetricsCount != loopCount {
		t.Error(fmt.Errorf("%d != %d", carbon.MetricsCount, loopCount))
	}
}

func TestCarbonFeed(t *testing.T) {
	feedBytes, err := ioutil.ReadFile(carbonTestFeedDataFilename)
	if err != nil {
		t.Error(err)
		return
	}

	server := newTestServer()

	err = server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	cli := NewClient()

	err = cli.FeedString(string(feedBytes))
	if err != nil {
		t.Error(err)
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestCarbonMultipleFeed(t *testing.T) {
	feedBytes, err := ioutil.ReadFile(carbonTestFeedDataFilename)
	if err != nil {
		t.Error(err)
		return
	}

	server := newTestServer()

	err = server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	cli := NewClient()

	for range 10 {
		err = cli.FeedString(string(feedBytes))
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Millisecond * 100)
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestCarbonMultipleFeedWithKeepConnection(t *testing.T) {
	feedBytes, err := ioutil.ReadFile(carbonTestFeedDataFilename)
	if err != nil {
		t.Error(err)
		return
	}

	serverWaitTimeout := time.Millisecond * 500

	server := newTestServer()
	server.SetConnectionWaitTimeout(serverWaitTimeout)

	err = server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	cli := NewClient()

	conn, err := cli.Open()
	if err != nil {
		t.Error(err)
		server.Stop()
		return
	}

	for range 10 {
		err = cli.FeedStringWithConnection(conn, string(feedBytes))
		if err != nil {
			t.Error(err)
		}
		time.Sleep(serverWaitTimeout * 2)
	}

	err = conn.Close()
	if err != nil {
		t.Error(err)
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}
