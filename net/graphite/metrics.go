// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	metricParseError                = "Invalid request : %s"
	metricsRenderCSVTimestampFormat = "20060102 15:04:05"
)

// Metrics is an instance for metrics of Carbon protocol.
type Metrics struct {
	Name       string
	DataPoints []*DataPoint
}

// NewMetrics returns a new metrics.
func NewMetrics() *Metrics {
	m := &Metrics{
		DataPoints: NewDataPoints(0),
	}
	return m
}

// NewMetricsWithPlainLine parses the specified line and returns the new metrics.
func NewMetricsWithPlainLine(line string) (*Metrics, error) {
	m := NewMetrics()
	err := m.ParsePlainLine(line)
	return m, err
}

// NewMetricsWithPlainText parses the specified data and returns the new metrics.
func NewMetricsWithPlainText(text string) ([]*Metrics, error) {
	var firstErr error
	lines := strings.Split(text, carbonPlainTextLineSep)
	ms := make([]*Metrics, 0)
	for _, line := range lines {
		if len(line) <= 0 {
			continue
		}
		m, err := NewMetricsWithPlainLine(line)
		if err != nil {
			if firstErr != nil {
				firstErr = err
			}
			continue
		}
		ms = append(ms, m)
	}

	return ms, firstErr
}

// SetName sets a name to the metrics.
func (self *Metrics) SetName(name string) {
	self.Name = name
}

// GetName returns the metrics name.
func (self *Metrics) GetName() string {
	return self.Name
}

// GetDataPointCount returns a count of the datapoints
func (self *Metrics) GetDataPointCount() int {
	return len(self.DataPoints)
}

// AddDataPoint add a new datapoint to the metrics
func (self *Metrics) AddDataPoint(dp *DataPoint) error {
	self.DataPoints = append(self.DataPoints, dp)
	return nil
}

// GetDataPoint retur a datapoint of the specified index.
func (self *Metrics) GetDataPoint(n int) (*DataPoint, error) {
	if (n < 0) || (len(self.DataPoints) <= n) {
		return nil, fmt.Errorf(errorInvalidRangeIndex, n, len(self.DataPoints))
	}
	return self.DataPoints[n], nil
}

// SortDataPoints sorts the current datapoints
func (self *Metrics) SortDataPoints() error {
	sort.Sort(DataPoints(self.DataPoints))
	return nil
}

// ParsePlainLine parses the specified line string of the following plain text protocol.
// Feeding In Your Data — Graphite 0.10.0 documentation
// http://graphite.readthedocs.io/en/latest/feeding-carbon.html
func (self *Metrics) ParsePlainLine(line string) error {
	strs := strings.Split(strings.Trim(line, carbonPlainTextLineTrim), carbonPlainTextLineFieldSep)
	if len(strs) != 3 {
		return fmt.Errorf(metricParseError, line)
	}

	var err error

	self.Name = strs[0]

	value, err := strconv.ParseFloat(strs[1], 64)
	if err != nil {
		return err
	}

	ts, err := TimeStringToTime(strs[2])
	if err != nil {
		return err
	}

	dp := NewDataPoint()
	dp.Value = value
	dp.Timestamp = *ts

	err = self.AddDataPoint(dp)
	if err != nil {
		return err
	}

	return nil
}

// ParseRenderCSV parses the specified line string of the following Render CSV protocol.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Metrics) ParseRenderCSV(line string) error {
	strs := strings.Split(line, ",")
	if len(strs) != 3 {
		return fmt.Errorf(metricParseError, line)
	}

	var err error

	self.Name = strings.TrimSpace(strs[0])

	ts, err := time.Parse(metricsRenderCSVTimestampFormat, strings.TrimSpace(strs[1]))
	if err != nil {
		return err
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(strs[2]), 64)
	if err != nil {
		return err
	}

	dp := NewDataPoint()
	dp.Value = value
	dp.Timestamp = ts

	err = self.AddDataPoint(dp)
	if err != nil {
		return err
	}

	return nil
}

// DataPointPlainTextString returns a string representation datapoint for the plaintext protocol.
func (self *Metrics) DataPointPlainTextString(n int) (string, error) {
	if len(self.DataPoints) < n {
		return "", fmt.Errorf(errorInvalidRangeIndex, n, len(self.DataPoints))
	}
	dp := self.DataPoints[n]
	return fmt.Sprintf("%s %s", self.Name, dp.PlainTextString()), nil
}
