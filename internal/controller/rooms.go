package controller

import (
	"chat-app/internal/models"
)

func (c *Controller) NewRoom(roomID string) *models.Room {
	c.mu.Lock()
	defer c.mu.Unlock()
	room := models.NewRoom(roomID)
	c.Rooms[roomID] = room
	return room
}

func (c *Controller) GetRoom(roomID string) (*models.Room, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	room, found := c.Rooms[roomID]
	return room, found
}
