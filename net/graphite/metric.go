// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graphite provides interfaces for Graphite protocols.
package graphite

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	metricParseError = "Could not parse %s"
)

// Metric is an instance for Metric of Carbon protocol.
type Metric struct {
	Path      string
	Value     float64
	Timestamp time.Time
}

// NewMetric returns a new Metric.
func NewMetric() *Metric {
	Metric := &Metric{}
	return Metric
}

// Parse parses the specified context.
func (self *Metric) Parse(line string) error {
	strs := strings.Split(line, " ")
	if len(strs) == 3 {
		return self.parsePlainText(strs)
	}

	return errors.New(fmt.Sprintf(metricParseError, line))
}

// parsePlainText parses the specified line string of the following plain text protocol.
// Feeding In Your Data — Graphite 0.10.0 documentation
// http://graphite.readthedocs.io/en/latest/feeding-carbon.html
func (self *Metric) parsePlainText(strs []string) error {
	var err error

	self.Path = strs[0]

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

// GoString returns a string representation value.
func (self *Metric) GoString() string {
	return fmt.Sprintf("%s %f %d", self.Path, self.Value, self.Timestamp.Unix())
}
