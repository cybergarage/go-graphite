// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"testing"
)

func testManagerBinding(t *testing.T, conf *Config) {
	mgr := NewManager()
	mgr.SetConfig(conf)

	err := mgr.Start()
	if err != nil {
		t.Error(err)
	}

	err = mgr.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestManagerWithDefaultConfig(t *testing.T) {
	conf := NewDefaultConfig()
	testManagerBinding(t, conf)
}

func TestManagerWithBindEachInterfaces(t *testing.T) {
	conf := NewDefaultConfig()
	conf.SetAutoInterfaceBindingEnabled(false)
	conf.SetEachInterfaceBindingEnabled(true)
	testManagerBinding(t, conf)
}

func TestManagerWithNoBindInterfaces(t *testing.T) {
	conf := NewDefaultConfig()
	conf.SetAutoInterfaceBindingEnabled(false)
	conf.SetEachInterfaceBindingEnabled(false)
	testManagerBinding(t, conf)
}
