// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"testing"
	"time"
)

type TestCarbon struct {
	*Carbon
	MetricsCount int
}

func NewTestCarbon() *TestCarbon {
	carbon := &TestCarbon{NewCarbon(), 0}
	carbon.Listener = carbon
	return carbon
}

func (self *TestCarbon) MetricRequestReceived(m *Metric, err error) {
	if err != nil {
		return
	}
	self.MetricsCount++
}

func TestNewCarbon(t *testing.T) {
	NewCarbon()
}

func TestCarbonParseMetric(t *testing.T) {
	carbon := NewTestCarbon()

	loopCount := 0
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("path%d", i)
		value := float64(i)
		ts := time.Now().Unix() + int64(i)

		line := fmt.Sprintf("%s %f %d", path, value, ts)

		m, err := carbon.ParseRequestString(line)
		if err != nil {
			t.Error(err)
		}

		if m.Path != path {
			t.Error(fmt.Errorf("%s != %s", m.Path, path))
		}

		if int64(m.Value) != int64(value) {
			t.Error(fmt.Errorf("%f != %f", m.Value, value))
		}

		if m.Timestamp.Unix() != ts {
			t.Error(fmt.Errorf("%d != %d", m.Timestamp.Unix(), ts))
		}

		loopCount++
	}

	if carbon.MetricsCount != loopCount {
		t.Error(fmt.Errorf("%d != %d", carbon.MetricsCount, loopCount))
	}
}
