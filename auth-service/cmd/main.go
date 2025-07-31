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

	authapp "github.com/oziev02/taskflow-microservices/auth-service/internal/application/auth"
	jwtmgr "github.com/oziev02/taskflow-microservices/auth-service/internal/infrastructure/jwt"
	"github.com/oziev02/taskflow-microservices/auth-service/internal/infrastructure/postgres"
	redisstore "github.com/oziev02/taskflow-microservices/auth-service/internal/infrastructure/redis"
	httphandler "github.com/oziev02/taskflow-microservices/auth-service/internal/interfaces/http"
)

func main() {
	_ = godotenv.Load() // в контейнере .env уже передаётся через env_file

	db := connectPostgres()
	defer db.Close()

	// JWT manager
	accessTTLMin, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL_MIN"))
	refreshTTLH, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL_H"))
	jwt := jwtmgr.NewManager(os.Getenv("JWT_SECRET"), accessTTLMin, refreshTTLH)

	// Redis
	store := redisstore.NewRefreshStore(os.Getenv("REDIS_ADDR"), mustAtoi(os.Getenv("REDIS_DB")))

	// Repo + Service
	repo := postgres.NewUserRepository(db)
	svc := authapp.NewService(repo, jwt, store)

	// HTTP
	r := chi.NewRouter()
	httphandler.RegisterRoutes(r, svc, jwt)

	port := os.Getenv("PORT")
	log.Printf("Starting auth-service on port %s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func connectPostgres() *sqlx.DB {
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"

	var db *sqlx.DB
	var err error

	// 30 попыток с паузой 2s (~1 минута ожидания)
	for i := 1; i <= 30; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			return db
		}
		log.Printf("postgres not ready (try %d/30): %v", i, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("failed to connect to postgres after retries: %v", err)
	return nil
}

func mustAtoi(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}
