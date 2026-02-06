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

// SetValue sets a value to the datapoint.
func (dp *DataPoint) SetValue(value float64) {
	dp.Value = value
}

// GetValue returns the value of the datapoint.
func (dp *DataPoint) GetValue() float64 {
	return dp.Value
}

// SetTimestamp sets a timestamp to the datapoint.
func (dp *DataPoint) SetTimestamp(value time.Time) {
	dp.Timestamp = value
}

// GetTimestamp returns the timestamp of the datapoint.
func (dp *DataPoint) GetTimestamp() time.Time {
	return dp.Timestamp
}

// PlainTextString returns a string representation datapoint for the plaintext protocol.
func (dp *DataPoint) PlainTextString() string {
	return fmt.Sprintf("%f %d", dp.Value, dp.UnixTimestamp())
}

// RenderCSVString returns a string representation datapoint for the render CSV format.
func (dp *DataPoint) RenderCSVString() string {
	return fmt.Sprintf("%s,%f", dp.Timestamp.Format(metricsRenderCSVTimestampFormat), dp.Value)
}

// TimestampString returns a string for the Render API format.
func (dp *DataPoint) TimestampString() string {
	return dp.Timestamp.Format(metricsRenderCSVTimestampFormat)
}

// UnixTimestamp returns the timestamp as a unix timestamp.
func (dp *DataPoint) UnixTimestamp() int64 {
	return dp.Timestamp.Unix()
}
