package server

import (
	"MDS/sessions"
	"MDS/tools"
	"sync"
)

// начальная инициализация
func bootstrap(config *tools.Config) {
	workersWaitGroup = &sync.WaitGroup{}
	shutdownChannel = make(chan int, 1)

	sessions.Bootstrap(config)
	sessions.RunGarbageCollector(config)
}
