// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package metric provides query interfaces for metric store.
package metric

import (
	"testing"

	"github.com/cybergarage/foreman-go/foreman/node"
)

func TestNewMetric(t *testing.T) {
	NewMetric()
}

func TestNewRegexMetric(t *testing.T) {
	node := node.NewBaseNode()

	m := NewMetric()
	_, err := m.GetRegexMetricForNode(node)
	if err != nil {
		t.Error(err)
	}
}
