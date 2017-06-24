// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

// PlaintextRequestListener represents a listener for plain text protocol of Carbon.
type PlaintextRequestListener interface {
	MetricRequestReceived(*Metric)
}

// CarbonListener represents a listener for all requests of Carbon.
type CarbonListener interface {
	PlaintextRequestListener
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

// Parse returns a metrics of the specified context.
func (self *Carbon) ParseRequest(context string) (*Metric, error) {
	m := NewMetric()
	err := m.Parse(context)
	if err != nil {
		return nil, err
	}

	if self.Listener != nil {
		self.Listener.MetricRequestReceived(m)
	}

	return m, nil
}
