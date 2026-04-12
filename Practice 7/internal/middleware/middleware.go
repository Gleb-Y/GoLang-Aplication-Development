package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"practice-7/pkg/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates JWT token from Authorization header
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

// RoleMiddleware checks if user has the required role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*utils.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if userClaims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	RequestsPerMinute int
}

// RateLimiter middleware limits requests per minute per user or IP
func RateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	type RequestTracker struct {
		count     int
		resetTime time.Time
	}

	trackers := make(map[string]*RequestTracker)
	var mu sync.Mutex

	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		// Try to get user ID from claims, otherwise use ClientIP
		key := c.ClientIP()
		if claims, exists := c.Get("claims"); exists {
			if userClaims, ok := claims.(*utils.Claims); ok {
				key = fmt.Sprintf("user:%d", userClaims.UserID)
			}
		}

		now := time.Now()
		tracker, exists := trackers[key]

		// Initialize or reset tracker
		if !exists || now.After(tracker.resetTime) {
			trackers[key] = &RequestTracker{
				count:     1,
				resetTime: now.Add(1 * time.Minute),
			}
			c.Next()
			return
		}

		// Check if limit exceeded
		if tracker.count >= config.RequestsPerMinute {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		tracker.count++
		c.Next()
	}
}
