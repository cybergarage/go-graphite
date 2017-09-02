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
	metricParseError       = "Could not parse %s"
	metricsTimestampFormat = "20060102 15:04:05"
)

// Metric is an instance for Metric of Carbon protocol.
type Metric struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

// NewMetric returns a new Metric.
func NewMetric() *Metric {
	m := &Metric{}
	return m
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

	self.Value, err = strconv.ParseFloat(strs[1], 64)
	if err != nil {
		return err
	}

	var unixTime int64
	unixTime, err = strconv.ParseInt(strs[2], 10, 64)
	if err != nil {
		return err
	}
	self.Timestamp = time.Unix(unixTime, 0)

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

	self.Timestamp, err = time.Parse(metricsTimestampFormat, strs[1])
	if err != nil {
		return err
	}

	self.Value, err = strconv.ParseFloat(strs[1], 64)
	if err != nil {
		return err
	}

	return nil
}

// GoString returns a string representation value.
func (self *Metric) GoString() string {
	return fmt.Sprintf("%s %f %d", self.Name, self.Value, self.Timestamp.Unix())
}

// TimestampString returns a string for the Render API format.
func (self *Metric) TimestampString() string {
	return self.Timestamp.Format(metricsTimestampFormat)
}
