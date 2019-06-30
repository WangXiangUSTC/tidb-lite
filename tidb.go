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
	"sync/atomic"
	"os"
	"runtime"
	"context"
	"strconv"

	"github.com/pingcap/log"
	"github.com/pingcap/tidb/statistics"
	"github.com/pingcap/tidb/store/tikv"
	plannercore "github.com/pingcap/tidb/planner/core"
	"github.com/pingcap/tidb/util/memory"
	"github.com/pingcap/tidb/privilege/privileges"
	"github.com/pingcap/tidb/server"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/sessionctx/variable"
	"github.com/pingcap/tidb/session"
	"github.com/pingcap/tidb/store/mockstore"
	"github.com/pingcap/tidb/store/tikv/gcworker"
	"github.com/pingcap/tidb/ddl"
	"github.com/pingcap/tidb/statistics/handle"
	"github.com/pingcap/tidb/bindinfo"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/util/logutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/parser/terror"
	"github.com/pingcap/tidb/util/printer"
	kvstore "github.com/pingcap/tidb/store"
	"go.uber.org/zap"
)

// TiDBServer ...
type TiDBServer struct {
	//cfg *Options

	cfg      *config.Config
	svr      *server.Server
	storage  kv.Storage
	dom      *domain.Domain

	closeGracefully bool
}

// NewTiDBServer returns a new TiDBServer
func NewTiDBServer(options *Options) *TiDBServer {
	
	//configWarning := loadConfig()
	cfg := config.GetGlobalConfig()
	cfg.Store = "mocktikv"
	cfg.Path = options.DataDir
	cfg.Port = uint(options.Port)

	//overrideConfig()
	if err := cfg.Valid(); err != nil {
		fmt.Fprintln(os.Stderr, "invalid config", err)
		os.Exit(1)
	}
	tidbServer := &TiDBServer {
		cfg: cfg,
	}

	tidbServer.registerStores()
	tidbServer.setGlobalVars()
	tidbServer.setupLog()
	
	tidbServer.printInfo()
	tidbServer.createStoreAndDomain()
	tidbServer.createServer()
	//signal.SetupSignalHandler(serverShutdown)
	

	go func() {
		tidbServer.runServer()
		tidbServer.cleanup(tidbServer.closeGracefully)
	}()
	
	return tidbServer
}

func (t *TiDBServer) Close() {
	t.serverShutdown(false)
}

func (t *TiDBServer) CloseGracefully() {
	t.serverShutdown(true)
}

func (t *TiDBServer) printInfo() {
	// Make sure the TiDB info is always printed.
	level := log.GetLevel()
	log.SetLevel(zap.InfoLevel)
	printer.PrintTiDBInfo()
	log.SetLevel(level)
}

func (t *TiDBServer) registerStores() {
	err := kvstore.Register("tikv", tikv.Driver{})
	terror.MustNil(err)
	tikv.NewGCHandlerFunc = gcworker.NewGCWorker
	err = kvstore.Register("mocktikv", mockstore.MockDriver{})
	terror.MustNil(err)
}

func (t *TiDBServer) createServer() {
	driver := server.NewTiDBDriver(t.storage)
	var err error
	t.svr, err = server.NewServer(t.cfg, driver)
	// Both domain and storage have started, so we have to clean them before exiting.
	terror.MustNil(err, t.closeDomainAndStorage)
	go t.dom.ExpensiveQueryHandle().SetSessionManager(t.svr).Run()
}

func (t *TiDBServer) runServer() {
	err := t.svr.Run()
	terror.MustNil(err)
}

func (t *TiDBServer) createStoreAndDomain() {
	fullPath := fmt.Sprintf("%s://%s", t.cfg.Store, t.cfg.Path)
	var err error
	t.storage, err = kvstore.New(fullPath)
	terror.MustNil(err)
	// Bootstrap a session to load information schema.
	t.dom, err = session.BootstrapSession(t.storage)
	terror.MustNil(err)
}

func (t *TiDBServer) setGlobalVars() {
	ddlLeaseDuration := parseDuration(t.cfg.Lease)
	session.SetSchemaLease(ddlLeaseDuration)
	runtime.GOMAXPROCS(int(t.cfg.Performance.MaxProcs))
	statsLeaseDuration := parseDuration(t.cfg.Performance.StatsLease)
	session.SetStatsLease(statsLeaseDuration)
	bindinfo.Lease = parseDuration(t.cfg.Performance.BindInfoLease)
	domain.RunAutoAnalyze = t.cfg.Performance.RunAutoAnalyze
	statistics.FeedbackProbability.Store(t.cfg.Performance.FeedbackProbability)
	handle.MaxQueryFeedbackCount.Store(int64(t.cfg.Performance.QueryFeedbackLimit))
	statistics.RatioOfPseudoEstimate.Store(t.cfg.Performance.PseudoEstimateRatio)
	ddl.RunWorker = t.cfg.RunDDL
	if t.cfg.SplitTable {
		atomic.StoreUint32(&ddl.EnableSplitTableRegion, 1)
	}
	plannercore.AllowCartesianProduct.Store(t.cfg.Performance.CrossJoin)
	privileges.SkipWithGrant = t.cfg.Security.SkipGrantTable

	priority := mysql.Str2Priority(t.cfg.Performance.ForcePriority)
	variable.ForcePriority = int32(priority)
	variable.SysVars[variable.TiDBForcePriority].Value = mysql.Priority2Str[priority]

	variable.SysVars[variable.TIDBMemQuotaQuery].Value = strconv.FormatInt(t.cfg.MemQuotaQuery, 10)
	variable.SysVars["lower_case_table_names"].Value = strconv.Itoa(t.cfg.LowerCaseTableNames)
	variable.SysVars[variable.LogBin].Value = variable.BoolToIntStr(config.GetGlobalConfig().Binlog.Enable)

	variable.SysVars[variable.Port].Value = fmt.Sprintf("%d", t.cfg.Port)
	variable.SysVars[variable.Socket].Value = t.cfg.Socket
	variable.SysVars[variable.DataDir].Value = t.cfg.Path
	variable.SysVars[variable.TiDBSlowQueryFile].Value = t.cfg.Log.SlowQueryFile

	// For CI environment we default enable prepare-plan-cache.
	plannercore.SetPreparedPlanCache(config.CheckTableBeforeDrop || t.cfg.PreparedPlanCache.Enabled)
	if plannercore.PreparedPlanCacheEnabled() {
		plannercore.PreparedPlanCacheCapacity = t.cfg.PreparedPlanCache.Capacity
		plannercore.PreparedPlanCacheMemoryGuardRatio = t.cfg.PreparedPlanCache.MemoryGuardRatio
		if plannercore.PreparedPlanCacheMemoryGuardRatio < 0.0 || plannercore.PreparedPlanCacheMemoryGuardRatio > 1.0 {
			plannercore.PreparedPlanCacheMemoryGuardRatio = 0.1
		}
		plannercore.PreparedPlanCacheMaxMemory.Store(t.cfg.Performance.MaxMemory)
		total, err := memory.MemTotal()
		terror.MustNil(err)
		if plannercore.PreparedPlanCacheMaxMemory.Load() > total || plannercore.PreparedPlanCacheMaxMemory.Load() <= 0 {
			plannercore.PreparedPlanCacheMaxMemory.Store(total)
		}
	}

	tikv.CommitMaxBackoff = int(parseDuration(t.cfg.TiKVClient.CommitTimeout).Seconds() * 1000)
	tikv.PessimisticLockTTL = uint64(parseDuration(t.cfg.PessimisticTxn.TTL).Seconds() * 1000)
}

func (t *TiDBServer) serverShutdown(isgraceful bool) {
	t.closeGracefully = isgraceful
	t.svr.Close()
}

func (t *TiDBServer) closeDomainAndStorage() {
	atomic.StoreUint32(&tikv.ShuttingDown, 1)
	t.dom.Close()
	err := t.storage.Close()
	terror.Log(errors.Trace(err))
}

func (t *TiDBServer) cleanup(graceful bool) {
	if t.closeGracefully {
		t.svr.GracefulDown(context.Background(), nil)
	} else {
		t.svr.TryGracefulDown()
	}

	t.closeDomainAndStorage()
}

func (t *TiDBServer) setupLog() {
	err := logutil.InitZapLogger(t.cfg.Log.ToLogConfig())
	terror.MustNil(err)

	err = logutil.InitLogger(t.cfg.Log.ToLogConfig())
	terror.MustNil(err)
}