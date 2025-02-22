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
	testNickname = "testuser"
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

func (suite *HandlersTestSuite) Test1Health() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	suite.NoError(err)
	suite.router.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.JSONEq(`{"status": "up"}`, w.Body.String())
}

func (suite *HandlersTestSuite) Test2BindRoom() {
	fmt.Println("starting Test2BindRoom")
	bindingUrl := fmt.Sprintf("ws%s/api/v1/rooms/%s/bind?nickname=%s", strings.TrimPrefix(suite.server.URL, "http"), testRoom.ID, testNickname)
	ws1, _, err := websocket.DefaultDialer.Dial(bindingUrl, nil)
	suite.NoError(err)
	fmt.Println("dial no err")
	defer ws1.Close()

	readingURL := fmt.Sprintf("ws%s/api/v1/rooms/%s/%s/send?content=hello", strings.TrimPrefix(suite.server.URL, "http"), testRoom.ID, testNickname)
	_, _, _ = websocket.DefaultDialer.Dial(readingURL, nil)

	_, msg, err := ws1.ReadMessage()
	suite.NoError(err)

	pattern := fmt.Sprintf(`^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\] %s: %s$`, testNickname, "hello")
	suite.Regexp(pattern, string(msg))
}

func (suite *HandlersTestSuite) Test3GetRooms() {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/rooms", nil)
	suite.NoError(err)
	suite.router.ServeHTTP(rec, req)
	suite.Equal(200, rec.Code)

	reqBody, err := io.ReadAll(rec.Result().Body)
	suite.NoError(err)
	reqRoom := []string{}
	err = json.Unmarshal(reqBody, &reqRoom)
	suite.NoError(err)
	suite.Equal(testRoom.ID, reqRoom[0])
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
