package controller

import (
	"context"
	"errors"
	"net/http"
	"sync"

	_ "chat-app/docs"
	"chat-app/internal/models"
	"chat-app/internal/repo"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller struct {
	ctx    context.Context
	Cancel func()
	router *gin.Engine
	repo   *repo.Repo

	mu    sync.Mutex
	Rooms map[string]*models.Room
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
	ctx, cancel := context.WithCancel(context.Background())
	c := &Controller{
		Rooms:  make(map[string]*models.Room),
		ctx:    ctx,
		Cancel: cancel,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.router == nil {
		return nil, errors.New("router cannot be nil")
	}
	if c.repo == nil {
		return nil, errors.New("repo cannot be nil")
	}

	return c, nil
}

func (c *Controller) RegisterRoutes() {
	c.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	c.router.GET("/", index)

	c.router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := c.router.Group("/api/v1")
	{
		api.GET("/health", c.Health)
		api.GET("/rooms", c.GetRooms)
		api.GET("/rooms/:room/bind", c.BindRoom)
		api.GET("/rooms/:room/:nickname/send", c.SendMessage)
	}
}

func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}
