// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

type renderFindMetricJSONResponseMetrics struct {
	IsLeaf int    `json:"is_leaf"`
	Name   string `json:"name"`
	Path   string `json:"path"`
}

type renderFindMetricJSONResponse struct {
	Metrics []renderFindMetricJSONResponseMetrics `json:"metrics"`
}

type renderMetricIndexJSONResponse struct {
	Metrics []string
}
