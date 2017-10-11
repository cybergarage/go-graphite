// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const (
	testDataPointsSortSampleCount = 100
)

func testDataPointsCheckOrder(t *testing.T, dps DataPoints) {
	sort.Sort(DataPoints(dps))

	dpsSize := len(dps)
	for n := 0; n < dpsSize; n++ {
		if n == 0 {
			continue
		}
		if !dps[n-1].Timestamp.Before(dps[n].Timestamp) {
			t.Error(fmt.Errorf("[%d]:%d > [%d]:%d ", (n - 1), dps[n-1].Timestamp.Unix(), n, dps[n].Timestamp.Unix()))
		}
	}
}

func TestDataPointsNormalOrderSort(t *testing.T) {
	now := time.Now()

	dps := NewDataPoints(testDataPointsSortSampleCount)
	for n := 0; n < testDataPointsSortSampleCount; n++ {
		dp := NewDataPoint()
		dp.Value = 100
		dp.Timestamp = time.Unix((now.Unix() + int64(n)), 0)
		dps[n] = dp
	}

	testDataPointsCheckOrder(t, dps)
}

func TestDataPointsReverseOrderSort(t *testing.T) {
	now := time.Now()

	dps := NewDataPoints(testDataPointsSortSampleCount)
	for n := 0; n < testDataPointsSortSampleCount; n++ {
		dp := NewDataPoint()
		dp.Value = 100
		dp.Timestamp = time.Unix((now.Unix() - int64(testDataPointsSortSampleCount-n)), 0)
		dps[n] = dp
	}

	testDataPointsCheckOrder(t, dps)
}

func TestDataPointsReverseRandomSort(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	dps := NewDataPoints(testDataPointsSortSampleCount)
	for n := 0; n < testDataPointsSortSampleCount; n++ {
		dp := NewDataPoint()
		dp.Value = 100
		dp.Timestamp = time.Unix((rand.Int63n(now.Unix())), 0)
		dps[n] = dp
	}

	testDataPointsCheckOrder(t, dps)
}
