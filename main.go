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
	"log"

	"github.com/couchbaselabs/tuqtng/network"
	"github.com/couchbaselabs/tuqtng/network/http"
	"github.com/couchbaselabs/tuqtng/server"
)

var addr = flag.String("addr", ":8093", "HTTP listen address")
var couchbaseSite = flag.String("couchbase", "", "Couchbase Cluster Address (http://...) or dir:PATH")
var poolName = flag.String("pool", "default", "Pool")

func main() {
	flag.Parse()

	// create a QueryChannel
	queryChannel := make(network.QueryChannel)

	// create one or more network endpoints
	httpEndpoint := http.NewHttpEndpoint(*addr)
	httpEndpoint.SendQueriesTo(queryChannel)

	err := server.Server(*couchbaseSite, *poolName, queryChannel)
	if err != nil {
		log.Fatalf("Unable to run server, err: %v", err)
	}
}
