package server

import (
	"MDS/tools"
	"context"
	"log"
	"net/http"
	"time"
)

// корректная остановка сервиса
func processExit(server *http.Server, config *tools.Config, exitCode int) {
	log.Println("Stopping...")

	// останавливаем веб сервер
	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeoutSec*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// ждем окончания расчета
	workersWaitGroup.Wait()

	log.Println("Stopped")

	// сообщаем что можно завершить программу
	shutdownChannel <- exitCode
}
