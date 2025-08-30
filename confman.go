// Copyright 2025 Nuclear Pasta. All rights reserved.
// Use of this source code is governed by the MIT license.
// license that can be found in the LICENSE file.

// Package confman implements a system for managing the configuration files of an application
// in a way compatible with most conventional systems.
package confman

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigPathForSystem attempts to find the configuration path for the current system.
//
// On darwin and a-Shell, this is $HOME/.config
// For other platforms, refer to [os.UserConfigDir]'s documentation.
func GetConfigPathForSystem() (path string, dot bool, err error) {
	dot = false

	switch runtime.GOOS {
	case "darwin", "wasip1":
		path, err = os.UserHomeDir()
		if err != nil {
			break
		}

		path = filepath.Join(path, ".config")

		if os.Getenv("APPNAME") == "a-Shell" {
			err = os.MkdirAll(path, os.ModeDir|0o777)
		}
	default:
		path, err = os.UserConfigDir()
		// err = fmt.Errorf("Cannot find config directory for '%s'", runtime.GOOS)
	}

	return
}
