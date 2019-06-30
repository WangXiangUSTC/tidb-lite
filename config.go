// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (

)

const (
	// DefaultPort is the default port to use
	DefaultPort = 4040
)

type Options struct {
	// DataDir is the directory to save db data
	DataDir string

	Port int
}

// NewOptions creates a new Options with default port
func NewOptions(dataDir string) *Options {
	return &Options {
		DataDir: dataDir,
		Port:    DefaultPort,
	}
}

// WithPort set the port
func (o *Options) WithPort(port int) {
	o.Port = port
}
