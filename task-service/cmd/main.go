package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	appsvc "github.com/oziev02/taskflow-microservices/task-service/internal/application/task"
	domain "github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
	"github.com/oziev02/taskflow-microservices/task-service/internal/infrastructure/kafka"
	"github.com/oziev02/taskflow-microservices/task-service/internal/infrastructure/postgres"
	redisCache "github.com/oziev02/taskflow-microservices/task-service/internal/infrastructure/redis"
	httphandler "github.com/oziev02/taskflow-microservices/task-service/internal/interfaces/http"
)

func main() {
	_ = godotenv.Load()

	// Подключение к Postgres
	db := connectPostgres()
	defer db.Close()

	// Kafka publisher
	kafkaPublisher := kafka.NewTaskPublisher(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_TOPIC_TASK_CREATED"),
		os.Getenv("KAFKA_TOPIC_TASK_UPDATED"),
		os.Getenv("KAFKA_TOPIC_TASK_DELETED"),
	)

	// Redis cache (TTL = 60s)
	var cache domain.Cache = nil
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		redisDB := mustAtoi(os.Getenv("REDIS_DB"))
		cache = redisCache.NewTaskCache(redisAddr, redisDB, 60*time.Second)
	}

	// Слои
	repo := postgres.NewTaskRepository(db)
	service := appsvc.NewService(repo, kafkaPublisher, cache)

	// HTTP
	router := chi.NewRouter()
	httphandler.RegisterTaskRoutes(router, service)

	port := os.Getenv("PORT")
	log.Printf("Starting task-service on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func connectPostgres() *sqlx.DB {
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	return db
}

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
