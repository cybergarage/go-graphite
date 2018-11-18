// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

const (
	// DefaultCarbonPort is the default port number for Carbon Server
	DefaultCarbonPort int = 2003
)

// PlaintextRequestListener represents a listener for plain text protocol of Carbon.
// See : Feeding In Your Data (http://graphite.readthedocs.io/en/latest/feeding-carbon.html)
type PlaintextRequestListener interface {
	InsertMetricsRequestReceived(*Metrics, error)
}

// CarbonListener represents a listener for all requests of Carbon.
type CarbonListener interface {
	PlaintextRequestListener
}

// Carbon is an instance for Carbon protocols.
type Carbon struct {
	addr           string
	port           int
	carbonListener CarbonListener
	tcpListener    net.Listener
}

// NewCarbon returns a new Carbon.
func NewCarbon() *Carbon {
	carbon := &Carbon{
		addr:           "",
		port:           DefaultCarbonPort,
		carbonListener: nil,
		tcpListener:    nil,
	}
	return carbon
}

// SetAddress sets a bind address to the server.
func (carbon *Carbon) SetAddress(addr string) {
	carbon.addr = addr
}

// GetAddress returns a bound address.
func (carbon *Carbon) GetAddress() string {
	return carbon.addr
}

// SetPort sets a bind porto the server.
func (carbon *Carbon) SetPort(port int) {
	carbon.port = port
}

// GetPort returns a bound port.
func (carbon *Carbon) GetPort() int {
	return carbon.port
}

// SetCarbonListener sets a default listener.
func (carbon *Carbon) SetCarbonListener(listener CarbonListener) {
	carbon.carbonListener = listener
}

// parseRequestLine parses the specified metrics request.
func (carbon *Carbon) parseRequestLine(lineString string) (*Metrics, error) {
	fmt.Printf("%s\n", lineString)

	m := NewMetrics()
	err := m.ParsePlainText(lineString)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		m = nil
	}

	if carbon.carbonListener != nil {
		carbon.carbonListener.InsertMetricsRequestReceived(m, err)
	}

	return m, err
}

// ParseRequestString returns a metrics of the specified context.
func (carbon *Carbon) ParseRequestString(context string) ([]*Metrics, error) {
	lines := strings.Split(context, "\n")
	ms := make([]*Metrics, len(lines))
	for n, line := range lines {
		if len(line) <= 0 {
			continue
		}
		m, err := carbon.parseRequestLine(line)
		if err != nil {
			return ms, err
		}
		ms[n] = m
	}

	return ms, nil
}

// ParseRequestBytes returns a metrics of the specified bytes.
func (carbon *Carbon) ParseRequestBytes(bytes []byte) ([]*Metrics, error) {
	return carbon.ParseRequestString(string(bytes))
}

// Start starts the Carbon server.
func (carbon *Carbon) Start() error {
	err := carbon.Stop()
	if err != nil {
		return err
	}

	err = carbon.open()
	if err != nil {
		return err
	}

	go carbon.serve()

	return nil
}

// Stop stops the Carbon server.
func (carbon *Carbon) Stop() error {
	err := carbon.close()
	if err != nil {
		return err
	}

	return nil
}

// open opens a socket for the Carbon server.
func (carbon *Carbon) open() error {
	var err error
	addr := net.JoinHostPort(carbon.addr, strconv.Itoa(carbon.port))
	carbon.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

// close closes a socket for the Carbon server.
func (carbon *Carbon) close() error {
	if carbon.tcpListener != nil {
		err := carbon.tcpListener.Close()
		if err != nil {
			return err
		}
	}

	carbon.tcpListener = nil

	return nil
}

// serve handles client requests.
func (carbon *Carbon) serve() error {
	defer carbon.close()

	l := carbon.tcpListener
	for {
		if l == nil {
			break
		}
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		reqBytes, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}

		carbon.ParseRequestBytes(reqBytes)
	}

	return nil
}
