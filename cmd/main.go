package main

import (
	_ "chat-app/docs" // replace with the actual path to your docs
	"log"

	"chat-app/internal/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Chat App API
// @version 1.0
// @description This is a sample server for a chat application.
// @host localhost:8080
// @BasePath /

func main() {
	r := gin.Default()

	ctrl, err := controller.NewController(controller.WithRouter(r))
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}
	ctrl.RegisterRoutes()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket endpoint
	r.GET("/ws", ctrl.WebSocketHandler)

	r.Run(":8080")
}
