package controller

import (
	"errors"

	"chat-app/internal/repo"

	"github.com/gin-gonic/gin"
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
	c.router.GET("/health", c.HealthHandler)
	c.router.GET("/rooms", c.GetRoomsHandler)
	c.router.POST("/rooms", c.CreateRoomHandler)
	c.router.GET("/rooms/:room/bind", c.BindRoomHandler)
}
