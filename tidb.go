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
	"fmt"
	"os"

	//"github.com/pingcap/log"
	"github.com/pingcap/tidb/server"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/util/signal"
	_ "github.com/go-sql-driver/mysql"
)

// TiDBServer ...
type TiDBServer struct {
	cfg *Options

	svr      *server.Server
}

// NewTiDBServer returns a new TiDBServer
func NewTiDBServer(options *Options) *TiDBServer {
	registerStores()
	//configWarning := loadConfig()
	cfg = config.GetGlobalConfig()
	cfg.Store = "mocktikv"
	cfg.Path = options.DataDir
	cfg.Port = uint(options.Port)

	overrideConfig()
	if err := cfg.Valid(); err != nil {
		fmt.Fprintln(os.Stderr, "invalid config", err)
		os.Exit(1)
	}

	setGlobalVars()
	setupLog()
	
	printInfo()
	createStoreAndDomain()
	createServer()
	signal.SetupSignalHandler(serverShutdown)
	go runServer()
	
	return &TiDBServer {
		cfg: options,
	}
	//cleanup()
	//exit()
}

