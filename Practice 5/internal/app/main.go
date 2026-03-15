package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"prac5/internal/handlers"
	"prac5/internal/middleware"
	"prac5/internal/repository"
	"prac5/internal/repository/_postgres"
	"prac5/internal/usecase"
	"prac5/pkg/modules"
)

func Run() {
	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(dbConfig)
	repositories := repository.NewRepositories(_postgre)
	userUsecase := usecase.NewUserUsecase(repositories.UserRepository)
	handler := handlers.NewHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Healthcheck)
	mux.HandleFunc("GET /users", handler.GetPaginatedUsers)
	mux.HandleFunc("GET /users/common-friends", handler.GetCommonFriends)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByID)
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("PATCH /users/{id}", handler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)

	withAuth := middleware.AuthMiddleware(mux)
	withLogging := middleware.LoggingMiddleware(withAuth)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", withLogging))
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "mydb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}
