// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package metric provides query interfaces for metric store.
package metric

import (
	"fmt"
	"testing"
	"time"
)

const (
	testMetricsTestMetricFormat             = "m%03d"
	testMetricsTestMetricCount              = 10
	testMetricsTestMetricsDatapointCount    = 10
	testMetricsTestAllMetricsDatapointCount = testMetricsTestMetricCount * testMetricsTestMetricsDatapointCount
)

func TestNewResultSet(t *testing.T) {
	rs := NewResultSet()

	dpsCount := rs.GetMetricsCount()
	if dpsCount != 0 {
		t.Error(fmt.Errorf("DataPoints is found : %d", dpsCount))
	}

	ms := rs.GetFirstMetrics()
	if ms != nil {
		t.Error(fmt.Errorf("DataPoints is not nil"))
		return
	}
}

func TestResultSetAddMetrics(t *testing.T) {
	rs := NewResultSet()

	for i := 0; i < testMetricsTestMetricCount; i++ {
		ms := NewMetricsWithSize(testMetricsTestMetricsDatapointCount)
		ms.Name = fmt.Sprintf(testMetricsTestMetricFormat, i)
		for j := 0; j < testMetricsTestMetricsDatapointCount; j++ {
			dp := NewDataPoint()
			dp.Timestamp = time.Now()
			dp.Value = float64(j)
			ms.Values[j] = dp
		}
		rs.AddMetrics(ms)
	}

	msCount := rs.GetMetricsCount()
	if msCount != testMetricsTestMetricCount {
		t.Error(fmt.Errorf("%d != %d", msCount, testMetricsTestMetricCount))
	}

	ms := rs.GetFirstMetrics()
	for ms != nil {
		dpCount := len(ms.Values)
		if dpCount != testMetricsTestMetricsDatapointCount {
			t.Error(fmt.Errorf("%d != %d", dpCount, testMetricsTestMetricsDatapointCount))
		}
		ms = rs.GetNextMetrics()
	}

}

func TestResultSetAddSameMetrics(t *testing.T) {
	rs := NewResultSet()

	for i := 0; i < testMetricsTestMetricCount; i++ {
		ms := NewMetricsWithSize(testMetricsTestMetricsDatapointCount)
		ms.Name = fmt.Sprintf(testMetricsTestMetricFormat, 0)
		for j := 0; j < testMetricsTestMetricsDatapointCount; j++ {
			dp := NewDataPoint()
			dp.Timestamp = time.Now()
			dp.Value = float64(j)
			ms.Values[j] = dp
		}
		rs.AddMetrics(ms)
	}

	msCount := rs.GetMetricsCount()
	if msCount != 1 {
		t.Error(fmt.Errorf("%d != %d", msCount, 1))
	}

	ms := rs.GetFirstMetrics()
	if ms == nil {
		t.Error(fmt.Errorf("%d != %d", 0, testMetricsTestAllMetricsDatapointCount))

	}

	dpCount := len(ms.Values)
	if dpCount != testMetricsTestAllMetricsDatapointCount {
		t.Error(fmt.Errorf("%d != %d", dpCount, testMetricsTestAllMetricsDatapointCount))
	}

}
