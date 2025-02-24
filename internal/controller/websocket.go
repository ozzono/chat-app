package controller

import (
	"chat-app/internal/models"
	"chat-app/pkg/bot"
	"chat-app/pkg/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// SendMessage sends a message to a specific room
//
//	@Summary		Send a message to a specific room
//	@Description	Send a message to a specific room identified by room ID
//	@Tags			websocket
//	@Accept			json
//	@Produce		json
//	@Param			room		path		string			true	"room ID"
//	@Param			nickname	path		string			true	"nickname"
//	@Param			payload		body		models.Message	true	"Payload with nickname and message"
//	@Success		200			{object}	map[string]string{}
//	@Failure		400			{object}	map[string]string{}
//	@Failure		404			{object}	map[string]string{}
//	@Failure		500			{object}	map[string]string{}
//	@Router			/api/v1/rooms/{room}/send [get]
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

	if strings.HasPrefix(content, "/") {
		botMsg, err := bot.ProcessCMD(content)
		if err != nil {
			log.Println("bot cmd process err", err)
			return
		}
		botMsg.Room = roomID
		defer func() {
			room.Worker.TaskQueue <- NewMsgTask(botMsg, room.Connection)
		}()
	}

	room.Worker.TaskQueue <- NewMsgTask(message, room.Connection)

	if err := c.repo.AddMessage(message); err != nil {
		log.Printf("error adding message to the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add message to the database"})
		return
	}
	log.Printf("Message sent to %s room: %s", roomID, message.Content)
	ctx.Done()
}

// BindRoom godoc
//
//	@Summary		Bind to chat room
//	@Description	Bind to a given chat room
//	@Tags			websocket
//	@Accept			json
//	@Produce		json
//	@Param			room		path		string	true	"Room name"
//	@Param			nickname	query		string	true	"Nickname"
//	@Success		200			{string}	string	"Connected"
//	@Router			/api/v1/rooms/{room}/bind [get]
func (c *Controller) BindRoom(ctx *gin.Context) {
	roomID := ctx.Param("room")

	conn, err := utils.NewSocketConnection(ctx.Writer, ctx.Request)
	if err != nil {
		log.Printf("error establishing websocket connection for room %s: %v", roomID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to establish websocket connection"})
		return
	}

	room, exists := c.GetRoom(roomID)
	if !exists {
		room = c.NewRoom(roomID)
		log.Printf("new websocket connection established for room: %s", roomID)
		go room.Worker.StartWorker(c.ctx)
		err = c.repo.AddRoom(room.ID)
		if err != nil {
			log.Printf("error room to db %s: %v", room.ID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add room to db"})
			return
		}
	}

	room.AddConnection(conn)

	if exists {
		msgs, err := c.repo.GetMessages(room.ID)
		if err != nil {
			log.Printf("error getting %s room msgs: %v", room.ID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get room msgs from db"})
			return
		}
		for _, msg := range msgs {
			m := msg.Fmt()
			if err = conn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
				log.Printf("failed do send history message to socket; msg %s - err: %v", m, err)
			}
		}
	}
	if err = conn.WriteMessage(websocket.TextMessage, []byte("chat loaded")); err != nil {
		log.Printf("failed do send history message to socket; msg %s - err: %v", "chat loaded", err)
	}
}
