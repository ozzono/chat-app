package controller

import (
	"chat-app/internal/models"
	"chat-app/pkg/queue"

	"github.com/gorilla/websocket"
)

func (c *Controller) NewRoom(roomID string, conn *websocket.Conn) *models.Room {
	c.mu.Lock()
	defer c.mu.Unlock()
	room := &models.Room{
		ID:         roomID,
		Worker:     queue.NewWorker(roomID),
		Connection: conn,
	}
	c.Rooms[roomID] = room
	return room
}

func (c *Controller) GetRoom(roomID string) (*models.Room, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	room, found := c.Rooms[roomID]
	return room, found
}
