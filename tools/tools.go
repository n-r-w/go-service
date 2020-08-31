package tools

import (
	"log"
	"sync"
)

var logMutex sync.Mutex

// Log - синхронный вывод в журнал
func Log(text string) {
	logMutex.Lock()
	log.Println(text)
	logMutex.Unlock()
}
