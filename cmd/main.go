package main

import (
	"log"

	"chat-app/internal/controller"
	"chat-app/internal/repo"

	"github.com/gin-gonic/gin"
)

//	@title			Chat App API
//	@version		1.0
//	@description	This is a sample server for a chat application.
//	@host			localhost:8080
//	@BasePath		/

func main() {
	r := gin.Default()

	repo, err := repo.NewRepo(":memory:")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	ctrl, err := controller.NewController(
		controller.WithRouter(r),
		controller.WithRepo(repo),
	)
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}
	ctrl.RegisterRoutes()

	// Swagger endpoint
	r.Run(":8080")
}
