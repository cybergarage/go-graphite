// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package metric provides query interfaces for metric store.
package metric

import (
	"fmt"
	"time"
)

// Metric represents a Foreman metric.
type Metric struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

// NewMetric returns a new metric.
func NewMetric() *Metric {
	m := &Metric{
		Timestamp: time.Now(),
		Value:     0.0,
	}
	return m
}

// NewMetricWithName returns a new metric with the specified name.
func NewMetricWithName(name string) *Metric {
	m := NewMetric()
	m.Name = name
	return m
}

// GetName returns the stored name
func (m *Metric) GetName() string {
	return m.Name
}

// GetValue returns the stored value.
func (m *Metric) GetValue() float64 {
	return m.Value
}

// GetTimestamp returns the stored timestamp
func (m *Metric) GetTimestamp() time.Time {
	return m.Timestamp
}

// String returns a string description of the instance
func (m *Metric) String() string {
	return fmt.Sprintf("%s : %f (%d)", m.Name, m.Value, m.Timestamp.Unix())
}
