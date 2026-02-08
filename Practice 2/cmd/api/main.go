package main

import (
	"assignment/internal/handlers"
	"assignment/internal/middleware"
	"assignment/internal/models"
	"log"
	"net/http"
)

func main() {
	store := models.NewTaskStore()
	taskHandler := handlers.NewTaskHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/tasks", taskHandler)
	mux.Handle("/tasks/", taskHandler)

	handler := middleware.LoggingMiddleware(middleware.AuthMiddleware(mux))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
