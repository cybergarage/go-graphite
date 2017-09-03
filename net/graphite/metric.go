// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graphite provides interfaces for Graphite protocols.
package graphite

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	metricParseError                = "Could not parse %s"
	metricsRenderCSVTimestampFormat = "20060102 15:04:05"
)

// Metric is an instance for Metric of Carbon protocol.
type Metric struct {
	Name       string
	DataPoints []*DataPoint
}

// NewMetric returns a new Metric.
func NewMetric() *Metric {
	m := &Metric{
		DataPoints: make([]*DataPoint, 0),
	}
	return m
}

// AddDataPoint add a new datapoint
func (self *Metric) AddDataPoint(dp *DataPoint) error {
	self.DataPoints = append(self.DataPoints, dp)
	return nil
}

// ParsePlainText parses the specified line string of the following plain text protocol.
// Feeding In Your Data — Graphite 0.10.0 documentation
// http://graphite.readthedocs.io/en/latest/feeding-carbon.html
func (self *Metric) ParsePlainText(line string) error {
	strs := strings.Split(line, " ")
	if len(strs) != 3 {
		return fmt.Errorf(metricParseError, line)
	}

	var err error

	self.Name = strs[0]

	value, err := strconv.ParseFloat(strs[1], 64)
	if err != nil {
		return err
	}

	var unixTime int64
	unixTime, err = strconv.ParseInt(strs[2], 10, 64)
	if err != nil {
		return err
	}

	dp := NewDataPoint()
	dp.Value = value
	dp.Timestamp = time.Unix(unixTime, 0)

	err = self.AddDataPoint(dp)
	if err != nil {
		return err
	}

	return nil
}

// ParseRenderCSV parses the specified line string of the following Render CSV protocol.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Metric) ParseRenderCSV(line string) error {
	strs := strings.Split(line, ", ")
	if len(strs) != 3 {
		return fmt.Errorf(metricParseError, line)
	}

	var err error

	self.Name = strs[0]

	ts, err := time.Parse(metricsRenderCSVTimestampFormat, strs[1])
	if err != nil {
		return err
	}

	value, err := strconv.ParseFloat(strs[1], 64)
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
func (self *Metric) DataPointPlainTextString(n int) (string, error) {
	if len(self.DataPoints) < n {
		return "", fmt.Errorf(errorInvalidIndex, n, len(self.DataPoints))
	}
	dp := self.DataPoints[n]
	return fmt.Sprintf("%s %s", self.Name, dp.PlainTextString()), nil
}
