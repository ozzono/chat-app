package models

import (
	"chat-app/pkg/queue"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID         string
	Connection []*websocket.Conn `json:"-"    gorm:"-"`
	Worker     *queue.Worker     `json:"-"    gorm:"-"`
	mu         sync.Mutex
}

func NewRoom(roomID string) *Room {
	return &Room{
		ID:         roomID,
		Worker:     queue.NewWorker(roomID),
		Connection: []*websocket.Conn{},
		mu:         sync.Mutex{},
	}
}
