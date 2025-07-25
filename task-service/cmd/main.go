package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"github.com/oziev02/taskflow-microservices/task-service/internal/application/task"
	taskdomain "github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
	"github.com/oziev02/taskflow-microservices/task-service/internal/infrastructure/kafka"
	"github.com/oziev02/taskflow-microservices/task-service/internal/infrastructure/postgres"
	"github.com/oziev02/taskflow-microservices/task-service/internal/interfaces/http"
)

func main() {
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Подключение к PostgreSQL
	db := connectPostgres()
	defer db.Close()

	// Инициализация Kafka publisher
	kafkaPublisher := kafka.NewTaskPublisher(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_TOPIC"),
	)

	// Репозиторий и use-case слой
	repo := postgres.NewTaskRepository(db)
	service := task.NewService(repo, kafkaPublisher)

	// HTTP маршруты
	router := chi.NewRouter()
	http.RegisterTaskRoutes(router, service)

	port := os.Getenv("PORT")
	log.Printf("Starting task-service on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func connectPostgres() *sqlx.DB {
	dbURL := "postgres://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("DB_NAME") +
		"?sslmode=disable"

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	return db
}
