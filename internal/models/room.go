package models

import (
	"chat-app/pkg/queue"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID         string
	Connection *websocket.Conn `json:"-"    gorm:"-"`
	Worker     *queue.Worker   `json:"-"    gorm:"-"`
}
