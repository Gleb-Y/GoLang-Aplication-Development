package app

import (
	"fmt"
	"log"
	"os"
	"practice-8/internal/handlers"
	"practice-8/internal/middleware"
	"practice-8/internal/repository"
	"practice-8/internal/repository/_postgres"
	"practice-8/internal/service"
	"practice-8/pkg/modules"

	"github.com/gin-gonic/gin"
)

func Run() {
	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(dbConfig)
	repositories := repository.NewRepositories(_postgre)
	userService := service.NewUserService(repositories.UserRepository)
	handler := handlers.NewHandler(userService)

	router := gin.Default()

	rateLimiterConfig := middleware.RateLimiterConfig{
		RequestsPerMinute: 60,
	}

	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Use(middleware.RateLimiter(rateLimiterConfig))
	{
		protected.GET("/rate", handler.GetRate)
	}

	fmt.Println("Server starting on :8080")
	log.Fatal(router.Run(":8080"))
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
