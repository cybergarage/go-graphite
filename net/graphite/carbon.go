// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides interfaces for Graphite protocols.
package graphite

import (
	"fmt"
	"io/ioutil"
	"net"
)

const (
	// CarbonDefaultPort is the default port number for Carbon Server
	CarbonDefaultPort int = 2003
)

// PlaintextRequestListener represents a listener for plain text protocol of Carbon.
type PlaintextRequestListener interface {
	MetricRequestReceived(*Metric, error)
}

// CarbonListener represents a listener for all requests of Carbon.
type CarbonListener interface {
	PlaintextRequestListener
}

// Carbon is an instance for Carbon protocols.
type Carbon struct {
	Port           int
	CarbonListener CarbonListener
	tcpListener    net.Listener
}

// NewCarbon returns a new Carbon.
func NewCarbon() *Carbon {
	carbon := &Carbon{Port: CarbonDefaultPort}
	return carbon
}

// ParseRequestString returns a metrics of the specified context.
func (self *Carbon) ParseRequestString(context string) (*Metric, error) {
	m := NewMetric()
	err := m.ParsePlainText(context)

	if err != nil {
		m = nil
	}

	if self.CarbonListener != nil {
		self.CarbonListener.MetricRequestReceived(m, err)
	}

	return m, err
}

// ParseRequestBytes returns a metrics of the specified bytes.
func (self *Carbon) ParseRequestBytes(bytes []byte) (*Metric, error) {
	return self.ParseRequestString(string(bytes))
}

// Start starts the Carbon server.
func (self *Carbon) Start() error {
	err := self.Stop()
	if err != nil {
		return err
	}

	err = self.open()
	if err != nil {
		return err
	}

	go self.serve()

	return nil
}

// Stop stops the Carbon server.
func (self *Carbon) Stop() error {
	err := self.close()
	if err != nil {
		return err
	}

	return nil
}

// open opens a socket for the Carbon server.
func (self *Carbon) open() error {
	var err error
	addr := fmt.Sprintf(":%d", self.Port)
	self.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return nil
}

// close closes a socket for the Carbon server.
func (self *Carbon) close() error {
	if self.tcpListener != nil {
		err := self.tcpListener.Close()
		if err != nil {
			return err
		}
	}

	self.tcpListener = nil

	return nil
}

// serve handles client requests.
func (self *Carbon) serve() error {
	defer self.close()

	l := self.tcpListener
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		reqBytes, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}

		self.ParseRequestBytes(reqBytes)
	}

	return nil
}
