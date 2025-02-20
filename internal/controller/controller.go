package controller

import (
	"errors"

	_ "chat-app/docs" // replace with the actual path to your docs
	"chat-app/internal/repo"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller struct {
	router *gin.Engine
	repo   *repo.Repo
}

type Option func(*Controller) error

func WithRouter(router *gin.Engine) Option {
	return func(c *Controller) error {
		c.router = router
		return nil
	}
}

func WithRepo(r *repo.Repo) Option {
	return func(c *Controller) error {
		c.repo = r
		return nil
	}
}

func NewController(opts ...Option) (*Controller, error) {
	c := &Controller{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.router == nil {
		return nil, errors.New("router cannot be nil")
	}
	return c, nil
}

func (c *Controller) RegisterRoutes() {
	c.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := c.router.Group("/api/v1")

	apiV1.GET("/ws", c.WebSocketHandler)

	apiV1.GET("/health", c.HealthHandler)
	apiV1.GET("/rooms", c.GetRoomsHandler)
	apiV1.POST("/rooms", c.CreateRoomHandler)
	apiV1.GET("/rooms/:room/bind", c.BindRoomHandler)
}
