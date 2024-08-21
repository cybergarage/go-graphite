// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package metric provides query interfaces for metric store.
package metric

import (
	"fmt"
	"time"
)

// DataPoint represents a Foreman DataPoint.
type DataPoint struct {
	Value     float64
	Timestamp time.Time
}

// NewDataPoint returns a new DataPoint.
func NewDataPoint() *DataPoint {
	dp := &DataPoint{}
	return dp
}

// String returns a string description of the instance
func (dp *DataPoint) String() string {
	return fmt.Sprintf("[%d] %f", dp.Timestamp.Unix(), dp.Value)
}
