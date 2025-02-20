package main

import (
	"log"

	"chat-app/internal/controller"

	"github.com/gin-gonic/gin"
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
	r.Run(":8080")
}
