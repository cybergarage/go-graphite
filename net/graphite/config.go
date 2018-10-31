// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import "reflect"

// Config represents a cofiguration for extended specifications.
type Config struct {
	EachInterfaceBindingEnabled bool
	AutoInterfaceBindingEnabled bool
	Addr                        string
	CarbonPort                  int
	RenderPort                  int
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		EachInterfaceBindingEnabled: true,
		AutoInterfaceBindingEnabled: true,
		Addr:                        "",
		CarbonPort:                  DefaultCarbonPort,
		RenderPort:                  DefaultRenderPort,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *Config) SetConfig(newConfig *Config) {
	conf.EachInterfaceBindingEnabled = newConfig.EachInterfaceBindingEnabled
	conf.AutoInterfaceBindingEnabled = newConfig.AutoInterfaceBindingEnabled
}

// SetAddress sets a bind address.
func (conf *Config) SetAddress(addr string) {
	conf.Addr = addr
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
func (conf *Config) SetCarbonPort(port int) error {
	conf.CarbonPort = port
	return nil
}

// GetCarbonPort returns a bind port for Carbon.
func (conf *Config) GetCarbonPort() int {
	return conf.CarbonPort
}

// SetRenderPort sets a bind port for Render.
func (conf *Config) SetRenderPort(port int) error {
	conf.RenderPort = port
	return nil
}

// GetRenderPort returns a bind port for Render.
func (conf *Config) GetRenderPort() int {
	return conf.RenderPort
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *Config) Equals(otherConf *Config) bool {
	return reflect.DeepEqual(conf, otherConf)
}
