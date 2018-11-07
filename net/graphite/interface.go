// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	errorNullInterface            = "Null interface"
	errorAvailableAddressNotFound = "Available address not found"
	errorAvailableInterfaceFound  = "Available interface not found"
)

// IsIPv6Address retusn true whether the specified address is a IPv6 address
func IsIPv6Address(addr string) bool {
	if len(addr) <= 0 {
		return false
	}

	if 0 <= strings.Index(addr, ":") {
		return true
	}
	return false
}

// IsIPv4Address retusn true whether the specified address is a IPv4 address
func IsIPv4Address(addr string) bool {
	if len(addr) <= 0 {
		return false
	}

	return !IsIPv6Address(addr)
}

// IsLocalAddress retusn true whether the specified address is a local addresses
func IsLocalAddress(addr string) bool {
	localAddrs := []string{
		"127.0.0.1",
		"::1",
	}

	for _, localAddr := range localAddrs {
		if localAddr == addr {
			return true
		}
	}

	return false
}

// GetInterfaceAddress retuns a IPv4 address of the specivied interface.
func GetInterfaceAddress(ifi *net.Interface) (string, error) {
	if ifi == nil {
		return "", fmt.Errorf(errorNullInterface)
	}

	addrs, err := ifi.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		addrStr := addr.String()
		saddr := strings.Split(addrStr, "/")
		if len(saddr) < 2 {
			continue
		}

		// Disabled IPv6 interface
		if IsIPv6Address(saddr[0]) {
			continue
		}

		return saddr[0], nil
	}

	return "", errors.New(errorAvailableAddressNotFound)
}

// GetAvailableInterfaces retuns all available interfaces in the node.
func GetAvailableInterfaces() ([]*net.Interface, error) {
	useIfs := make([]*net.Interface, 0)

	localIfs, err := net.Interfaces()
	if err != nil {
		return useIfs, err
	}

	for _, localIf := range localIfs {
		if (localIf.Flags & net.FlagLoopback) != 0 {
			continue
		}
		if (localIf.Flags & net.FlagUp) == 0 {
			continue
		}
		if (localIf.Flags & net.FlagMulticast) == 0 {
			continue
		}

		_, addrErr := GetInterfaceAddress(&localIf)
		if addrErr != nil {
			continue
		}

		useIf := localIf
		useIfs = append(useIfs, &useIf)
	}

	if len(useIfs) <= 0 {
		return useIfs, errors.New(errorAvailableInterfaceFound)
	}

	return useIfs, err
}

// HasMultipleAvailableInterfaces retuns true when the system has multiple interfaces, otherwise false.
func HasMultipleAvailableInterfaces() bool {
	ifes, err := GetAvailableInterfaces()
	if err != nil {
		return false
	}

	if len(ifes) <= 1 {
		return false
	}

	return true
}

// GetAvailableAddresses retuns all available IPv4 addresses in the node
func GetAvailableAddresses() ([]string, error) {
	addrs := make([]string, 0)

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return addrs, err
	}

	for _, ifi := range ifis {
		addr, err := GetInterfaceAddress(ifi)
		if err != nil {
			continue
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func getMatchAddressBlockCount(ifAddr string, targetAddr string) int {
	const addrSep = "."
	targetAddrs := strings.Split(targetAddr, addrSep)
	ifAddrs := strings.Split(ifAddr, addrSep)

	if len(targetAddrs) != len(ifAddrs) {
		return -1
	}

	addrSize := len(targetAddrs)
	for n := 0; n < len(targetAddrs); n++ {
		if targetAddrs[n] != ifAddrs[n] {
			return n
		}
	}

	return addrSize
}

// GetAvailableInterfaceForAddr returns an interface of the specified address.
func GetAvailableInterfaceForAddr(fromAddr string) (*net.Interface, error) {
	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return nil, err
	}

	switch len(ifis) {
	case 0:
		return nil, errors.New(errorAvailableInterfaceFound)
	case 1:
		return ifis[0], nil
	}

	ifAddrs := make([]string, len(ifis))
	for n := 0; n < len(ifAddrs); n++ {
		ifAddrs[n], _ = GetInterfaceAddress(ifis[n])
	}

	selIf := ifis[0]
	selIfMatchBlocks := getMatchAddressBlockCount(fromAddr, ifAddrs[0])
	for n := 0; n < len(ifAddrs); n++ {
		matchBlocks := getMatchAddressBlockCount(fromAddr, ifAddrs[n])
		if matchBlocks < selIfMatchBlocks {
			continue
		}
		selIf = ifis[n]
		selIfMatchBlocks = matchBlocks
	}

	return selIf, nil
}