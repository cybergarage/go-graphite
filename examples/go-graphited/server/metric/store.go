// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package metric provides query interfaces for metric store.
package metric

import (
	"time"
)

// StoreListener represents a listener for metric store.
type StoreListener interface {
	StoreMetricAdded(*Metric) error
}

// Storing represents an abstract interface of metric store
type Storing interface {
	SetStoreListener(StoreListener) error

	SetRetentionInterval(value time.Duration) error
	GetRetentionInterval() (time.Duration, error)

	SetRetentionPeriod(value time.Duration) error
	GetRetentionPeriod() (time.Duration, error)

	Open() error
	Close() error
	Clear() error

	AddMetric(m *Metric) error
	AddMetricWithoutNotification(m *Metric) error
	NotifyMetric(m *Metric) error

	Vacuum() error

	String() string
}

// Store represents an metric store.
type Store struct {
	Storing
}
