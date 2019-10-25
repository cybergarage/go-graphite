// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"reflect"
	"time"
)

const (
	// DefaultBindingRetryCount is a default retry count when the server can't bind the specified port.
	DefaultBindingRetryCount = 0
	// DefaultConnectionTimeout is a default timeout for Render and Carbon server.
	DefaultConnectionTimeout = time.Second * 60
	// DefaultConnectionWaitTimeout is a default wait timeout for Render and Carbon server.
	DefaultConnectionWaitTimeout = time.Second * 10
)

// Config represents a cofiguration for extended specifications.
type Config struct {
	EachInterfaceBindingEnabled bool
	AutoInterfaceBindingEnabled bool
	BindingRetryCount           int
	Addr                        string
	CarbonPort                  int
	RenderPort                  int
	ConnectionTimeout           time.Duration
	ConnectionWaitTimeout       time.Duration
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		EachInterfaceBindingEnabled: true,
		AutoInterfaceBindingEnabled: true,
		Addr:                        "",
		CarbonPort:                  DefaultCarbonPort,
		RenderPort:                  DefaultRenderPort,
		BindingRetryCount:           DefaultBindingRetryCount,
		ConnectionTimeout:           DefaultConnectionTimeout,
		ConnectionWaitTimeout:       DefaultConnectionWaitTimeout,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *Config) SetConfig(newConfig *Config) {
	conf.EachInterfaceBindingEnabled = newConfig.EachInterfaceBindingEnabled
	conf.AutoInterfaceBindingEnabled = newConfig.AutoInterfaceBindingEnabled
	conf.CarbonPort = newConfig.CarbonPort
	conf.RenderPort = newConfig.RenderPort
	conf.BindingRetryCount = newConfig.BindingRetryCount
}

// SetAddress sets a configuration address.
func (conf *Config) SetAddress(addr string) {
	conf.Addr = addr
}

// GetAddress returns a configuration address.
func (conf *Config) GetAddress() string {
	return conf.Addr
}

// SetEachInterfaceBindingEnabled sets a flag for binding functions.
func (conf *Config) SetEachInterfaceBindingEnabled(flag bool) {
	conf.EachInterfaceBindingEnabled = flag
}

// IsEachInterfaceBindingEnabled returns true whether the binding functions is enabled, otherwise false.
func (conf *Config) IsEachInterfaceBindingEnabled() bool {
	return conf.EachInterfaceBindingEnabled
}

// SetAutoInterfaceBindingEnabled sets a flag for the auto interface binding.
func (conf *Config) SetAutoInterfaceBindingEnabled(flag bool) {
	conf.AutoInterfaceBindingEnabled = flag
}

// IsAutoInterfaceBindingEnabled returns true whether the the auto interface binding is enabled, otherwise false.
func (conf *Config) IsAutoInterfaceBindingEnabled() bool {
	return conf.AutoInterfaceBindingEnabled
}

// SetCarbonPort sets a bind port for Carbon.
func (conf *Config) SetCarbonPort(port int) {
	conf.CarbonPort = port
}

// GetCarbonPort returns a bind port for Carbon.
func (conf *Config) GetCarbonPort() int {
	return conf.CarbonPort
}

// SetRenderPort sets a bind port for Render.
func (conf *Config) SetRenderPort(port int) {
	conf.RenderPort = port
}

// GetRenderPort returns a bind port for Render.
func (conf *Config) GetRenderPort() int {
	return conf.RenderPort
}

// SetBindingRetryCount sets a bind retry count.
func (conf *Config) SetBindingRetryCount(n int) {
	conf.BindingRetryCount = n
}

// GetBindingRetryCount returns a bind retry count.
func (conf *Config) GetBindingRetryCount() int {
	return conf.BindingRetryCount
}

// SetConnectionTimeout sets the connection timeout for the carbon and the render server.
func (conf *Config) SetConnectionTimeout(d time.Duration) {
	conf.ConnectionTimeout = d
}

// GetConnectionTimeout return the connection timeout of he carbon and the render server.
func (conf *Config) GetConnectionTimeout() time.Duration {
	return conf.ConnectionTimeout
}

// SetConnectionWaitTimeout sets the connection timeout for the carbon and the render server.
func (conf *Config) SetConnectionWaitTimeout(d time.Duration) {
	conf.ConnectionWaitTimeout = d
}

// GetConnectionWaitTimeout return the connection timeout of he carbon and the render server.
func (conf *Config) GetConnectionWaitTimeout() time.Duration {
	return conf.ConnectionWaitTimeout
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *Config) Equals(otherConf *Config) bool {
	return reflect.DeepEqual(conf, otherConf)
}
