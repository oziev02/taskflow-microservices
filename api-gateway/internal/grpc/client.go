package grpc

import (
	"log"

	"github.com/oziev02/taskflow-microservices/api-gateway/proto/taskpb"

	"google.golang.org/grpc"
)

func NewTaskServiceClient() taskpb.TaskServiceClient {
	conn, err := grpc.Dial("task-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	return taskpb.NewTaskServiceClient(conn)
}
