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
	"context"
	"time"
	"database/sql"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	_ "github.com/go-sql-driver/mysql"
)

// Flag Names
const (
	/*
	nmVersion          = "V"
	nmConfig           = "config"
	nmConfigCheck      = "config-check"
	nmConfigStrict     = "config-strict"
	nmStore            = "store"
	nmStorePath        = "path"
	nmHost             = "host"
	nmAdvertiseAddress = "advertise-address"
	nmPort             = "P"
	nmCors             = "cors"
	nmSocket           = "socket"
	nmEnableBinlog     = "enable-binlog"
	nmRunDDL           = "run-ddl"
	nmLogLevel         = "L"
	nmLogFile          = "log-file"
	nmLogSlowQuery     = "log-slow-query"
	nmReportStatus     = "report-status"
	nmStatusHost       = "status-host"
	nmStatusPort       = "status"
	nmMetricsAddr      = "metrics-addr"
	nmMetricsInterval  = "metrics-interval"
	nmDdlLease         = "lease"
	nmTokenLimit       = "token-limit"
	nmPluginDir        = "plugin-dir"
	nmPluginLoad       = "plugin-load"

	nmProxyProtocolNetworks      = "proxy-protocol-networks"
	nmProxyProtocolHeaderTimeout = "proxy-protocol-header-timeout"
	*/
)

var (
	/*
	version = false
	store = "mocktikv"
	storePath = "/tmp/tidb"
	host = "0.0.0.0"
	port = DefaultPort
	runDDL = true
	ddlLease = "45s"
	tokenLimit = 1000
	
	reportStatus = false
	*/
	//statusHost = "0.0.0.0"
	//statusPort = 10080

	//version      = flagBoolean(nmVersion, false, "print version information and exit")
	//configPath   = flag.String(nmConfig, "", "config file path")
	//configCheck  = flagBoolean(nmConfigCheck, false, "check config file validity and exit")
	//configStrict = flagBoolean(nmConfigStrict, false, "enforce config file validity")

	// Base
	//store            = flag.String(nmStore, "mocktikv", "registered store name, [tikv, mocktikv]")
	//storePath        = flag.String(nmStorePath, "/tmp/tidb", "tidb storage path")
	//host             = flag.String(nmHost, "0.0.0.0", "tidb server host")
	//advertiseAddress = flag.String(nmAdvertiseAddress, "", "tidb server advertise IP")
	//port             = flag.String(nmPort, "4000", "tidb server port")
	//cors             = flag.String(nmCors, "", "tidb server allow cors origin")
	//socket           = flag.String(nmSocket, "", "The socket file to use for connection.")
	//enableBinlog     = flagBoolean(nmEnableBinlog, false, "enable generate binlog")
	//runDDL           = flagBoolean(nmRunDDL, true, "run ddl worker on this tidb-server")
	//ddlLease         = flag.String(nmDdlLease, "45s", "schema lease duration, very dangerous to change only if you know what you do")
	//tokenLimit       = flag.Int(nmTokenLimit, 1000, "the limit of concurrent executed sessions")
	//pluginDir        = flag.String(nmPluginDir, "/data/deploy/plugin", "the folder that hold plugin")
	//pluginLoad       = flag.String(nmPluginLoad, "", "wait load plugin name(separated by comma)")

	// Log
	//logLevel     = flag.String(nmLogLevel, "info", "log level: info, debug, warn, error, fatal")
	//logFile      = flag.String(nmLogFile, "", "log file path")
	//logSlowQuery = flag.String(nmLogSlowQuery, "", "slow query file path")

	// Status
	//reportStatus    = flagBoolean(nmReportStatus, true, "If enable status report HTTP service.")
	//statusHost      = flag.String(nmStatusHost, "0.0.0.0", "tidb server status host")
	//statusPort      = flag.String(nmStatusPort, "10080", "tidb server status port")
	//metricsAddr     = flag.String(nmMetricsAddr, "", "prometheus pushgateway address, leaves it empty will disable prometheus push.")
	//metricsInterval = flag.Uint(nmMetricsInterval, 15, "prometheus client push interval in second, set \"0\" to disable prometheus push.")

	// PROXY Protocol
	//proxyProtocolNetworks      = flag.String(nmProxyProtocolNetworks, "", "proxy protocol networks allowed IP or *, empty mean disable proxy protocol support")
	//proxyProtocolHeaderTimeout = flag.Uint(nmProxyProtocolHeaderTimeout, 5, "proxy protocol header read timeout, unit is second.")
)

var (
	//cfg      *config.Config
	//storage  kv.Storage
	//dom      *domain.Domain
	//svr      *server.Server
	//graceful bool
)

/*
func StartTiDB() {
	flag.Parse()
	if *version {
		fmt.Println(printer.GetTiDBInfo())
		os.Exit(0)
	}

	registerStores()
	loadConfig()
	//overrideConfig()
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
	runServer()
	cleanup()
	exit()
}
*/

func main() {
	/*
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Wait()
		StartTiDB()
	}()
	time.Sleep(5*time.Second)

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4", "root", "", "127.0.0.1", 4000)
	dbConn, err := sql.Open("mysql", dbDSN)
	if err != nil {
		panic(err.Error())
	}

	version, err := GetDBVersion(context.Background(), dbConn)
	if err != nil {
		panic(err.Error())
	}
	log.Info("", zap.String("tidb version", version))

	svr.Close()

	wg.Wait()
	*/
}

// GetDBVersion returns the database's version
func GetDBVersion(ctx context.Context, db *sql.DB) (string, error) {
	/*
		example in TiDB:
		mysql> select version();
		+--------------------------------------+
		| version()                            |
		+--------------------------------------+
		| 5.7.10-TiDB-v2.1.0-beta-173-g7e48ab1 |
		+--------------------------------------+

		example in MySQL:
		mysql> select version();
		+-----------+
		| version() |
		+-----------+
		| 5.7.21    |
		+-----------+
	*/
	query := "SELECT version()"
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer result.Close()

	var version sql.NullString
	for result.Next() {
		err := result.Scan(&version)
		if err != nil {
			return "", errors.Trace(err)
		}
		break
	}

	if version.Valid {
		return version.String, nil
	}

	return "", nil
}


// parseDuration parses lease argument string.
func parseDuration(lease string) time.Duration {
	dur, err := time.ParseDuration(lease)
	if err != nil {
		dur, err = time.ParseDuration(lease + "s")
	}
	if err != nil || dur < 0 {
		log.Fatal("invalid lease duration", zap.String("lease", lease))
	}
	return dur
}

/*
func loadConfig() {
	cfg = config.GetGlobalConfig()
	cfg.Store = store
	cfg.Host = host
	//cfg.AdvertiseAddress = advertiseAddress
	cfg.Port = uint(port)
	cfg.Store = store
	cfg.Path = storePath
	cfg.RunDDL = runDDL
	cfg.Lease = ddlLease
	cfg.TokenLimit = uint(tokenLimit)

	// Status
	//cfg.Status.ReportStatus = reportStatus
	//cfg.Status.StatusHost = statusHost
	//cfg.Status.StatusPort = uint(statusPort)
	
	if *configPath != "" {
		// Not all config items are supported now.
		config.SetConfReloader(*configPath, reloadConfig, hotReloadConfigItems...)

		err := cfg.Load(*configPath)
		// This block is to accommodate an interim situation where strict config checking
		// is not the default behavior of TiDB. The warning message must be deferred until
		// logging has been set up. After strict config checking is the default behavior,
		// This should all be removed.
		if _, ok := err.(*config.ErrConfigValidationFailed); ok && !*configStrict {
			return err.Error()
		}
		terror.MustNil(err)
	}
}
*/

/*
func reloadConfig(nc, c *config.Config) {
	// Just a part of config items need to be reload explicitly.
	// Some of them like OOMAction are always used by getting from global config directly
	// like config.GetGlobalConfig().OOMAction.
	// These config items will become available naturally after the global config pointer
	// is updated in function ReloadGlobalConfig.
	if nc.Performance.MaxProcs != c.Performance.MaxProcs {
		runtime.GOMAXPROCS(int(nc.Performance.MaxProcs))
	}
	if nc.Performance.MaxMemory != c.Performance.MaxMemory {
		plannercore.PreparedPlanCacheMaxMemory.Store(nc.Performance.MaxMemory)
	}
	if nc.Performance.CrossJoin != c.Performance.CrossJoin {
		plannercore.AllowCartesianProduct.Store(nc.Performance.CrossJoin)
	}
	if nc.Performance.FeedbackProbability != c.Performance.FeedbackProbability {
		statistics.FeedbackProbability.Store(nc.Performance.FeedbackProbability)
	}
	if nc.Performance.QueryFeedbackLimit != c.Performance.QueryFeedbackLimit {
		handle.MaxQueryFeedbackCount.Store(int64(nc.Performance.QueryFeedbackLimit))
	}
	if nc.Performance.PseudoEstimateRatio != c.Performance.PseudoEstimateRatio {
		statistics.RatioOfPseudoEstimate.Store(nc.Performance.PseudoEstimateRatio)
	}
}
*/

/*
func overrideConfig() {
	//actualFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		actualFlags[f.Name] = true
	})

	// Base
	if actualFlags[nmHost] {
		cfg.Host = *host
	}
	if actualFlags[nmAdvertiseAddress] {
		cfg.AdvertiseAddress = *advertiseAddress
	}
	var err error
	if actualFlags[nmPort] {
		var p int
		p, err = strconv.Atoi(*port)
		terror.MustNil(err)
		cfg.Port = uint(p)
	}
	if actualFlags[nmCors] {
		fmt.Println(cors)
		cfg.Cors = *cors
	}
	if actualFlags[nmStore] {
		cfg.Store = *store
	}
	if actualFlags[nmStorePath] {
		cfg.Path = *storePath
	}
	if actualFlags[nmSocket] {
		cfg.Socket = *socket
	}
	if actualFlags[nmEnableBinlog] {
		cfg.Binlog.Enable = *enableBinlog
	}
	if actualFlags[nmRunDDL] {
		cfg.RunDDL = *runDDL
	}
	if actualFlags[nmDdlLease] {
		cfg.Lease = *ddlLease
	}
	if actualFlags[nmTokenLimit] {
		cfg.TokenLimit = uint(*tokenLimit)
	}
	if actualFlags[nmPluginLoad] {
		cfg.Plugin.Load = *pluginLoad
	}
	if actualFlags[nmPluginDir] {
		cfg.Plugin.Dir = *pluginDir
	}

	// Log
	if actualFlags[nmLogLevel] {
		cfg.Log.Level = *logLevel
	}
	if actualFlags[nmLogFile] {
		cfg.Log.File.Filename = *logFile
	}
	if actualFlags[nmLogSlowQuery] {
		cfg.Log.SlowQueryFile = *logSlowQuery
	}

	// Status
	if actualFlags[nmReportStatus] {
		cfg.Status.ReportStatus = *reportStatus
	}
	if actualFlags[nmStatusHost] {
		cfg.Status.StatusHost = *statusHost
	}
	if actualFlags[nmStatusPort] {
		var p int
		p, err = strconv.Atoi(*statusPort)
		terror.MustNil(err)
		cfg.Status.StatusPort = uint(p)
	}
	if actualFlags[nmMetricsAddr] {
		cfg.Status.MetricsAddr = *metricsAddr
	}
	if actualFlags[nmMetricsInterval] {
		cfg.Status.MetricsInterval = *metricsInterval
	}

	// PROXY Protocol
	if actualFlags[nmProxyProtocolNetworks] {
		cfg.ProxyProtocol.Networks = *proxyProtocolNetworks
	}
	if actualFlags[nmProxyProtocolHeaderTimeout] {
		cfg.ProxyProtocol.HeaderTimeout = *proxyProtocolHeaderTimeout
	}
}
*/












