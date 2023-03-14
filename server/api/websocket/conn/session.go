package conn

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	Id             string `json:"id"`
	ConnectionTime time.Time
}

func NewSession() *Session {
	session := Session{}
	session.Id = uuid.New().String()
	session.ConnectionTime = time.Now()
	return &session
}
