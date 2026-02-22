package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"prac3/internal/handlers"
	"prac3/internal/middleware"
	"prac3/internal/repository"
	"prac3/internal/repository/_postgres"
	"prac3/internal/usecase"
	"prac3/pkg/modules"
	"time"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
	repositories := repository.NewRepositories(_postgre)
	userUsecase := usecase.NewUserUsecase(repositories.UserRepository)
	handler := handlers.NewHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Healthcheck)
	mux.HandleFunc("GET /users", handler.GetUsers)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByID)
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("PATCH /users/{id}", handler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)

	withAuth := middleware.AuthMiddleware(mux)
	withLogging := middleware.LoggingMiddleware(withAuth)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", withLogging))
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "postgres",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
