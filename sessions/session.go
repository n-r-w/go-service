package sessions

import (
	"time"

	"github.com/pborman/uuid"
)

// SessionID идентификатор сессии
type SessionID string

const (
	// MaxSessionIDLength - максимальная длина идентификатора сессии
	MaxSessionIDLength = 100
	// MaxBodyLength - максимальная длина тела запроса
	MaxBodyLength = 100
	// SessionLifetimeSec - время жизни результата расчета сессии в секундах
	SessionLifetimeSec = 5
)

// Session - информация о сессии
type Session struct {
	ID      SessionID // уникальный id
	ValidTo time.Time // хранить до указанного времени

	Error      error  // ошибка (если была)
	Data       []byte // результат расчета (если не было ошибки)
	HTTPStatus int    // http статус расчета
}

// NewSession - конструктор сессии
func NewSession() *Session {
	s := new(Session)
	s.ID = SessionID(uuid.NewRandom().String())
	return s
}
