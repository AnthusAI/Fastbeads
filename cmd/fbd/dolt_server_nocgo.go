//go:build !dolt

package main

import (
	"context"
	"errors"
)

// DoltServerHandle is a stub for builds without Dolt support.
type DoltServerHandle struct{}

// DoltDefaultSQLPort is the default SQL port for dolt server
const DoltDefaultSQLPort = 3306

// DoltDefaultRemotesAPIPort is the default remotesapi port for dolt server
const DoltDefaultRemotesAPIPort = 50051

// ErrDoltRequiresCGO is returned when Dolt features are requested without Dolt support.
var ErrDoltRequiresCGO = errors.New("dolt backend requires the dolt build tag (and CGO); use pre-built binaries from GitHub releases or rebuild with -tags dolt")

// StartDoltServer returns an error when Dolt is not enabled.
func StartDoltServer(ctx context.Context, dataDir, logFile string, sqlPort, remotePort int) (*DoltServerHandle, error) {
	return nil, ErrDoltRequiresCGO
}

// Stop is a no-op stub
func (h *DoltServerHandle) Stop() error {
	return nil
}

// SQLPort returns 0 when Dolt is not enabled.
func (h *DoltServerHandle) SQLPort() int {
	return 0
}

// RemotesAPIPort returns 0 when Dolt is not enabled.
func (h *DoltServerHandle) RemotesAPIPort() int {
	return 0
}

// Host returns empty string when Dolt is not enabled.
func (h *DoltServerHandle) Host() string {
	return ""
}

// DoltServerAvailable returns false when Dolt is not enabled.
func DoltServerAvailable() bool {
	return false
}
