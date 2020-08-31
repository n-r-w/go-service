package sessions

import (
	"MDS/tools"
	"strconv"
	"time"
)

// RunGarbageCollector - запуск сбора мусора
func RunGarbageCollector(config *tools.Config) {
	go worker(config.GarbageCollectSec)
}

// Таймер сборки мусора
func worker(periodSec time.Duration) {
	duration := time.Second * periodSec
	timer := time.NewTimer(duration)
	for {
		<-timer.C
		collect()
		timer.Reset(duration)
	}
}

// Сборка мусора
func collect() {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()

	for ID, session := range sessions.completedSessions {
		if session.ValidTo.Before(time.Now()) {
			delete(sessions.completedSessions, ID)
			tools.Log("Removed by timeout (total " + strconv.Itoa(len(sessions.completedSessions)) + "): " + string(ID))
		}
	}
}
