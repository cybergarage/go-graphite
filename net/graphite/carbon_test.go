// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

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
	carbon.CarbonListener = carbon
	return carbon
}

func (self *TestCarbon) InsertMetricsRequestReceived(m *Metrics, err error) {
	if err != nil {
		return
	}
	self.MetricsCount++
}

func TestNewCarbon(t *testing.T) {
	NewCarbon()
}

func TestCarbonParseMetrics(t *testing.T) {
	carbon := NewTestCarbon()

	loopCount := 0
	for n := 0; n < 10; n++ {
		path := fmt.Sprintf("path%d", n)
		value := float64(n)
		ts := time.Now().Unix() + int64(n)

		line := fmt.Sprintf("%s %f %d", path, value, ts)

		ms, err := carbon.ParseRequestString(line)
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
