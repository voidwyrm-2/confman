// Copyright 2025 Nuclear Pasta. All rights reserved.
// Use of this source code is governed by the MIT license.
// license that can be found in the LICENSE file.

package confman

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// OpenRead opens the specified file for reading.
// The associated file descriptor has mode O_RDONLY.
//
// If the specified doesn't exist, OpenRead will check for a default.
// If a default is found, it will call that and then try to read.
// Otherwise, an error will be returned.
func (c *Config) OpenRead(name string) (io.ReadCloser, error) {
	if err := c.verifyExists(name); err != nil {
		return nil, err
	}

	return c.OpenRaw(name, os.O_RDONLY, 0)
}

// OpenReadAuto is the same as OpenRead, but the file is closed automatically with [Close].
// This is useful when a file needs to be open for the entirety of an application's runtime, such as when using [log.NewLogger].
func (c *Config) OpenReadAuto(name string) (io.Reader, error) {
	r, err := c.OpenRead(name)
	if err != nil {
		return nil, err
	}

	c.addCloser(r)

	return r, err
}

// Read reads the entirety of the specified file, then returns it as a slice of bytes.
func (c *Config) Read(name string) ([]byte, error) {
	fr, err := c.OpenRead(name)
	if err != nil {
		return nil, err
	}

	defer fr.Close()

	return io.ReadAll(fr)
}

// Read reads the entirety of the specified file, then returns it as a string.
func (c *Config) ReadString(name string) (string, error) {
	b, err := c.Read(name)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// ReadJson parses the specified file as JSON, then stores the result into the value pointed to by v.
func (c *Config) ReadJson(name string, v any) error {
	r, err := c.OpenRead(name)
	if err != nil {
		return err
	}

	defer r.Close()

	return json.NewDecoder(r).Decode(v)
}

// ReadJson parses the specified file as TOMl, then stores the result into the value pointed to by v.
func (c *Config) ReadToml(name string, v any) error {
	r, err := c.OpenRead(name)
	if err != nil {
		return err
	}

	defer r.Close()

	_, err = toml.NewDecoder(r).Decode(v)
	return err
}

// ReadJson parses the specified file as a CSV, then returns the records as a slice of slices of strings
func (c *Config) ReadCsv(name string) ([][]string, error) {
	r, err := c.OpenRead(name)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	return csv.NewReader(r).ReadAll()
}

// ReadJson parses the specified file as XML, then stores the result into the value pointed to by v.
func (c *Config) ReadXml(name string, v any) error {
	r, err := c.OpenRead(name)
	if err != nil {
		return err
	}

	defer r.Close()

	return xml.NewDecoder(r).Decode(v)
}
