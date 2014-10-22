//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"flag"
	realhttp "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/couchbase/indexing/secondary/common"
	"github.com/couchbaselabs/clog"
	"github.com/couchbaselabs/tuqtng/ast"
	"github.com/couchbaselabs/tuqtng/network"
	"github.com/couchbaselabs/tuqtng/network/http"
	"github.com/couchbaselabs/tuqtng/server"
)

var VERSION = "0.0.0" // Build-time overriddable.

var addr = flag.String("addr", ":8093", "HTTP listen address")
var couchbaseSite = flag.String("couchbase", "", "Couchbase Cluster Address (http://...) or dir:PATH")
var defaultPoolName = flag.String("pool", "default", "Default Pool")
var logKeys = flag.String("log", "", "Log keywords, comma separated")
var devMode = flag.Bool("dev", false, "Developer Mode")
var profileMode = flag.Bool("profile", false, "Profile Mode")
var staticPath = flag.String("staticPath", "static", "Path to static content")
var queryTimeout = flag.Duration("queryTimeout", -1*time.Second, "Query execution timeout, negative values disable timeout")

var devModeDefaultLogKeys = []string{"HTTP", "NETWORK", "PIPELINE", "CATALOG", "PLANNER", "SCAN", "OPTIMIZER", "PARSER"}
var disableInfo = flag.Bool("disableInfo", false, "Enable query info line")

func main() {
	flag.Parse()

	if *devMode {
		ast.EnableDeveloperFunctions()
		if *logKeys == "" {
			clog.ParseLogFlags(devModeDefaultLogKeys)
			clog.Log("Developer Mode Enabled (default developer logging)")
		} else {
			clog.ParseLogFlag(*logKeys)
			clog.Log("Developer Mode Enabled (custom command-line logging)")
		}
	}

	clog.Log("Info line disabled %v", *disableInfo)

	if *profileMode {
		clog.Log("Enabling HTTP Profiling on :6060")
		go func() {
			realhttp.ListenAndServe(":6060", nil)
		}()
	}

	go dumpOnSignalForPlatform()
	go common.ExitOnStdinClose()

	// create a QueryChannel
	queryChannel := make(network.QueryChannel)

	// create one or more network endpoints
	httpEndpoint := http.NewHttpEndpoint(*addr, *staticPath, !(*disableInfo))
	httpEndpoint.SendQueriesTo(queryChannel)

	err := server.Server(VERSION, *couchbaseSite, *defaultPoolName, queryChannel, queryTimeout)
	if err != nil {
		clog.Fatalf("Unable to run server, err: %v", err)
	}
}

func dumpOnSignal(signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	for _ = range c {
		pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	}
}
