package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func NewSocketConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error upgrading connection")
	}

	return conn, nil

}
