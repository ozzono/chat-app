package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chat-app/internal/repo"
	"chat-app/pkg/queue"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite
	router *gin.Engine
	server *httptest.Server
	repo   *repo.Repo
}

func (suite *HandlersTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	// Initialize the repository
	var err error
	suite.repo, err = repo.NewRepo(":memory:")
	if err != nil {
		suite.T().Fatal(err)
	}

	ctrl, _ := NewController(WithRouter(suite.router), WithRepo(suite.repo))
	ctrl.RegisterRoutes()
	suite.server = httptest.NewServer(suite.router)
}

func (suite *HandlersTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *HandlersTestSuite) SetupTest() {
	// Reset the rooms map before each test
	rooms = make(map[string]*Room)

	// Reset the database before each test
	suite.repo.DB.Exec("DELETE FROM rooms")
}

func (suite *HandlersTestSuite) TestHealthHandler() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)
	assert.JSONEq(suite.T(), `{"status": "up"}`, w.Body.String())
}

func (suite *HandlersTestSuite) TestGetRoomsHandler() {
	// Create a room first to ensure the list is not empty
	w := httptest.NewRecorder()
	body := `{"name": "testroom1"}`
	req, _ := http.NewRequest("POST", "/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	// Test getting rooms
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/rooms", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)
	assert.JSONEq(suite.T(), `["testroom1"]`, w.Body.String())
}

func (suite *HandlersTestSuite) TestCreateRoomHandler() {
	w := httptest.NewRecorder()
	body := `{"name": "testroom2"}`
	req, _ := http.NewRequest("POST", "/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)
	assert.JSONEq(suite.T(), `{"msg": "Room created"}`, w.Body.String())
}

func (suite *HandlersTestSuite) TestBindRoomHandler() {
	// Create a room first
	w := httptest.NewRecorder()
	body := `{"name": "testroom3"}`
	req, _ := http.NewRequest("POST", "/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	// Initialize the room in the rooms map
	rooms["testroom3"] = &Room{
		Connections: make(map[string]*websocket.Conn),
		TaskQueue:   make(chan queue.Task),
	}

	// Start the worker for the room
	worker := queue.NewWorker("testroom3")
	worker.TaskQueue = rooms["testroom3"].TaskQueue
	go worker.StartWorker(context.Background())

	// Bind to the room using a real HTTP server
	u := "ws" + strings.TrimPrefix(suite.server.URL, "http") + "/rooms/testroom3/bind?nickname=user1"
	ws1, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(suite.T(), err)
	defer ws1.Close()

	u = "ws" + strings.TrimPrefix(suite.server.URL, "http") + "/rooms/testroom3/bind?nickname=user2"
	ws2, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(suite.T(), err)
	defer ws2.Close()

	// Send a message from ws1
	err = ws1.WriteMessage(websocket.TextMessage, []byte("hello"))
	assert.NoError(suite.T(), err)

	// Read the message from ws2
	_, msg, err := ws2.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "hello", string(msg))
}

func (suite *HandlersTestSuite) TestBindRoomHandler_NicknameTaken() {
	// Create a room first
	w := httptest.NewRecorder()
	body := `{"name": "testroom4"}`
	req, _ := http.NewRequest("POST", "/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	// Initialize the room in the rooms map
	rooms["testroom4"] = &Room{
		Connections: make(map[string]*websocket.Conn),
		TaskQueue:   make(chan queue.Task),
	}

	// Start the worker for the room
	worker := queue.NewWorker("testroom4")
	worker.TaskQueue = rooms["testroom4"].TaskQueue
	go worker.StartWorker(context.Background())

	// Bind to the room using a real HTTP server
	u := "ws" + strings.TrimPrefix(suite.server.URL, "http") + "/rooms/testroom4/bind?nickname=user1"
	ws1, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(suite.T(), err)
	defer ws1.Close()

	// Attempt to bind with the same nickname
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/rooms/testroom4/bind?nickname=user1", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
	assert.JSONEq(suite.T(), `{"error": "Nickname is already taken"}`, w.Body.String())
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
