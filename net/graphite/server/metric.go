// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

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

// ParsePlainText parses the specified line string of the following plain text protocol.
// Feeding In Your Data — Graphite 0.10.0 documentation
// http://graphite.readthedocs.io/en/latest/feeding-carbon.html
func (self *Metric) ParsePlainText(line string) error {
	var err error

	strs := strings.Split(line, " ")
	if len(strs) < 3 {
		return errors.New(fmt.Sprintf(metricParseError, line))
	}

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
