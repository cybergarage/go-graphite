// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"testing"
)

func testManagerBoundAddressWithConfig(t *testing.T, conf *Config) {
	mgr := NewManager()
	mgr.SetConfig(conf)

	err := mgr.Start()
	if err != nil {
		t.Error(err)
	}

	addrs := mgr.GetBoundAddresses()
	if len(addrs) <= 0 {
		t.Errorf("%d", len(addrs))
	}

	for n, addr := range addrs {
		if !IsIPv4Address(addr) {
			t.Errorf("[%d] : %s", n, addr)
		}
	}
	err = mgr.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestManagerBoundAddressWithDefaultConfig(t *testing.T) {
	conf := NewDefaultConfig()
	testManagerBoundAddressWithConfig(t, conf)
}
