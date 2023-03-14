package conn

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Conn struct {
	Conn    *websocket.Conn
	Session *Session
	w       http.ResponseWriter
	r       *http.Request
	mutex   *sync.Mutex
}

func NewConn(conn *websocket.Conn, w http.ResponseWriter, r *http.Request) *Conn {
	connection := &Conn{
		Conn:    conn,
		Session: NewSession(),
		w:       w,
		r:       r,
		mutex:   &sync.Mutex{},
	}

	return connection
}

func (conn *Conn) WriteText(message []byte) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	err := conn.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Conn) WriteJSON(v interface{}) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	err := conn.Conn.WriteJSON(v)
	if err != nil {
		return err
	}
	return nil
}
