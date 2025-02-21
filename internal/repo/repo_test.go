package repo

import (
	"chat-app/internal/models"
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	testRoom = "testroom"
)

type RepoTestSuite struct {
	suite.Suite
	repo *Repo
}

func (suite *RepoTestSuite) SetupSuite() {
	var err error
	suite.repo, err = NewRepo(":memory:")
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *RepoTestSuite) Test1Room() {
	err := suite.repo.AddRoom(testRoom)
	suite.NoError(err)

	rooms, err := suite.repo.GetRooms()
	roomsMap := models.ToMap(rooms)
	suite.NoError(err)
	testRoom, found := roomsMap[testRoom]
	suite.True(found)
	suite.Equal(testRoom.ID, testRoom.ID)
}

func (suite *RepoTestSuite) Test2Message() {
	msg := models.Message{
		Room:     testRoom,
		Nickname: "user1",
		Content:  "Hello, world!",
	}
	err := suite.repo.AddMessage(msg)
	suite.NoError(err)

	messages, err := suite.repo.GetMessages(testRoom)
	suite.NoError(err)
	suite.Len(messages, 1)
	suite.Equal(msg.Content, messages[0].Content)
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
