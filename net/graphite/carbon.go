// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"bufio"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

const (
	// DefaultCarbonPort is the default port number for Carbon.
	DefaultCarbonPort int = 2003
	// DefaultCarbonConnectionTimeout is a default timeout for Carbon.
	DefaultCarbonConnectionTimeout time.Duration = DefaultConnectionTimeout
)

const (
	carbonPlainTextLineSep      = "\n"
	carbonPlainTextLineTrim     = "\n\r"
	carbonPlainTextLineFieldSep = " "
)

// PlainTextRequestListener represents a listener for plain text protocol of Carbon.
// See : Feeding In Your Data (http://graphite.readthedocs.io/en/latest/feeding-carbon.html)
type PlainTextRequestListener interface {
	InsertMetricsRequestReceived([]*Metrics, error)
}

// CarbonListener represents a listener for all requests of Carbon.
type CarbonListener interface {
	PlainTextRequestListener
}

// Carbon is an instance for Carbon protocols.
type Carbon struct {
	addr              string
	port              int
	connectionTimeout time.Duration
	carbonListener    CarbonListener
	tcpListener       net.Listener
}

// NewCarbon returns a new Carbon.
func NewCarbon() *Carbon {
	carbon := &Carbon{
		addr:              "",
		port:              DefaultCarbonPort,
		connectionTimeout: DefaultCarbonConnectionTimeout,
		carbonListener:    nil,
		tcpListener:       nil,
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

// SetConnectionTimeout sets the connection timeout.
func (carbon *Carbon) SetConnectionTimeout(d time.Duration) {
	carbon.connectionTimeout = d
}

// GetConnectionTimeout return the connection timeout.
func (carbon *Carbon) GetConnectionTimeout() time.Duration {
	return carbon.connectionTimeout
}

// SetCarbonListener sets a default listener.
func (carbon *Carbon) SetCarbonListener(listener CarbonListener) {
	carbon.carbonListener = listener
}

// FeedPlainTextString returns a metrics of the specified text.
func (carbon *Carbon) FeedPlainTextString(reqString string) ([]*Metrics, error) {
	ms, err := NewMetricsWithPlainText(reqString)

	if carbon.carbonListener != nil {
		carbon.carbonListener.InsertMetricsRequestReceived(ms, err)
	}

	return ms, err
}

// FeedPlainTextBytes returns a metrics of the specified bytes.
func (carbon *Carbon) FeedPlainTextBytes(reqBytes []byte) ([]*Metrics, error) {
	return carbon.FeedPlainTextString(string(reqBytes))
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

		go carbon.receive(conn)
	}

	return nil
}

func (carbon *Carbon) receive(conn net.Conn) error {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(carbon.connectionTimeout))

	reader := bufio.NewReader(conn)
	reqBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	carbon.FeedPlainTextBytes(reqBytes)

	return nil
}
