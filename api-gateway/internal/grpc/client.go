package grpc

import (
	"log"

	"github.com/oziev02/taskflow-microservices/api-gateway/proto/taskpb"
	"google.golang.org/grpc"
)

type TaskClient struct {
	Client taskpb.TaskServiceClient
}

func NewTaskClient(addr string) *TaskClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to task-service: %v", err)
	}

	client := taskpb.NewTaskServiceClient(conn)
	return &TaskClient{Client: client}
}
