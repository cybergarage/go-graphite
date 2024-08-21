// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"time"
)

const (
	MetricStoreEmpty   = "empty"
	MetricStoreSqlite  = "sqlite"
	MetricStoreTsmap   = "tsmap"
	MetricStoreRingmap = "ringmap"
)

const (
	RetentionIntervalFiveMinute = time.Duration(5) * time.Minute
	DefaultRetentionInterval    = RetentionIntervalFiveMinute
	DefaultRetentionPeriod      = time.Duration(60) * time.Minute
	QueryDefaultFromOffset      = -time.Duration(60) * time.Minute
)
