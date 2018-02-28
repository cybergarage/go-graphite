// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"time"
)

// DataPoint is an instance for datapoint of Render protocol.
type DataPoint struct {
	Value     float64
	Timestamp time.Time
}

// NewDataPoint returns a new Metrics.
func NewDataPoint() *DataPoint {
	p := &DataPoint{}
	return p
}

// PlainTextString returns a string representation datapoint for the plaintext protocol.
func (self *DataPoint) PlainTextString() string {
	return fmt.Sprintf("%f %d", self.Value, self.UnixTimestamp())
}

// RenderCSVString returns a string representation datapoint for the render CSV format.
func (self *DataPoint) RenderCSVString() string {
	return fmt.Sprintf("%s,%f", self.Timestamp.Format(metricsRenderCSVTimestampFormat), self.Value)
}

// TimestampString returns a string for the Render API format.
func (self *DataPoint) TimestampString() string {
	return self.Timestamp.Format(metricsRenderCSVTimestampFormat)
}

// UnixTimestampg returns a Unix timestamp value.
func (self *DataPoint) UnixTimestamp() int64 {
	return self.Timestamp.Unix()
}
