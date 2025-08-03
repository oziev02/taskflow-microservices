package main

import (
	"github.com/gin-gonic/gin"
	"log"

	"github.com/oziev02/taskflow-microservices/api-gateway/internal/grpc"
	"github.com/oziev02/taskflow-microservices/api-gateway/internal/handler"
	"github.com/oziev02/taskflow-microservices/api-gateway/internal/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	taskHandler := handler.NewTaskHandler(grpc.NewTaskServiceClient())

	r.POST("/tasks", taskHandler.CreateTask)

	log.Println("Starting API Gateway on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
