// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	NewMetrics()
}

func TestMetricsParsePlainLine(t *testing.T) {
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("path%d", i)
		value := float64(i) * 100
		ts := time.Now().Unix() + int64(i)

		line := fmt.Sprintf("%s %f %d", path, value, ts)

		m := NewMetrics()
		err := m.ParsePlainLine(line)
		if err != nil {
			t.Error(err)
		}

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
	}
}

func TestMetricsParsePlainText(t *testing.T) {
	feedBytes, err := ioutil.ReadFile(carbonTestFeedDataFilename)
	if err != nil {
		t.Error(err)
		return
	}

	feedString := strings.Trim(string(feedBytes), carbonPlainTextLineTrim)

	ms, err := NewMetricsWithPlainText(feedString)
	if err != nil {
		t.Error(err)
		return
	}

	feedLines := strings.Split(feedString, carbonPlainTextLineSep)

	if len(ms) != len(feedLines) {
		t.Errorf("%d != %d", len(ms), len(feedLines))
	}
}
