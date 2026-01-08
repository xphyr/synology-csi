// Copyright 2021 Synology Inc.

package utils

import (
	"fmt"
	"net"
	"strings"
)

type AuthType string

const (
	UnitGB = 1024 * 1024 * 1024
	UnitMB = 1024 * 1024

	ProtocolSmb     = "smb"
	ProtocolIscsi   = "iscsi"
	ProtocolNfs     = "nfs"
	ProtocolDefault = ProtocolIscsi

	AuthTypeReadWrite AuthType = "rw"
	AuthTypeReadOnly  AuthType = "ro"
	AuthTypeNoAccess  AuthType = "no"
)

func SliceContains(items []string, s string) bool {
	for _, item := range items {
		if s == item {
			return true
		}
	}
	return false
}

func MBToBytes(size int64) int64 {
	return size * UnitMB
}

func BytesToMB(size int64) int64 {
	return size / UnitMB
}

// Ceiling
func BytesToMBCeil(size int64) int64 {
	return (size + UnitMB - 1) / UnitMB
}

func StringToBoolean(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "yes" || value == "true" || value == "1"
}
func StringToSlice(value string) []string {
	return strings.Fields(value)
}

func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// Haven't supported IPv6 yet.
func LookupIPv4(name string) ([]string, error) {
	ips, _ := net.LookupIP(name)

	retIps := []string{}
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			retIps = append(retIps, ipv4.String())
		}
	}
	if len(retIps) > 0 {
		return retIps, nil
	}

	return nil, fmt.Errorf("failed to LookupIPv4 by local resolver for: %s", name)
}
