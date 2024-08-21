// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

// ResultSet represents an abstract interface of metric store for Foreman.
type ResultSet interface {
	AddMetrics(*Metrics) error
	GetMetricsCount() int
	GetFirstMetrics() *Metrics
	GetNextMetrics() *Metrics
}

// goResultSet represents a result set.
type goResultSet struct {
	ResultSet
	metrics []*Metrics
	// For GetMetricsByName() to access by O(1)
	metricsMap map[string]*Metrics
	// For GetFirstMetrics() and GetNextMetrics()
	iteratorPos int
}

// NewResultSet returns a new result set.
func NewResultSet() ResultSet {
	rs := &goResultSet{
		metrics:     make([]*Metrics, 0),
		metricsMap:  make(map[string]*Metrics),
		iteratorPos: 0,
	}
	return rs
}

// AddMetrics adds a new data point.
func (rs *goResultSet) AddMetrics(ms *Metrics) error {
	existMs := rs.GetMetricsByName(ms.Name)
	if existMs != nil {
		return existMs.AddDataPoints(ms.Values)
	}
	rs.metrics = append(rs.metrics, ms)
	rs.metricsMap[ms.Name] = ms
	return nil
}

// GetMetricsCount returns a number of the data points.
func (rs *goResultSet) GetMetricsCount() int {
	return len(rs.metrics)
}

// GetMetricsByName returns a metrics by name.
func (rs *goResultSet) GetMetricsByName(name string) *Metrics {
	ms, ok := rs.metricsMap[name]
	if !ok {
		return nil
	}
	return ms
}

// GetFirstMetrics returns a first data points.
func (rs *goResultSet) GetFirstMetrics() *Metrics {
	rs.iteratorPos = 0
	return rs.GetNextMetrics()
}

// GetNextMetrics returns a first data points.
func (rs *goResultSet) GetNextMetrics() *Metrics {
	if len(rs.metrics) <= rs.iteratorPos {
		return nil
	}
	dp := rs.metrics[rs.iteratorPos]
	rs.iteratorPos++
	return dp
}
