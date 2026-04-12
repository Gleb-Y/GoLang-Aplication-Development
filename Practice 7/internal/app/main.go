package app

import (
	"fmt"
	"log"
	"os"
	"practice-7/internal/handlers"
	"practice-7/internal/middleware"
	"practice-7/internal/repository"
	"practice-7/internal/repository/_postgres"
	"practice-7/internal/usecase"
	"practice-7/pkg/modules"

	"github.com/gin-gonic/gin"
)

func Run() {
	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(dbConfig)
	repositories := repository.NewRepositories(_postgre)
	userUsecase := usecase.NewUserUsecase(repositories.UserRepository)
	handler := handlers.NewHandler(userUsecase)

	router := gin.Default()

	// Rate limiter config
	rateLimiterConfig := middleware.RateLimiterConfig{
		RequestsPerMinute: 60,
	}

	// Public routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Protected routes (require JWT)
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Use(middleware.RateLimiter(rateLimiterConfig))
	{
		protected.GET("/me", handler.GetMe)
	}

	// Admin routes (require JWT and admin role)
	admin := router.Group("/")
	admin.Use(middleware.JWTAuthMiddleware())
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.Use(middleware.RateLimiter(rateLimiterConfig))
	{
		admin.PATCH("/users/promote/:id", handler.PromoteUser)
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
