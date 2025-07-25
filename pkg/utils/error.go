// Copyright 2021 Synology Inc.

package utils

import (
	"fmt"
)

type OutOfFreeSpaceError string
type AlreadyExistError string
type BadParametersError string
type NoSuchLunError string
type LunReachMaxCountError string
type TargetReachMaxCountError string
type NoSuchSnapshotError string
type BadLunTypeError string
type SnapshotReachMaxCountError string
type IscsiDefaultError struct {
	ErrCode int
}
type NoSuchShareError string
type ShareReachMaxCountError string
type ShareSystemBusyError string
type ShareDefaultError struct {
	ErrCode int
}

func (OutOfFreeSpaceError) Error() string {
	return "Out of free space"
}
func (AlreadyExistError) Error() string {
	return "Already Existed"
}
func (BadParametersError) Error() string {
	return "Invalid input value"
}

// ISCSI errors
func (NoSuchLunError) Error() string {
	return "No such LUN"
}

func (LunReachMaxCountError) Error() string {
	return "Number of LUN reach limit"
}

func (TargetReachMaxCountError) Error() string {
	return "Number of target reach limit"
}

func (NoSuchSnapshotError) Error() string {
	return "No such snapshot uuid"
}

func (BadLunTypeError) Error() string {
	return "Bad LUN type"
}

func (SnapshotReachMaxCountError) Error() string {
	return "Number of snapshot reach limit"
}

func (e IscsiDefaultError) Error() string {
	return fmt.Sprintf("ISCSI API error. Error code: %d", e.ErrCode)
}

// Share errors
func (NoSuchShareError) Error() string {
	return "No such share"
}

func (ShareReachMaxCountError) Error() string {
	return "Number of share reach limit"
}

func (ShareSystemBusyError) Error() string {
	return "Share system is temporary busy"
}

func (e ShareDefaultError) Error() string {
	return fmt.Sprintf("Share API error. Error code: %d", e.ErrCode)
}
