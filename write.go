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

// OpenWrite opens the specified file for reading.
// The associated file descriptor has mode O_WRONLY|O_TRUNC.
//
// If the specified doesn't exist, OpenWrite will check for a default.
// If a default is found, it will call that and then try to read.
// Otherwise, an error will be returned.
func (c *Config) OpenWrite(name string) (io.WriteCloser, error) {
	if err := c.verifyExists(name); err != nil {
		return nil, err
	}

	return c.OpenRaw(name, os.O_WRONLY|os.O_TRUNC, 0)
}

// OpenWriteAuto is the same as OpenWrite, but the file is closed automatically with [Close].
// This is useful when a file needs to be open for the entirety of an application's runtime, such as when using [log.New].
func (c *Config) OpenWriteAuto(name string) (io.Writer, error) {
	w, err := c.OpenWrite(name)
	if err != nil {
		return nil, err
	}

	c.addCloser(w)

	return w, err
}

// Write writes the entirety of data into the specified file, truncating it.
func (c *Config) Write(name string, data []byte) (int, error) {
	w, err := c.OpenWrite(name)
	if err != nil {
		return 0, err
	}

	defer w.Close()

	return w.Write(data)
}

// WriteString writes the entirety of str into the specified file, truncating it.
func (c *Config) WriteString(name, str string) (int, error) {
	return c.Write(name, []byte(str))
}

// WriteJson serializes v as JSON and writes it into the specified file.
func (c *Config) WriteJson(name string, v any) error {
	w, err := c.OpenWrite(name)
	if err != nil {
		return err
	}

	defer w.Close()

	return json.NewEncoder(w).Encode(v)
}

// WriteJson serializes v as TOML and writes it into the specified file.
func (c *Config) WriteToml(name string, v any) error {
	w, err := c.OpenWrite(name)
	if err != nil {
		return err
	}

	defer w.Close()

	return toml.NewEncoder(w).Encode(v)
}

// WriteJson serializes records as a CSV and writes it into the specified file.
func (c *Config) WriteCsv(name string, records [][]string) error {
	w, err := c.OpenWrite(name)
	if err != nil {
		return err
	}

	defer w.Close()

	return csv.NewWriter(w).WriteAll(records)
}

// WriteJson serializes v as XML and writes it into the specified file.
func (c *Config) WriteXml(name string, v any) error {
	w, err := c.OpenWrite(name)
	if err != nil {
		return err
	}

	defer w.Close()

	return xml.NewEncoder(w).Encode(v)
}
