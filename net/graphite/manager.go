// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"net"
	"time"
)

const (
	errorManagerNotRunning = "Manager is not running"
)

// A Manager represents a multicast server manager.
type Manager struct {
	*Config

	httpListeners  map[string]RenderHTTPRequestListener
	CarbonListener CarbonListener
	RenderListener RenderRequestListener

	Servers []*Server
}

// NewManager returns a new Manager.
func NewManager() *Manager {
	mgr := &Manager{
		Config: NewDefaultConfig(),

		httpListeners:  map[string]RenderHTTPRequestListener{},
		CarbonListener: nil,
		RenderListener: nil,

		Servers: make([]*Server, 0),
	}
	return mgr
}

// SetHTTPRequestListener sets a extra HTTP request listner.
func (mgr *Manager) SetHTTPRequestListener(path string, listener RenderHTTPRequestListener) error {
	if len(path) == 0 || listener == nil {
		return fmt.Errorf(errorInvalidHTTPRequestListener, path, listener)
	}
	mgr.httpListeners[path] = listener

	for _, server := range mgr.Servers {
		server.SetHTTPRequestListeners(mgr.httpListeners)
	}

	return nil
}

// SetCarbonListener sets a default listener.
func (mgr *Manager) SetCarbonListener(l CarbonListener) error {
	mgr.CarbonListener = l

	for _, server := range mgr.Servers {
		server.SetCarbonListener(l)
	}

	return nil
}

// SetRenderListener sets a default listener.
func (mgr *Manager) SetRenderListener(l RenderRequestListener) error {
	mgr.RenderListener = l

	for _, server := range mgr.Servers {
		server.SetRenderListener(l)
	}

	return nil
}

// GetBoundAddress returns a listen address.
func (mgr *Manager) GetBoundAddress() (string, error) {
	// FIXME : Return an appropriate address instead of addrs[0]
	addrs, err := mgr.GetBoundAddresses()
	if err != nil {
		return "", err
	}

	if 0 < len(addrs) {
		return addrs[0], nil
	}

	return "", fmt.Errorf(errorManagerNotRunning)
}

// GetBoundAddresses returns the listen addresses.
func (mgr *Manager) GetBoundAddresses() ([]string, error) {
	if !mgr.IsRunning() {
		return nil, fmt.Errorf(errorManagerNotRunning)
	}

	boundAddrs := make([]string, 0)

	if mgr.IsEachInterfaceBindingEnabled() {
		for _, server := range mgr.Servers {
			boundAddrs = append(boundAddrs, server.GetBoundAddress())
		}
	} else {
		addrs, err := GetAvailableAddresses()
		if err == nil {
			boundAddrs = append(boundAddrs, addrs...)
		}
	}

	return boundAddrs, nil
}

// GetBoundInterfaces returns the listen interfaces.
func (mgr *Manager) GetBoundInterfaces() ([]*net.Interface, error) {
	if !mgr.IsRunning() {
		return nil, fmt.Errorf(errorManagerNotRunning)
	}

	boundIfs := make([]*net.Interface, 0)

	if mgr.IsEachInterfaceBindingEnabled() {
		for _, server := range mgr.Servers {
			boundIfs = append(boundIfs, server.GetBoundInterface())
		}
	} else {
		ifis, err := GetAvailableInterfaces()
		if err == nil {
			boundIfs = append(boundIfs, ifis...)
		}
	}

	return boundIfs, nil
}

// StartWithInterface starts this server on the specified interface.
func (mgr *Manager) StartWithInterface(ifi *net.Interface) (*Server, error) {
	server := NewServer()
	server.SetConfig(mgr.Config)
	server.SetHTTPRequestListeners(mgr.httpListeners)
	server.SetCarbonListener(mgr.CarbonListener)
	server.SetRenderListener(mgr.RenderListener)

	startupError := fmt.Errorf(errorManagerNotRunning)
	for n := 0; n <= mgr.BindingRetryCount; n++ {
		startupError = server.Start()
		if (startupError == nil) || (mgr.BindingRetryCount <= n) {
			break
		}
		time.Sleep(time.Second * 1)
	}

	if startupError != nil {
		return nil, startupError
	}

	server.SetBoundInterface(ifi)

	// Set bound interface addrss
	if ifi != nil {
		addrs, err := GetInterfaceAddress(ifi)
		if err != nil {
			return nil, err
		}
		server.SetBoundAddress(addrs)
	} else {
		addrs, err := GetAvailableAddresses()
		if err != nil {
			return nil, err
		}
		// FIXME : Set more appropriate addrss instead of addrs[0]
		if 0 < len(addrs) {
			server.SetBoundAddress(addrs[0])
		}
	}

	mgr.Servers = append(mgr.Servers, server)

	return server, nil
}

// Start starts servers on the all avairable interfaces.
func (mgr *Manager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	shouldBindEachInterfaces := mgr.IsEachInterfaceBindingEnabled()
	if mgr.IsAutoInterfaceBindingEnabled() {
		shouldBindEachInterfaces = len(ifis) <= 1
	}

	if shouldBindEachInterfaces {
		for _, ifi := range ifis {
			_, err = mgr.StartWithInterface(ifi)
			if err != nil {
				break
			}
		}
	} else {
		_, err = mgr.StartWithInterface(nil)
	}

	if err != nil {
		mgr.Stop()
	}

	return err
}

// Stop stops this server.
func (mgr *Manager) Stop() error {
	var lastErr error
	for _, server := range mgr.Servers {
		err := server.Stop()
		if err != nil {
			lastErr = err
		}
	}
	mgr.Servers = make([]*Server, 0)
	return lastErr
}

// Stop stops this server.
func (mgr *Manager) getAppropriateServerForInterface(ifi *net.Interface) (*Server, error) {
	if len(mgr.Servers) == 0 {
		return nil, fmt.Errorf(errorManagerNotRunning)
	}

	for _, server := range mgr.Servers {
		if server == nil {
			continue
		}
		if server.GetBoundInterface() == ifi {
			return server, nil
		}
	}

	return mgr.Servers[0], nil
}

// IsRunning returns true whether the local servers are running, otherwise false.
func (mgr *Manager) IsRunning() bool {
	return len(mgr.Servers) != 0
}
