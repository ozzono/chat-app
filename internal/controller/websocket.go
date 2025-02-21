package controller

import (
	"chat-app/internal/models"
	"chat-app/pkg/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SendMessage sends a message to a specific room
// @Summary Send a message to a specific room
// @Description Send a message to a specific room identified by room ID
// @Tags websocket
// @Accept json
// @Produce json
// @Param room path string true "room ID"
// @Param nickname path string true "nickname"
// @Param payload body models.Message true "Payload with nickname and message"
// @Success 200 {object} map[string]string{}
// @Failure 400 {object} map[string]string{}
// @Failure 404 {object} map[string]string{}
// @Failure 500 {object} map[string]string{}
// @Router /api/v1/rooms/{room}/send [get]
func (c *Controller) SendMessage(ctx *gin.Context) {
	roomID := ctx.Param("room")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "roomID path parameter is required"})
		return
	}

	nickname := ctx.Param("nickname")
	if nickname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nickname path parameter is required"})
		return
	}

	content := ctx.Query("content")
	if content == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "content query parameter is required"})
		return
	}

	room, found := c.GetRoom(roomID)
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	message := models.Message{
		Nickname:  nickname,
		Room:      roomID,
		Timestamp: time.Now().UTC(),
		Content:   content,
	}

	room.Worker.TaskQueue <- NewMsgTask(message, room.Connection)

	if err := c.repo.AddMessage(message); err != nil {
		log.Printf("Error adding message to the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add message to the database"})
		return
	}
	log.Printf("Message sent to %s room: %s", roomID, message.Content)
	ctx.Done()
	// ctx.JSON(http.StatusOK, gin.H{"msg": "message sent successfully"})
}

// BindRoom godoc
// @Summary Bind to chat room
// @Description Bind to a given chat room
// @Tags chat
// @Accept  json
// @Produce  json
// @Param room path string true "Room name"
// @Param nickname query string true "Nickname"
// @Success 200 {string} string "Connected"
// @Router /api/v1/rooms/{room}/bind [get]
func (c *Controller) BindRoom(ctx *gin.Context) {
	roomID := ctx.Param("room")

	conn, err := utils.NewSocketConnection(ctx.Writer, ctx.Request)
	if err != nil {
		log.Printf("error establishing websocket connection for room %s: %v", roomID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to establish websocket connection"})
		return
	}

	room, exists := c.GetRoom(roomID)
	defer func() {
		room.AddConnection(conn)
		log.Printf("new websocket connection established for room: %s", roomID)
	}()
	if exists {
		return
	}
	room = c.NewRoom(roomID)
	go room.Worker.StartWorker(c.ctx)
	err = c.repo.AddRoom(room.ID)
	if err != nil {
		log.Printf("error room to db %s: %v", room.ID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add room to db"})
		return
	}
}

// // WebSocket handler to register connections
// func (c *Controller) RegisterConnection(ctx *gin.Context) {
// 	// Extract listener ID from query
// 	listenerID := ctx.Query("listener")
// 	if listenerID == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "listener query parameter is required"})
// 		return
// 	}

// 	log.Printf("New WebSocket connection established for listener: %s", listenerID)

// 	c.ManageConnection(conn, listenerID)
// 	ctx.Status(http.StatusOK)
// }
