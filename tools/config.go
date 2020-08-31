package tools

import "time"

// Config - параметры сервера
type Config struct {
	ConnHost           string
	ConnPort           int
	GarbageCollectSec  time.Duration
	ShutdownTimeoutSec time.Duration
}
