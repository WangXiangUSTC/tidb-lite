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

package tidblite

import (
	"fmt"
)

const (
	// DefaultSocket is the default socket to use
	DefaultSocket   = "/tmp/tidb-lite-socket"
	DefaultConnOpts = "charset=utf8mb4"
)

type Options struct {
	// DataDir is the directory to save db data
	DataDir string

	Socket string

	Port int

	//Connection options pass to the database, separated with & like charset=utf8mb4&parseTime=True
	ConnOpts string
}

// NewOptions creates a new Options with default port
func NewOptions(dataDir string) *Options {
	return &Options{
		DataDir:  dataDir,
		Socket:   fmt.Sprintf("%s-%s.sock", DefaultSocket, randomStringWithCharset(5)),
		ConnOpts: DefaultConnOpts,
	}
}

// WithPort set the port
func (o *Options) WithPort(port int) *Options {
	o.Port = port
	return o
}

// WithSocket set the socket
func (o *Options) WithSocket(socket string) *Options {
	o.Socket = socket
	return o
}

// WithConnOpts overrides default connection options
func (o *Options) WithConnOpts(connOpts string) *Options {
	o.ConnOpts = connOpts
	return o
}
