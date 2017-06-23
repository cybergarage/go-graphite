// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

import (
	"time"
)

// Metric is an instance for Metric of Carbon protocol.
type Metric struct {
	Path      string
	Value     float64
	Timestamp *time.Time
}

// NewMetric returns a new Metric.
func NewMetric() *Metric {
	Metric := &Metric{}
	return Metric
}
