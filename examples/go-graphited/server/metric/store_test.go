// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"fmt"
	"testing"
	"time"
)

const (
	testStoreMetricsCount       = 100
	testStoreMetricsPrefix      = "path"
	testStoreMetricsInterval    = DefaultRetentionInterval
	testStoreMetricsPeriodCount = 10
)

func testStore(t *testing.T, store *Store) {
	store.SetRetentionInterval(testStoreMetricsInterval)

	err := store.Open()
	if err != nil {
		t.Error(err)
	}

	// Setup metrics

	var m [testStoreMetricsCount]*Metric
	for n := 0; n < testStoreMetricsCount; n++ {
		m[n] = NewMetric()
		m[n].Name = fmt.Sprintf("%s%d", testStoreMetricsPrefix, n)
	}

	// Testing range

	now := time.Now()
	diff := now.Unix() % int64(testStoreMetricsInterval.Seconds())
	from := now.Add(-(time.Duration(diff) * time.Second))
	until := from

	// Insert metrics

	for i := 0; i < testStoreMetricsPeriodCount; i++ {
		for j := 0; j < testStoreMetricsCount; j++ {
			m[j].Timestamp = until
			m[j].Value = float64(i * j)
			err = store.AddMetric(m[j])
			if err != nil {
				t.Error(err)
			}
		}
		until = until.Add(testStoreMetricsInterval)
	}
}

func TestNewSQLiteStore(t *testing.T) {
	testStore(t, NewSQLiteStore())
}
