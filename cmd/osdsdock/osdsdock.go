// Copyright 2019 The SODA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This module implements a entry into the SODA REST service.

*/

package main

import (
	"flag"

	"github.com/sodafoundation/dock/pkg/dock"
	"github.com/sodafoundation/dock/pkg/model"
	"github.com/sodafoundation/dock/pkg/utils/constants"

	"github.com/sodafoundation/dock/pkg/db"
	. "github.com/sodafoundation/dock/pkg/utils/config"
	"github.com/sodafoundation/dock/pkg/utils/daemon"
	"github.com/sodafoundation/dock/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load()

	// Parse some configuration fields from command line. and it will override the value which is got from config file.
	flag.StringVar(&CONF.OsdsDock.ApiEndpoint, "api-endpoint", CONF.OsdsDock.ApiEndpoint, "Listen endpoint of dock service")
	flag.StringVar(&CONF.OsdsDock.DockType, "dock-type", CONF.OsdsDock.DockType, "Type of dock service")
	flag.BoolVar(&CONF.OsdsDock.Daemon, "daemon", CONF.OsdsDock.Daemon, "Run app as a daemon with -daemon=true")
	flag.DurationVar(&CONF.OsdsDock.LogFlushFrequency, "log-flush-frequency", CONF.OsdsDock.LogFlushFrequency, "Maximum number of seconds between log flushes")
	flag.Parse()

	daemon.CheckAndRunDaemon(CONF.OsdsDock.Daemon)
}

func main() {
	// Open SODA dock service log file.
	logs.InitLogs(CONF.OsdsDock.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// FixMe: osdsdock attacher service needs to specify the endpoint via configuration file,
	//  so add this temporarily.
	listenEndpoint := constants.OpensdsDockBindEndpoint
	if CONF.OsdsDock.DockType == model.DockTypeAttacher {
		listenEndpoint = CONF.OsdsDock.ApiEndpoint
	}
	// Construct dock module grpc server struct and run dock server process.
	ds := dock.NewDockServer(CONF.OsdsDock.DockType, listenEndpoint)
	if err := ds.Run(); err != nil {
		panic(err)
	}
}
