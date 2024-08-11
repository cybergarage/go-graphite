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
	metricParseError                = "invalid request : %s"
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
	lines := strings.FieldsFunc(text,
		func(r rune) bool {
			return r == '\n' || r == '\r'
		})
	ms := make([]*Metrics, 0)
	for _, line := range lines {
		if len(line) == 0 {
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
func (m *Metrics) SetName(name string) {
	m.Name = name
}

// GetName returns the metrics name.
func (m *Metrics) GetName() string {
	return m.Name
}

// GetDataPointCount returns a count of the datapoints.
func (m *Metrics) GetDataPointCount() int {
	return len(m.DataPoints)
}

// AddDataPoint add a new datapoint to the metrics.
func (m *Metrics) AddDataPoint(dp *DataPoint) error {
	m.DataPoints = append(m.DataPoints, dp)
	return nil
}

// GetDataPoint retur a datapoint of the specified index.
func (m *Metrics) GetDataPoint(n int) (*DataPoint, error) {
	if (n < 0) || (len(m.DataPoints) <= n) {
		return nil, fmt.Errorf(errorInvalidRangeIndex, n, len(m.DataPoints))
	}
	return m.DataPoints[n], nil
}

// SortDataPoints sorts the current datapoints.
func (m *Metrics) SortDataPoints() error {
	sort.Sort(DataPoints(m.DataPoints))
	return nil
}

// ParsePlainLine parses the specified line string of the following plain text protocol.
// Feeding In Your Data — Graphite 0.10.0 documentation
// http://graphite.readthedocs.io/en/latest/feeding-carbon.html
func (m *Metrics) ParsePlainLine(line string) error {
	strs := strings.Split(strings.Trim(line, carbonPlainTextLineTrim), carbonPlainTextLineFieldSep)
	if len(strs) != 3 {
		return fmt.Errorf(metricParseError, line)
	}

	var err error

	m.Name = strs[0]

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

	err = m.AddDataPoint(dp)
	if err != nil {
		return err
	}

	return nil
}

// ParseRenderCSV parses the specified line string of the following Render CSV protocol.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (m *Metrics) ParseRenderCSV(line string) error {
	strs := strings.Split(line, ",")
	if len(strs) != 3 {
		return fmt.Errorf(metricParseError, line)
	}

	var err error

	m.Name = strings.TrimSpace(strs[0])

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

	err = m.AddDataPoint(dp)
	if err != nil {
		return err
	}

	return nil
}

// DataPointPlainTextString returns a string representation datapoint for the plaintext protocol.
func (m *Metrics) DataPointPlainTextString(n int) (string, error) {
	if len(m.DataPoints) < n {
		return "", fmt.Errorf(errorInvalidRangeIndex, n, len(m.DataPoints))
	}
	dp := m.DataPoints[n]
	return fmt.Sprintf("%s %s", m.Name, dp.PlainTextString()), nil
}
