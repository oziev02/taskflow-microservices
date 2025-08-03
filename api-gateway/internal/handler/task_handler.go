package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oziev02/taskflow-microservices/api-gateway/proto/taskpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskHandler struct {
	Client taskpb.TaskServiceClient
}

func NewTaskHandler(client taskpb.TaskServiceClient) *TaskHandler {
	return &TaskHandler{Client: client}
}

func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
	tasks := r.Group("/tasks")
	tasks.POST("/", h.CreateTask)
	tasks.GET("/", h.GetAllTasks)
	tasks.PUT("/:id", h.UpdateTask)
	tasks.DELETE("/:id", h.DeleteTask)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req taskpb.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Client.CreateTask(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	res, err := h.Client.GetAllTasks(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req taskpb.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}
	req.Id = id

	res, err := h.Client.UpdateTask(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, err = h.Client.DeleteTask(context.Background(), &taskpb.DeleteTaskRequest{
		Id: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
