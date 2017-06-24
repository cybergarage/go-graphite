// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

// PlaintextListener represents a listener for plain text protocol of Carbon.
type PlaintextListener interface {
	MetricRequestReceived(*Metric)
}

// CarbonListener represents a listener for all requests of Carbon.
type CarbonListener interface {
	PlaintextListener
}

// Carbon is an instance for Carbon protocols.
type Carbon struct {
	Listener CarbonListener
}

// NewCarbon returns a new Carbon.
func NewCarbon() *Carbon {
	carbon := &Carbon{}
	return carbon
}

// Parse returns a metrics array of the specified context.
// Feeding In Your Data — Graphite 0.10.0 documentation
// http://graphite.readthedocs.io/en/latest/feeding-carbon.html
func (self *Carbon) Parse(context string) (bool, []Metric) {
	return true, nil
}

// Parse returns a metrics array of the specified context.
// Feeding In Your Data — Graphite 0.10.0 documentation
// http://graphite.readthedocs.io/en/latest/feeding-carbon.html
func (self *Carbon) parsePlainText(context string) (bool, *Metric) {
	return true, nil
}
