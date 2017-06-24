// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"errors"
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

func (self *TestCarbon) MetricRequestReceived(*Metric) {
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

		m, err := carbon.ParseRequest(line)
		if err != nil {
			t.Error(err)
		}

		if m.Path != path {
			t.Error(errors.New(fmt.Sprintf("%s != %s", m.Path, path)))
		}

		if int64(m.Value) != int64(value) {
			t.Error(errors.New(fmt.Sprintf("%f != %f", m.Value, value)))
		}

		if m.Timestamp.Unix() != ts {
			t.Error(errors.New(fmt.Sprintf("%d != %d", m.Timestamp.Unix(), ts)))
		}

		loopCount++
	}

	if carbon.MetricsCount != loopCount {
		t.Error(errors.New(fmt.Sprintf("%d != %d", carbon.MetricsCount, loopCount)))
	}
}
