package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chat-app/internal/models"
	"chat-app/internal/repo"
	"chat-app/pkg/queue"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
)

var (
	testRoom = models.Room{
		ID:     "testroom",
		Worker: queue.NewWorker("testroom"),
	}
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

	ctrl, err := NewController(
		WithRouter(suite.router),
		WithRepo(suite.repo),
	)
	suite.NoError(err)
	ctrl.RegisterRoutes()
	suite.server = httptest.NewServer(suite.router)
}

func (suite *HandlersTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *HandlersTestSuite) Test1HealthHandler() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	suite.NoError(err)
	suite.router.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.JSONEq(`{"status": "up"}`, w.Body.String())
}

func (suite *HandlersTestSuite) Test2CreateRoomHandler() {
	w := httptest.NewRecorder()
	body := roomPayload(testRoom.ID)
	req, _ := http.NewRequest("POST", "/api/v1/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	suite.Equal(201, w.Code)
	suite.JSONEq(`{"msg": "Room created"}`, w.Body.String())
}

func (suite *HandlersTestSuite) Test3GetRoomsHandler() {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/rooms", nil)
	suite.NoError(err)
	suite.router.ServeHTTP(rec, req)
	suite.Equal(200, rec.Code)

	reqBody, err := io.ReadAll(rec.Result().Body)
	suite.NoError(err)
	reqRoom := []models.Room{}
	err = json.Unmarshal(reqBody, &reqRoom)
	suite.NoError(err)
	suite.Equal(testRoom.ID, reqRoom[0].ID)
}

func (suite *HandlersTestSuite) Test4BindRoomHandler() {
	fmt.Println("starting Test4BindRoomHandler")
	u := "ws" + strings.TrimPrefix(suite.server.URL, "http") + "/api/v1/rooms/testroom/bind?nickname=user1"
	ws1, _, err := websocket.DefaultDialer.Dial(u, nil)
	suite.NoError(err)
	fmt.Println("dial no err")
	defer ws1.Close()

	err = ws1.WriteMessage(websocket.TextMessage, []byte("hello"))
	suite.NoError(err)
	fmt.Println("write message no err")

	_, msg, err := ws1.ReadMessage()
	fmt.Println("read message no err")

	suite.NoError(err)
	suite.Equal("hello", string(msg))

	fmt.Println("close?")
	ws1.Close()
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func roomPayload(id string) string {
	return `{"id": "` + id + `"}`
}
