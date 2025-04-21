package system

import "errors"

var (
	// ErrUnsupportedOS is returned when the operating system is not supported
	ErrUnsupportedOS = errors.New("unsupported operating system")

	// ErrNoPackageManager is returned when no package manager is found
	ErrNoPackageManager = errors.New("no supported package manager found")
) 