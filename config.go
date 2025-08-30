// Copyright 2025 Nuclear Pasta. All rights reserved.
// Use of this source code is governed by the MIT license.
// license that can be found in the LICENSE file.

package confman

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type fileDefault struct {
	f    func(*Config, io.Writer) error
	perm os.FileMode
}

// A Config represents the config of an application, with functions for convienently reading and writing various types of data to files in the configuration directory.
type Config struct {
	defaults map[string]fileDefault
	closers  []io.Closer
	path     string
}

// OpenSpecific creates a Config pointing to the specified path.
func OpenSpecific(path string) (*Config, error) {
	path = filepath.Clean(path)

	if path == "" {
		panic("path cannot be empty")
	}

	conf := &Config{
		defaults: map[string]fileDefault{},
		path:     path,
	}

	err := conf.create()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// OpenHome creates a Config pointing to name beginning with a period inside of the system's home directory.
func OpenHome(name string) (*Config, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", err
	}

	path := filepath.Join(home, name)

	conf, err := OpenSpecific(path)
	if err != nil {
		return nil, "", err
	}

	return conf, path, nil
}

// Open creates a Config pointing to name inside of the system's default config directory.
func Open(name string) (*Config, string, error) {
	configPath, dot, err := GetConfigPathForSystem()
	if err != nil {
		return nil, "", err
	}

	if _, err := os.Stat(configPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
		}

		return nil, "", err
	}

	name = strings.TrimLeft(name, ".")

	if dot {
		name = "." + name
	}

	path := filepath.Join(configPath, name)

	conf, err := OpenSpecific(path)
	if err != nil {
		return nil, "", err
	}

	return conf, path, nil
}

// Closes all files created with [OpenReadAuto], [OpenWriteAuto], or [OpenCreateAuto].
//
// A list of errors from each close call is returned.
func (c *Config) Close() []error {
	errors := []error{}

	for _, closer := range c.closers {
		if err := closer.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// Delete removes the directory pointed to by this Config along with all the data inside it.
//
// This is a very dangerous function, it can lead to unrecoverable data loss.
func (c *Config) Delete() error {
	return os.RemoveAll(c.path)
}

func (c *Config) addCloser(closer io.Closer) {
	c.closers = append(c.closers, closer)
}

func (c *Config) child(name string) string {
	return filepath.Join(c.path, name)
}

// Stat returns a [FileInfo] describing the specified file.
// If there is an error, it will be of type [*PathError].
func (c *Config) Stat(name string) (os.FileInfo, error) {
	return os.Stat(c.child(name))
}

// Exists checks if specified exists or not.
// Under normal circumstances, it should not return an error.
// If there is an error, it will be of type [*PathError].
func (c *Config) Exists(name string) (bool, error) {
	if _, err := c.Stat(name); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// DeleteFile removes the specified file from the configuration directory.
//
// This is a very dangerous function, it can lead to unrecoverable data loss.
func (c *Config) DeleteFile(name string) error {
	return os.RemoveAll(c.path)
}

func (c *Config) create() error {
	if exists, err := c.Exists(""); err != nil {
		return err
	} else if exists {
		return nil
	}

	return os.Mkdir(c.path, os.ModeDir|0o777)
}

func (c *Config) verifyExists(name string) error {
	exists, err := c.Exists(name)
	if err != nil {
		return err
	} else if !exists {
		if def, ok := c.defaults[name]; ok {
			w, err := c.OpenCreate(name, def.perm)
			if err != nil {
				fmt.Println("err:", err.Error())
				return err
			}

			return def.f(c, w)
		}

		return fmt.Errorf("'%s' does not exist", name)
	}

	return nil
}

func (c *Config) verifyNotExists(name string) error {
	exists, err := c.Exists(name)
	if err != nil {
		return err
	} else if exists {
		return fmt.Errorf("'%s' already exists", name)
	}

	return nil
}

// DefaultFunc sets the callback for if the specified file doesn't exist when call any Config.Read, Config.Write function, or any of their variants.
func (c *Config) DefaultFunc(name string, perm os.FileMode, f func(*Config, io.Writer) error) {
	if _, ok := c.defaults[name]; ok {
		panic(fmt.Sprintf("'%s' already has a default", name))
	}

	c.defaults[name] = fileDefault{perm: perm, f: f}
}

// Default sets the bytes written to the specified file if it doesn't exist when call any Config.Read, Config.Write, or any of their variants.
//
// This is a convienence function over top of [DefaultFunc]
func (c *Config) Default(name string, perm os.FileMode, data []byte) {
	c.DefaultFunc(name, perm, func(_ *Config, w io.Writer) error {
		_, err := w.Write(data)
		return err
	})
}

// DefaultString sets the string written to the specified file if it doesn't exist when call any Config.Read, Config.Write, or any of their variants.
//
// This is a convienence function over top of [Default]
func (c *Config) DefaultString(name string, perm os.FileMode, str string) {
	c.Default(name, perm, []byte(str))
}

// OpenRaw is the generalized open call; most users will use [OpenRead],
// [OpenWrite], or [OpenCreate] instead. It opens the named file with specified flag and file permissions.
func (c *Config) OpenRaw(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	return os.OpenFile(c.child(name), flag, perm)
}
