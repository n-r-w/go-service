package server

import (
	"MDS/sessions"
	"MDS/tools"
	"net/http"
)

// обработка запроса на выдачу результата расчета
func processResultRequest(body *[]byte, w http.ResponseWriter) {
	ID := sessions.SessionID(string(*body))
	data, status, logError, err := sessions.Get(ID)

	w.WriteHeader(status)

	if logError != nil {
		tools.Log(logError.Error())
	}

	if err != nil {
		w.Write([]byte(err.Error()))
	} else if data != nil {
		w.Write(data)
		tools.Log("Result delivered: " + string(ID))
	}

}
