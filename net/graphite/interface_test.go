// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"testing"
)

func TestGetAvailableInterfaces(t *testing.T) {
	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
	}
	if len(ifs) <= 0 {
		t.Errorf("available interface is not found")
	}
}

func TestGetAvailableAddresses(t *testing.T) {
	addrs, err := GetAvailableAddresses()
	if err != nil {
		t.Error(err)
	}
	if len(addrs) <= 0 {
		t.Errorf("available address is not found")
	}
}
