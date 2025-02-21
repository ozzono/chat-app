package controller

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"chat-app/internal/models"
	"chat-app/pkg/queue"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Health godoc
// @Summary Health check
// @Description Get the health status of the service
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func (c *Controller) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "up"})
}

// GetRooms godoc
// @Summary Get chat rooms
// @Description List available chat rooms
// @Tags chat
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Router /api/v1/rooms [get]
func (c *Controller) GetRooms(ctx *gin.Context) {
	rooms, err := c.repo.GetRooms()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rooms)
}

// // CreateRoom godoc
// // @Summary Create chat room
// // @Description Create a new chat room
// // @Tags chat
// // @Accept  json
// // @Produce  json
// // @Param room body map[string]string true "Room name"
// // @Success 201 {string} string "Room created"
// // @Router /api/v1/rooms [post]
// func (c *Controller) CreateRoom(ctx *gin.Context) {
// 	var req struct {
// 		ID string `json:"id"`
// 	}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if req.ID == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
// 		return
// 	}

// 	if err := c.repo.AddRoom(req.ID); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	upgrader := websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}

// 	ctx.Writer.Header().Del("Content-Type")
// 	ctx.Writer.Header().Set("Connection", "Upgrade")
// 	ctx.Writer.Header().Set("Upgrade", "websocket")
// 	log.Println("header", ctx.Request.Header)
// 	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to establish websocket connection"})
// 		return
// 	}

// 	// conn, err := utils.NewSocketConnection(ctx.Writer, ctx.Request)
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish WebSocket connection"})
// 	// 	return
// 	// }

// 	room := c.NewRoom(req.ID, conn)
// 	room.Worker.StartWorker(c.ctx)

// 	ctx.JSON(http.StatusCreated, gin.H{"msg": "Room created"})
// }

type messageTask struct {
	message   models.Message
	connPool  []*websocket.Conn
	execCount int
	mu        sync.Mutex
}

func NewMsgTask(msg models.Message, conn []*websocket.Conn) queue.Task {
	return &messageTask{
		message:  msg,
		connPool: conn,
	}
}

func (t *messageTask) Log() {
	log.Printf("message from %s reached execution limit", t.message.Nickname)
}

func (t *messageTask) Action(ctx context.Context) error {

	if t.connPool == nil {
		return nil
	}
	if t.ExecCount() >= queue.ExecutionLimit {
		return errors.New("message reached execution limit")
	}

	defer t.AddExecCount()

	for i, conn := range t.connPool {
		err := conn.WriteMessage(websocket.TextMessage, []byte(t.message.Fmt()))
		if err != nil {
			log.Printf("error sending message [%d] from %s: %v", i, t.message.Nickname, err)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}

	return nil
}

func (t *messageTask) ExecCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.execCount
}

func (t *messageTask) AddExecCount() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.execCount++
}
