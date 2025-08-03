package handler

import (
	"context"
	"net/http"

	"github.com/oziev02/taskflow-microservices/api-gateway/proto/taskpb"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	client taskpb.TaskServiceClient
}

func NewTaskHandler(client taskpb.TaskServiceClient) *TaskHandler {
	return &TaskHandler{client: client}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &taskpb.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
	}

	resp, err := h.client.CreateTask(context.Background(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": resp.GetTask()})
}
