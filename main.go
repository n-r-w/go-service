package main

import (
	"MDS/server"
	"MDS/tools"
)

const (
	connHost           = "localhost"
	connPort           = 5000
	garbageCollectSec  = 9999999
	shutdownTimeoutSec = 10
)

func main() {
	config := new(tools.Config)
	config.ConnHost = connHost
	config.ConnPort = connPort
	config.GarbageCollectSec = garbageCollectSec
	config.ShutdownTimeoutSec = shutdownTimeoutSec

	server.Start(config)
}
