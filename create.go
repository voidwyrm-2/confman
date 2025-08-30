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

// OpenCreate creates and opens the specified file for reading or writing.
// The associated file descriptor has mode O_WRONLY|O_CREATE|O_EXCL.
//
// An error will be returned if the specified file already exists.
func (c *Config) OpenCreate(name string, perm os.FileMode) (io.ReadWriteCloser, error) {
	return c.OpenRaw(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, perm)
}

// OpenCreateAuto is the same as OpenCreate, but the file is closed automatically with [Close].
// This is useful when a file needs to be open for the entirety of an application's runtime, such as when using [log.NewLogger].
func (c *Config) OpenCreateAuto(name string, perm os.FileMode) (io.ReadWriter, error) {
	w, err := c.OpenCreate(name, perm)
	if err != nil {
		return nil, err
	}

	c.addCloser(w)

	return w, err
}

// Create creates the specified file, then writes the entirety of data into it.
func (c *Config) Create(name string, data []byte, perm os.FileMode) (int, error) {
	w, err := c.OpenCreate(name, perm)
	if err != nil {
		return 0, err
	}

	defer w.Close()

	return w.Write(data)
}

// CreateString creates the specified file, then writes the entirety of str into the file.
func (c *Config) CreateString(name, str string, perm os.FileMode) (int, error) {
	return c.Create(name, []byte(str), perm)
}

// CreateJson creates the specified file, serializes v as JSON, then writes it into the file.
func (c *Config) CreateJson(name string, v any, perm os.FileMode) error {
	w, err := c.OpenCreate(name, perm)
	if err != nil {
		return err
	}

	defer w.Close()

	return json.NewEncoder(w).Encode(v)
}

// CreateToml creates the specified file, serializes v as TOML, then writes it into the file.
func (c *Config) CreateToml(name string, v any, perm os.FileMode) error {
	w, err := c.OpenCreate(name, perm)
	if err != nil {
		return err
	}

	defer w.Close()

	return toml.NewEncoder(w).Encode(v)
}

// CreateCSV creates the specified file, serializes records as a CSV, then writes it into the file.
func (c *Config) CreateCsv(name string, records [][]string, perm os.FileMode) error {
	w, err := c.OpenCreate(name, perm)
	if err != nil {
		return err
	}

	defer w.Close()

	return csv.NewWriter(w).WriteAll(records)
}

// CreateXml creates the specified file, serializes v as XML, then writes it into it.
func (c *Config) CreateXml(name string, v any, perm os.FileMode) error {
	w, err := c.OpenCreate(name, perm)
	if err != nil {
		return err
	}

	defer w.Close()

	return xml.NewEncoder(w).Encode(v)
}
