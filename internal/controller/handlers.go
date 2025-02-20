package controller

import (
	"context"
	"net/http"
	"time"

	"chat-app/internal/models"
	"chat-app/pkg/queue"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Room struct {
	Connections map[string]*websocket.Conn
	TaskQueue   chan queue.Task
}

var rooms = make(map[string]*Room)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HealthHandler godoc
// @Summary Health check
// @Description Get the health status of the service
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func (c *Controller) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "up"})
}

// GetRoomsHandler godoc
// @Summary Get chat rooms
// @Description List available chat rooms
// @Tags chat
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Router /api/v1/rooms [get]
func (c *Controller) GetRoomsHandler(ctx *gin.Context) {
	rooms, err := c.repo.GetRooms()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roomNames := make([]string, len(rooms))
	for i, room := range rooms {
		roomNames[i] = room.Name
	}

	ctx.JSON(http.StatusOK, roomNames)
}

// CreateRoomHandler godoc
// @Summary Create chat room
// @Description Create a new chat room
// @Tags chat
// @Accept  json
// @Produce  json
// @Param room body map[string]string true "Room name"
// @Success 201 {string} string "Room created"
// @Router /api/v1/rooms [post]
func (c *Controller) CreateRoomHandler(ctx *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.CreateRoom(req.Name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Initialize the room in the rooms map
	rooms[req.Name] = &Room{
		Connections: make(map[string]*websocket.Conn),
		TaskQueue:   make(chan queue.Task),
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg": "Room created"})
}

// BindRoomHandler godoc
// @Summary Bind to chat room
// @Description Bind to a given chat room
// @Tags chat
// @Accept  json
// @Produce  json
// @Param room path string true "Room name"
// @Param nickname query string true "Nickname"
// @Success 200 {string} string "Connected"
// @Router /api/v1/rooms/{room}/bind [get]
func (c *Controller) BindRoomHandler(ctx *gin.Context) {
	roomName := ctx.Param("room")
	nickname := ctx.Query("nickname")

	if nickname == "" {
		ctx.JSON(400, gin.H{"error": "Nickname is required"})
		return
	}

	room, exists := rooms[roomName]
	if !exists {
		ctx.JSON(404, gin.H{"error": "Room not found"})
		return
	}

	if _, taken := room.Connections[nickname]; taken {
		ctx.JSON(400, gin.H{"error": "Nickname is already taken"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	room.Connections[nickname] = conn
	defer func() {
		delete(room.Connections, nickname)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		message := models.Message{
			Room:      roomName,
			Nickname:  nickname,
			Timestamp: time.Now(),
			Content:   string(msg),
		}
		room.TaskQueue <- &messageTask{message: message, room: room}
	}
}

// WebSocketHandler godoc
// @Summary Handle WebSocket connections
// @Description Establish and handle WebSocket connections
// @Tags websocket
// @Accept  json
// @Produce  json
// @Router /api/v1/ws [get]
func (c *Controller) WebSocketHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

type messageTask struct {
	message models.Message
	room    *Room
}

func (t *messageTask) Action(ctx context.Context) error {
	for nick, conn := range t.room.Connections {
		if nick != t.message.Nickname {
			conn.WriteMessage(websocket.TextMessage, []byte(t.message.Content))
		}
	}
	return nil
}

func (t *messageTask) ExecCount() int {
	return 0
}

func (t *messageTask) AddExecCount() {}
