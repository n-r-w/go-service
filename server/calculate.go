package server

import (
	"MDS/sessions"
	"MDS/tools"
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

// обработка запроса на расчет
func processReportRequest(body *[]byte, w http.ResponseWriter) {
	ID := sessions.Create()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ID))

	workersWaitGroup.Add(1)
	go requestWorker(body, ID)
}

// функция расчета
func requestWorker(body *[]byte, ID sessions.SessionID) {
	defer workersWaitGroup.Done()

	tools.Log("New Session: " + string(ID))
	session := sessions.Take(ID)
	if session == nil {
		panic("session is nil")
	}

	var err error
	defer func() {
		if err != nil {
			session.Error = err
			session.HTTPStatus = http.StatusInternalServerError
		}

		session.ValidTo = time.Now().Add(time.Second * sessions.SessionLifetimeSec)
		sessions.Done(session)
		tools.Log("Created: " + string(ID))
	}()

	// заглушка для расчета: парсим XML и создаем его снова на основе результата парсинга
	var template []byte
	template, err = ioutil.ReadFile("document.xml")
	if err == nil {
		buf := bytes.NewBuffer(template)
		dec := xml.NewDecoder(buf)

		var n Node
		err = dec.Decode(&n)
		if err != nil {
			return
		}

		var output []byte
		output, err = xml.Marshal(n)
		if err != nil {
			return
		}

		session.Data = []byte("processed: " + string(output))
		session.HTTPStatus = http.StatusCreated
	}

}
