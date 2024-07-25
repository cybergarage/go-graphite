// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

// PlainTextRequestListener represents a listener for plain text protocol of Carbon.
// See : Feeding In Your Data (http://graphite.readthedocs.io/en/latest/feeding-carbon.html)
type PlainTextRequestListener interface {
	InsertMetricsRequestReceived([]*Metrics, error)
}

// CarbonListener represents a listener for all requests of Carbon.
type CarbonListener interface {
	PlainTextRequestListener
}
