// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

type renderFindMetricsJSONMetrics struct {
	IsLeaf int    `json:"is_leaf"`
	Name   string `json:"name"`
	Path   string `json:"path"`
}

type renderFindMetricJSONResponse struct {
	Metrics []renderFindMetricsJSONMetrics `json:"metrics"`
}

type renderMetricIndexJSONMetrics struct {
	Name string `json:"metrics"`
}

type renderMetricIndexJSONResponse []string
