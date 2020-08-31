package sessions

import (
	"MDS/tools"
	"errors"
	"net/http"
	"sync"
)

// Sessions - информация о сессиях
var sessions container

// информация о сессиях
type container struct {
	waitingSessions    map[SessionID]*Session
	processingSessions map[SessionID]*Session
	completedSessions  map[SessionID]*Session

	mutex sync.Mutex
}

// Bootstrap - Начальная инициалиация
func Bootstrap(config *tools.Config) {
	sessions.waitingSessions = make(map[SessionID]*Session)
	sessions.processingSessions = make(map[SessionID]*Session)
	sessions.completedSessions = make(map[SessionID]*Session)
}

// Take - взять сессию в работу
func Take(ID SessionID) (result *Session) {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()

	session := sessions.waitingSessions[ID]
	if session == nil {
		return nil
	}

	// перемещается в processingSessions
	sessions.processingSessions[ID] = session
	delete(sessions.waitingSessions, ID)

	return session
}

// Done - Закончить работу с сессией
func Done(s *Session) {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()

	// перемещается в completedSessions
	sessions.completedSessions[s.ID] = s
	delete(sessions.processingSessions, s.ID)
}

// Get - получить результат расчета
func Get(ID SessionID) (data []byte, httpStatus int, logError error, err error) {
	if len(ID) == 0 {
		err = errors.New("Empty session")
		logError = err
		httpStatus = http.StatusBadRequest
	} else if len(ID) > MaxSessionIDLength {
		err = errors.New("Session ID too long")
		logError = err
		httpStatus = http.StatusBadRequest
	}
	if err != nil {
		return
	}

	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()

	session := sessions.completedSessions[ID]
	if session == nil {
		if sessions.processingSessions[ID] == nil {
			err = errors.New("Session not found: " + string(ID))
			logError = err
			httpStatus = http.StatusNotFound
		} else {
			err = errors.New("Data not ready: " + string(ID))
			logError = err
			httpStatus = http.StatusNoContent
		}
	} else {
		delete(sessions.completedSessions, ID)
		err = session.Error
		data = session.Data
		httpStatus = session.HTTPStatus
	}

	return
}

// Create - создать новую сесиию
func Create() SessionID {
	session := NewSession()
	sessions.mutex.Lock()
	sessions.waitingSessions[session.ID] = session
	sessions.mutex.Unlock()
	return session.ID
}
