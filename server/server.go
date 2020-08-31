package server

import (
	"MDS/tools"
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	// Собственно сам http сервер
	server *http.Server
	// группа для ожидания окончания всех расчетов
	workersWaitGroup *sync.WaitGroup
	// канал для корректного завершения сервиса
	shutdownChannel chan int
)

// Start - запуск сервера
func Start(config *tools.Config) {
	// начальная инициализация
	bootstrap(config)
	// создаем http сервер
	server = &http.Server{Addr: config.ConnHost + ":" + strconv.Itoa(config.ConnPort), Handler: nil}

	// остановка по нажатию кнопки
	go keyboardReader(config)

	// запускаем
	http.HandleFunc("/", mainHandler)
	log.Println("Listening on", config.ConnHost+":"+strconv.Itoa(config.ConnPort))
	log.Println("Press Q+Enter for exit: ")
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Println(err)
	} else {
		// сервер был остановлен через Shutdown ждем сигнала о завершении
		<-shutdownChannel
	}
}

// главный обработчик запросов
func mainHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	switch r.Method {
	case http.MethodPost:
		processReportRequest(&body, w)
	case http.MethodGet:
		processResultRequest(&body, w)
	default:
		processUnqnownRequest(r.Method, w)
	}
}

// обработка непонятных запросов
func processUnqnownRequest(method string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	tools.Log("Unknown request: " + method)
}

// остановка по нажатию кнопки
func keyboardReader(config *tools.Config) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if len(text) > 0 && "q" == strings.ToLower(text[0:1]) {
			processExit(server, config, 0)
		}
	}
}
