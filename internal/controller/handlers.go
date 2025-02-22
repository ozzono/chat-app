package controller

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"chat-app/internal/models"
	"chat-app/pkg/queue"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Health godoc
//
//	@Summary		Health check
//	@Description	Get the health status of the service
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/api/v1/health [get]
func (c *Controller) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "up"})
}

// GetRooms godoc
//
//	@Summary		Get chat rooms
//	@Description	List available chat rooms
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.UIRoom
//	@Router			/api/v1/rooms [get]
func (c *Controller) GetRooms(ctx *gin.Context) {
	var rooms = []models.UIRoom{}
	for _, r := range c.Rooms {
		rooms = append(rooms, models.UIRoom{Room: r.ID, Users: len(r.Connection)})
	}
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Users > rooms[j].Users
	})
	ctx.JSON(http.StatusOK, rooms)
}

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
