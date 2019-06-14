// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	// DefaultRenderPort is the default port number for Render
	DefaultRenderPort int = 8080
	// DefaultRenderConnectionTimeout is a default timeout for Render.
	DefaultRenderConnectionTimeout time.Duration = DefaultConnectionTimeout
)

// Render is an instance for Graphite render protocols.
type Render struct {
	addr               string
	port               int
	connectionTimeout  time.Duration
	renderListener     RenderRequestListener
	server             *http.Server
	extraHTTPListeners map[string]RenderHTTPRequestListener
}

// NewRender returns a new Render.
func NewRender() *Render {
	server := &Render{
		addr:               "",
		port:               DefaultRenderPort,
		connectionTimeout:  DefaultRenderConnectionTimeout,
		renderListener:     nil,
		server:             nil,
		extraHTTPListeners: make(map[string]RenderHTTPRequestListener),
	}

	return server
}

// SetAddress sets a bind address to the server.
func (render *Render) SetAddress(addr string) {
	render.addr = addr
}

// GetAddress returns a bound address.
func (render *Render) GetAddress() string {
	return render.addr
}

// SetPort sets a bind porto the server.
func (render *Render) SetPort(port int) {
	render.port = port
}

// GetPort returns a bound port.
func (render *Render) GetPort() int {
	return render.port
}

// SetConnectionTimeout sets the connection timeout.
func (render *Render) SetConnectionTimeout(d time.Duration) {
	render.connectionTimeout = d
}

// GetConnectionTimeout return the connection timeout.
func (render *Render) GetConnectionTimeout() time.Duration {
	return render.connectionTimeout
}

// SetRenderListener sets a default listener.
func (render *Render) SetRenderListener(listener RenderRequestListener) {
	render.renderListener = listener
}

// SetHTTPRequestListener sets a extra HTTP request listener.
func (render *Render) SetHTTPRequestListener(path string, listener RenderHTTPRequestListener) error {
	if len(path) <= 0 || listener == nil {
		return fmt.Errorf(errorInvalidHTTPRequestListener, path, listener)
	}

	render.extraHTTPListeners[path] = listener

	return nil
}

// SetHTTPRequestListeners sets a extra HTTP request listeners.
func (render *Render) SetHTTPRequestListeners(listeners map[string]RenderHTTPRequestListener) error {
	var lastError error
	for path, listener := range listeners {
		err := render.SetHTTPRequestListener(path, listener)
		if err != nil {
			lastError = err
		}
	}
	return lastError
}

// Start starts the HTTP server.
func (render *Render) Start() error {
	err := render.Stop()
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(render.addr, strconv.Itoa(render.port))

	render.server = &http.Server{
		Addr:        addr,
		ReadTimeout: render.connectionTimeout,
		Handler:     render,
	}

	c := make(chan error)
	go func() {
		c <- render.server.ListenAndServe()
	}()

	select {
	case err = <-c:
	case <-time.After(time.Millisecond * 500):
		err = nil
	}

	return err
}

// Stop stops the HTTP server.
func (render *Render) Stop() error {
	if render.server == nil {
		return nil
	}

	err := render.server.Close()
	if err != nil {
		return err
	}

	return nil
}
