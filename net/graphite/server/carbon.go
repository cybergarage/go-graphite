// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package server

import (
	"fmt"
	"net"
)

const (
	// CarbonDefaultPort is the default port number for Carbon Server
	CarbonDefaultPort int = 2003
)

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
	Port     int
	Listener CarbonListener
}

// NewCarbon returns a new Carbon.
func NewCarbon() *Carbon {
	carbon := &Carbon{Port: CarbonDefaultPort}
	return carbon
}

// ParseRequest returns a metrics of the specified context.
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

// Start starts the Carbon server.
func (self *Carbon) Start() error {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", self.Port))
	if err != nil {
		return err
	}

	for {
		_, err := ln.Accept()
		if err != nil {
			return err
		}
	}

	return nil
}

// Stop stops the Carbon server.
func (self *Carbon) Stop() error {
	return nil
}
